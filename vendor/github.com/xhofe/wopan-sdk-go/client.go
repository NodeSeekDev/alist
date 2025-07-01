package wopan

import (
	"encoding/json"
	"net/http"
	"path"
	"sync"

	"github.com/go-resty/resty/v2"
)

type WoClient struct {
	accessToken  string
	refreshToken string

	client            *resty.Client
	crypto            *Crypto
	ua                string
	jsonMarshalFunc   func(v interface{}) ([]byte, error)
	jsonUnmarshalFunc func(data []byte, v interface{}) error

	Phone            string
	ZoneURL          string
	zoneURLOnce      sync.Once
	ClassifyRuleData *ClassifyRuleData

	onRefreshToken func(accessToken, refreshToken string)
}

func New(opts ...Option) *WoClient {
	w := &WoClient{
		client:            resty.New(),
		crypto:            NewCrypto(),
		jsonMarshalFunc:   json.Marshal,
		jsonUnmarshalFunc: json.Unmarshal,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func DefaultWithAccessToken(accessToken string) *WoClient {
	w := Default()
	w.SetAccessToken(accessToken)
	return w
}

func DefaultWithRefreshToken(refreshToken string) *WoClient {
	w := Default()
	w.SetRefreshToken(refreshToken)
	return w
}

func Default() *WoClient {
	return New(WithUA(DefaultUA))
}

func (w *WoClient) SetUA(ua string) {
	w.ua = ua
}

func (w *WoClient) SetJsonMarshalFunc(f func(v interface{}) ([]byte, error)) {
	w.jsonMarshalFunc = f
}

func (w *WoClient) SetJsonUnmarshalFunc(f func(data []byte, v interface{}) error) {
	w.jsonUnmarshalFunc = f
}

func (w *WoClient) SetAccessToken(token string) {
	w.accessToken = token
	_ = w.crypto.SetAccessToken(token)
}

func (w *WoClient) SetRefreshToken(token string) {
	w.refreshToken = token
}

func (w *WoClient) GetToken() (string, string) {
	return w.accessToken, w.refreshToken
}

func (w *WoClient) SetHttpClient(httpClient *http.Client) *WoClient {
	w.client = resty.NewWithClient(httpClient)
	return w
}

func (w *WoClient) SetUserAgent(userAgent string) *WoClient {
	w.client.SetHeader("User-Agent", userAgent)
	return w
}

func (w *WoClient) SetDebug(d bool) *WoClient {
	w.client.SetDebug(d)
	return w
}

func (w *WoClient) EnableTrace() *WoClient {
	w.client.EnableTrace()
	return w
}

func (w *WoClient) SetProxy(proxy string) *WoClient {
	w.client.SetProxy(proxy)
	return w
}

func (w *WoClient) NewRequest() *resty.Request {
	return w.client.R()
}

func (w *WoClient) GetFileType(filename string) string {
	ext := path.Ext(filename)
	if ext == "" {
		return "5"
	}
	ext = ext[1:]
	err := w.InitClassifyRule()
	if err != nil {
		return "5"
	}
	if _type, ok := w.ClassifyRuleData.FileTypes[ext]; ok {
		return _type.Type
	}
	return "5"
}

func (w *WoClient) OnRefreshToken(f func(accessToken, refreshToken string)) {
	w.onRefreshToken = f
}
