package mopan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-version"
)

type RestyOption func(request *resty.Request)
type Json map[string]any

func NewMoClientWithAuthorization(authorization string) *MoClient {
	return NewMoClient().SetAuthorization(authorization)
}

func NewMoClientWithRestyClient(client *resty.Client) *MoClient {
	return &MoClient{
		Client:     client,
		DeviceInfo: DefaultDeviceInfo,
	}
}

func NewMoClient() *MoClient {
	return &MoClient{
		Client:     resty.New(),
		DeviceInfo: DefaultDeviceInfo,
	}
}

type MoClient struct {
	Authorization string
	DeviceInfo    DeviceInfo

	Client *resty.Client

	onAuthorizationExpired func(err error) error
	flag                   int32
}

// 当Token失效时回调
func (c *MoClient) SetOnAuthorizationExpired(f func(err error) error) *MoClient {
	c.onAuthorizationExpired = f
	return c
}

func (c *MoClient) SetDeviceInfo(info *DeviceInfo) *MoClient {
	if info != nil {
		c.DeviceInfo = *info
	}
	return c
}

func (c *MoClient) GetDeviceInfo() DeviceInfo {
	return c.DeviceInfo
}

func (c *MoClient) SetAuthorization(authorization string) *MoClient {
	if !strings.HasPrefix(authorization, "Bearer") {
		authorization = "Bearer " + authorization
	}
	c.Authorization = authorization
	return c
}

func (c *MoClient) SetClient(client *http.Client) *MoClient {
	c.Client = resty.NewWithClient(client)
	return c
}

func (c *MoClient) SetRestyClient(client *resty.Client) *MoClient {
	c.Client = client
	return c
}

func (c *MoClient) SetProxy(proxy string) *MoClient {
	c.Client.SetProxy(proxy)
	return c
}

func (c *MoClient) request(url string, data Json, resp any, option ...RestyOption) ([]byte, error) {
	secretKey := GetSecretKey()

	headers := map[string]string{
		"Authorization": c.Authorization,
		"remoteInfo":    c.DeviceInfo.Encrypt(secretKey),
	}

	ver := version.Must(version.NewVersion(c.DeviceInfo.MpVersion))
	if ver.GreaterThanOrEqual(version.Must(version.NewVersion("1.1.202"))) {
		encryptedKey := MustRsaEncryptBase64Str(secretKey, PublicKeyV2)
		headers["version"] = c.DeviceInfo.MpVersion
		headers["encrypted-key"] = encryptedKey
	} else {
		encryptedKey := MustRsaEncryptBase64Str(secretKey, PublicKeyV1)
		headers["encrypted-key"] = encryptedKey
	}

	req := c.Client.R().SetHeaders(headers)

	if data != nil {
		req.SetHeader("Content-Type", "application/json")
		temp, err := c.Client.JSONMarshal(data)
		if err != nil {
			return nil, err
		}
		enc, _ := AesEncryptBase64(temp, []byte(secretKey))
		req.SetBody(enc)
	}

	for _, opt := range option {
		opt(req)
	}
	resp_, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	body := resp_.Body()
	// 解密数据
	if bytes.HasPrefix(body, []byte{'"'}) && bytes.HasSuffix(body, []byte{'"'}) {
		body, err = AesDecryptBase64(bytes.Trim(body, "\""), []byte(secretKey))
		if err != nil {
			return nil, err
		}
	}

	var result Resp
	c.Client.JSONUnmarshal(body, &result)

	if resp_.StatusCode() == http.StatusUnauthorized {
		result.Code = 401
	}

	if result.Code != 200 {
		return nil, &result
	}

	if resp != nil {
		if err := c.Client.JSONUnmarshal(result.Data, &resp); err != nil {
			return nil, err
		}
	}
	return result.Data, nil
}

func (c *MoClient) Request(url string, data Json, resp any, option ...RestyOption) ([]byte, error) {
	v, err := c.request(url, data, resp, option...)
	if err != nil {
		if err, ok := err.(*Resp); ok {
			// 401 错误处理
			if err.Code == 401 && c.onAuthorizationExpired != nil {
				if atomic.CompareAndSwapInt32(&c.flag, 0, 1) {
					err2 := c.onAuthorizationExpired(err)
					atomic.SwapInt32(&c.flag, 0)
					if err2 != nil {
						return nil, errors.Join(err, err2)
					}
				}

				for atomic.LoadInt32(&c.flag) != 0 {
					runtime.Gosched()
				}

				return c.request(url, data, resp, option...)
			}
		}
		return nil, err
	}
	return v, nil
}

type Resp struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Status  bool            `json:"status"`
}

func (r *Resp) Error() string {
	return fmt.Sprintf("Code:%d, Message:%s", r.Code, r.Message)
}
