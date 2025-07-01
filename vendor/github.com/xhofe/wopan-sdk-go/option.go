package wopan

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Option func(w *WoClient)

func WithAccessToken(token string) Option {
	return func(w *WoClient) {
		w.SetAccessToken(token)
	}
}

func WithRefreshToken(token string) Option {
	return func(w *WoClient) {
		w.SetRefreshToken(token)
	}
}

func WithUA(ua string) Option {
	return func(w *WoClient) {
		w.SetUA(ua)
	}
}

func WithJsonMarshalFunc(f func(v interface{}) ([]byte, error)) Option {
	return func(w *WoClient) {
		w.SetJsonMarshalFunc(f)
	}
}

func WithJsonUnmarshalFunc(f func(data []byte, v interface{}) error) Option {
	return func(w *WoClient) {
		w.SetJsonUnmarshalFunc(f)
	}
}

func WithClient(hc *http.Client) Option {
	return func(c *WoClient) {
		c.SetHttpClient(hc)
	}
}

func WithRestyClient(rc *resty.Client) Option {
	return func(c *WoClient) {
		c.client = rc
	}
}

func WithDebug() Option {
	return func(c *WoClient) {
		c.SetDebug(true)
	}
}

func WithTrace() Option {
	return func(c *WoClient) {
		c.EnableTrace()
	}
}

func WithProxy(proxy string) Option {
	return func(c *WoClient) {
		c.SetProxy(proxy)
	}
}

type RestyOption func(request *resty.Request)
