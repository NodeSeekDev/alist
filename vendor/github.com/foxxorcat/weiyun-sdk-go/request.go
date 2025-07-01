package weiyunsdkgo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

func NewBody(cmdName string, data, tokenInfo Json) Json {
	json := Json{
		"ReqMsg_body": Json{
			"ext_req_head": Json{
				"token_info": tokenInfo,
				"language_info": Json{
					"language_type": 2052,
				},
			},
			".weiyun." + cmdName + "MsgReq_body": data,
		},
	}
	return json
}

func NewHeader(cmd int, tokenInfo Json) Json {
	json := Json{
		"seq":       time.Now().UnixNano(),
		"cmd":       cmd,
		"wx_openid": tokenInfo["openid"],
		"user_flag": tokenInfo["token_type"],

		"type":          1,
		"appid":         30013,
		"version":       3,
		"major_version": 3,
		"minor_version": 3,
		"fix_version":   3,
	}
	return json
}

func (c *WeiYunClient) request(protocol, cmdName string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	tokenInfo := c.ParseTokenInfo()
	req := c.Client.R().SetQueryParams(map[string]string{
		//"refer": "chrome_windows",
		"g_tk": GetCookieValue("wyctoken", c.GetCookies()),
		"cmd":  strconv.Itoa(cmd),
	})

	if protocol == "upload" {
		// 取出后删除
		var fileReader io.Reader
		if data != nil {
			if r, ok := data["fileReader"].(io.Reader); ok {
				delete(data, "fileReader")
				fileReader = r
			}
		}

		// 严格排序，只能手动构建
		boundary := "----WebKitFormBoundaryIifrOqiswelC8nfe"
		formData := bytes.NewBuffer(make([]byte, 0, 1024*1024+4096))
		formData.WriteString("--" + boundary + "\r\n")
		formData.WriteString("Content-Disposition: form-data; name=\"json\"\r\n\r\n")
		formData.WriteString(MustJsonMarshalToString(Json{
			"req_header": Json{
				"cmd":           cmd,
				"appid":         30013,
				"major_version": 3,
				"minor_version": 0,
				"fix_version":   0,
				"version":       0,
				"user_flag":     0,
			},
			"req_body": Json{
				"ReqMsg_body": Json{
					"weiyun." + cmdName + "MsgReq_body": data,
				},
			},
		}))
		if fileReader != nil {
			formData.WriteString("\r\n--" + boundary + "\r\n")
			formData.WriteString("Content-Disposition: form-data; name=\"upload\"; filename=\"blob\"\r\n")
			formData.WriteString("Content-Type: application/octet-stream\r\n\r\n")
			if _, err := io.Copy(formData, fileReader); err != nil {
				return nil, err
			}
		}
		formData.WriteString("\r\n--" + boundary + "--\r\n")

		req.SetBody(formData).SetHeader("Content-Type", "multipart/form-data; boundary="+boundary)
	} else {
		req.SetBody(Json{
			"req_header": MustJsonMarshalToString(NewHeader(cmd, tokenInfo)),
			"req_body":   MustJsonMarshalToString(NewBody(cmdName, data, tokenInfo)),
		})
	}

	for _, opt := range opts {
		opt(req)
	}

	var (
		respRaw Resp
		resp_   *resty.Response
		err     error
	)

	if protocol == "upload" {
		resp_, err = req.SetResult(&respRaw.Data).Post("https://upload.weiyun.com/ftnup_v2/weiyun")
		if err != nil {
			return nil, err
		}
	} else {
		resp_, err = req.
			SetPathParams(map[string]string{
				"protocol": protocol,
				"name":     cmdName,
			}).
			SetResult(&respRaw).
			Post("https://www.weiyun.com/webapp/json/{protocol}/{name}")
		if err != nil {
			return nil, err
		}
	}

	// http code 处理
	if resp_.StatusCode() != 200 {
		if resp_.StatusCode() == 403 {
			return resp_.Body(), ErrCode403
		}
		return resp_.Body(), fmt.Errorf("err http code: %d", resp_.StatusCode())
	}

	// resp.code 处理
	if respRaw.HasError() {
		return resp_.Body(), &respRaw
	}

	// 绑定body
	if resp != nil {
		if protocol == "upload" {
			jsoniter.Get(respRaw.Data.RspBody.RspMsgBody, "weiyun."+cmdName+"MsgRsp_body").ToVal(&resp)
		} else {
			err = c.Client.JSONUnmarshal(respRaw.GetBody(), resp)
			if err != nil {
				return resp_.Body(), err
			}
		}
	}
	return resp_.Body(), nil
}

func (c *WeiYunClient) RefreshCtoken() error {
	resp, err := c.Client.R().Get("https://www.weiyun.com/disk")
	if err != nil {
		return err
	}

	// 302跳转
	if resp.RawResponse.Request.URL != resp.Request.RawRequest.URL {
		return ErrCookieExpiration
	}
	return nil
}

func (c *WeiYunClient) KeepAlive() error {
	// TODO:
	// 登录时选择记住登录，可延长存活时间
	return c.RefreshCtoken()
}

func (c *WeiYunClient) Request(protocol, name string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	resp_, err := c.request(protocol, name, cmd, data, resp, opts...)
	if err == ErrCode403 {
		if atomic.CompareAndSwapInt32(&c.flag, 0, 1) {
			err2 := c.RefreshCtoken()
			// 如果是微信登录，尝试刷新Token
			if errors.Is(err2, ErrCookieExpiration) && c.LoginType() == 1 {
				_, err2 = c.WeiXinRefreshToken()
				if err2 == nil {
					// 验证Token刷新
					err2 = c.RefreshCtoken()
				}
				if err2 != nil {
					errors.Join(ErrCookieExpiration, err2)
				}
			}

			atomic.SwapInt32(&c.flag, 0)
			if err2 != nil {
				err = errors.Join(err, err2)
				c.onCookieExpired(err)
				return resp_, err
			}
		}
		for atomic.LoadInt32(&c.flag) != 0 {
			runtime.Gosched()
		}
		resp_, err = c.request(protocol, name, cmd, data, resp, opts...)
	}
	return resp_, err
}

func (c *WeiYunClient) WeiyunQdiskRequest(name string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	return c.Request("weiyunQdisk", name, cmd, data, resp, opts...)
}

func (c *WeiYunClient) WeiyunQdiskClientRequest(name string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	return c.Request("weiyunQdiskClient", name, cmd, data, resp, opts...)
}

func (c *WeiYunClient) WeiyunFileLibClientRequest(name string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	return c.Request("weiyunFileLibClient", name, cmd, data, resp, opts...)
}

func (c *WeiYunClient) UploadRequest(name string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	return c.Request("upload", name, cmd, data, resp, opts...)
}

func (c *WeiYunClient) WeiyunSafeBoxRequest(name string, cmd int, data Json, resp any, opts ...RestyOption) ([]byte, error) {
	return c.Request("weiyunSafeBox", name, cmd, data, resp, opts...)
}
