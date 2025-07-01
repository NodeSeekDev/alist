package weiyunsdkgo

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// 微信扫码登录
// cookies:
// wx_login_ticket, access_token, refresh_token, openid, wy_appid, wyctoken, wx_uid, key_type
func (c *WeiYunClient) WeiXinQRLogin(ctx context.Context, showQR func([]byte) error) ([]*http.Cookie, error) {
RESTART:
	resp, err := c.Client.R().SetContext(ctx).
		SetQueryParams(map[string]string{
			"cmd":        "web_login",
			"no_referer": "1",
			"low_login":  "1",
		}).
		Get("https://user.weiyun.com/newcgi/web_wx_login.fcg")
	if err != nil {
		return nil, err
	}

	code := regexp.MustCompile(`/connect/qrcode/([^\s"']+)`).FindStringSubmatch(resp.String())[1]
	appid := regexp.MustCompile(`appid=([^&"']+)`).FindStringSubmatch(resp.String())[1]
	g_tk := regexp.MustCompile(`g_tk=([^&"']+)`).FindStringSubmatch(resp.String())[1]
	//state := regexp.MustCompile(`state=([^&]+)`).FindStringSubmatch(resp.String())[1]

	// 获取二维码图像
	resp, err = c.Client.R().SetContext(ctx).Get("https://open.weixin.qq.com/connect/qrcode/" + code)
	if err != nil {
		return nil, err
	}

	// 显示二维码
	if err := showQR(resp.Body()); err != nil {
		return nil, err
	}

	for {
		resp, err = c.Client.R().SetContext(ctx).Get("https://lp.open.weixin.qq.com/connect/l/qrconnect?uuid=" + code)
		if err != nil {
			return nil, err
		}
		codes := regexp.MustCompile(`window.wx_errcode=(\d+);window.wx_code='(.*)'`).FindStringSubmatch(resp.String())
		if len(codes) == 0 {
			return nil, fmt.Errorf("return parameter matching error")
		}
		wxErrCode, wxCode := codes[1], codes[2]

		switch wxErrCode {
		case "408": // 等待扫描
		case "402": // 等待超时
			goto RESTART
		case "403": // 拒绝
			return nil, fmt.Errorf("user cancels this login")
		case "404": // 等待确认
			time.Sleep(time.Second)
		case "405": // 确认登录
			resp, err := c.Client.R().SetDoNotParseResponse(true).SetContext(ctx).
				SetQueryParams(map[string]string{
					"g_tk":   g_tk,
					"appid":  appid,
					"action": "web_login",
					"code":   wxCode,
					//"state":  state,
				}).Get("https://user.weiyun.com/newcgi/weixin_oauth20.fcg")
			if err != nil {
				return nil, err
			}

			cookies := resp.RawResponse.Request.Cookies()
			c.SetCookies(cookies)
			return cookies, nil
		default:
			return nil, fmt.Errorf("err code:%s", wxErrCode)
		}
	}
}

type WeiXinRefreshTokenData struct {
	Openid       string `json:"openid"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// 微信刷新Token
func (c *WeiYunClient) WeiXinRefreshToken() (*WeiXinRefreshTokenData, error) {
	var resp WeiXinRefreshTokenData
	resp_, err := c.Client.R().
		SetQueryParams(map[string]string{
			"grant_type":    "refresh_token",
			"appid":         GetCookieValue("wy_appid", c.GetCookies()),
			"refresh_token": GetCookieValue("refresh_token", c.GetCookies()),
		}).
		ForceContentType("application/json; charset=UTF-8").
		SetResult(&resp).
		Get("https://api.weixin.qq.com/sns/oauth2/refresh_token")
	if err != nil {
		return nil, err
	}

	if jsoniter.Get(resp_.Body(), "errcode").ToInt() != 0 {
		return nil, &Resp{Code: jsoniter.Get(resp_.Body(), "errcode").ToInt(), Msg: jsoniter.Get(resp_.Body(), "errmsg").ToString()}
	}
	cks := c.GetCookies()
	SetCookieValue("openid", resp.Openid, cks)
	SetCookieValue("access_token", resp.AccessToken, cks)
	SetCookieValue("refresh_token", resp.RefreshToken, cks)
	c.SetCookies(cks)
	return &resp, nil
}

func (c *WeiYunClient) QQLoginInit(ctx context.Context) (appid, daid, callbackURL string, err error) {
	// var resp *resty.Response
	// resp, err = c.Client.R().SetContext(ctx).Get("https://www.weiyun.com")
	// if err != nil {
	// 	return
	// }
	// appid = regexp.MustCompile(`appid=([^&"']+)`).FindStringSubmatch(resp.String())[1]
	// daid = regexp.MustCompile(`daid=([^&"']+)`).FindStringSubmatch(resp.String())[1]
	appid = "527020901"
	daid = "372"
	callbackURL = "https://www.weiyun.com/web/callback/common_qq_login_ok.html?login_succ"
	return
}

// QQ扫码登录
func (c *WeiYunClient) QQQRLogin(ctx context.Context, showQR func([]byte) error) ([]*http.Cookie, error) {
RESTART:
	appid, daid, callbackURL, err := c.QQLoginInit(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.R().SetContext(ctx).SetQueryParams(map[string]string{
		"appid": appid,
		"daid":  daid,
		"s_url": callbackURL,

		"low_login":         "1",
		"low_login_foreced": "1",
		"qlogin_auto_login": "1",
		// 页面显示
		"style":           "20",
		"hide_title":      "1",
		"target":          "self",
		"link_target":     "blank",
		"hide_close_icon": "1",
		"pt_no_auth":      "1",
	}).Get("https://ssl.xui.ptlogin2.weiyun.com/cgi-bin/xlogin")
	if err != nil {
		return nil, err
	}
	pt_login_sig := GetCookieValue("pt_login_sig", resp.Cookies())

	resp, err = c.Client.R().SetContext(ctx).SetQueryParams(map[string]string{
		"appid":      appid,
		"daid":       daid,
		"pt_3rd_aid": "0",

		"t":  RandomT(),
		"u1": callbackURL,

		"e": "2", // 二维码外边框大小
		"l": "M",
		"s": "8", // 二维码大小
		"d": "72",
		"v": "4",
	}).Get("https://ssl.ptlogin2.weiyun.com/ptqrshow")
	if err != nil {
		return nil, err
	}

	if err := showQR(resp.Body()); err != nil {
		return nil, err
	}

	qrsig := GetCookieValue("qrsig", resp.Cookies())
	ptqrtoken := GetHash33(qrsig)

	for {
		resp, err = c.Client.R().SetContext(ctx).SetQueryParams(map[string]string{
			"u1":        callbackURL,
			"aid":       appid, // appid
			"daid":      daid,
			"ptqrtoken": ptqrtoken,
			"action":    fmt.Sprintf("0-0-%s", strconv.FormatInt(time.Now().UnixNano(), 10)[:13]),
			"login_sig": pt_login_sig,

			"from_ui":          "1",
			"low_login_enable": "1",
			"low_login_hour":   "720",

			"ptlang":        "2052", //2052：简体中文 1028：繁体中文 1033：英文
			"ptredirect":    "0",
			"pt_uistyle":    "40",
			"pt_js_version": "v1.46.0",

			// "js_ver":  "23071715",
			// "js_type": "1",
			// "h":       "1",
			// "t":       "1",
			// "g":       "1",
			// "o1vId":   "e1ddfd8550156335a97b7ef2acaeef0d", // deviceID
		}).Get("https://ssl.ptlogin2.weiyun.com/ptqrlogin")
		if err != nil {
			return nil, err
		}

		ptuiCB := regexp.MustCompile(`ptuiCB\((.*)\)`).FindStringSubmatch(resp.String())[1]
		codes := MustSliceConvert(strings.Split(ptuiCB, ","), func(code string) string {
			return strings.Trim(strings.TrimSpace(code), "'")
		})
		switch codes[0] {
		case "66", "67": // 未失效  已扫描,但还未点击确认
			time.Sleep(time.Second)
		case "65": // 已失效
			goto RESTART
		case "0": // 已经点击确认,并登录成功
			cookies := resp.Cookies()
			// 获取 p_skey
			resp, err = c.Client.R().SetContext(ctx).Get(codes[2])
			if err != nil {
				return nil, err
			}
			cookies = append(cookies, resp.Cookies()...)
			c.SetCookies(cookies)
			return cookies, nil
		default:
			return nil, fmt.Errorf(resp.String())
		}
	}
}

// 本地qq快速登录
func (c *WeiYunClient) QQFastLogin(ctx context.Context, qq string) ([]*http.Cookie, error) {
	appid, daid, callbackURL, err := c.QQLoginInit(ctx)
	if err != nil {
		return nil, err
	}
	// 获取pt_local_tk
	resp, err := c.Client.R().SetContext(ctx).
		SetQueryParams(map[string]string{
			"appid": appid,
			"daid":  daid,
			"s_url": callbackURL,

			"low_login":         "1",
			"qlogin_auto_login": "1",
			"low_login_foreced": "1",
			// 页面显示
			"style":           "20",
			"hide_title":      "1",
			"target":          "self",
			"link_target":     "blank",
			"hide_close_icon": "1",
			"pt_no_auth":      "1",
		}).Get("https://ssl.xui.ptlogin2.weiyun.com/cgi-bin/xlogin")
	if err != nil {
		return nil, err
	}
	pt_local_tk := GetCookieValue("pt_local_token", resp.Cookies())

	// 获取所有本地QQ信息
	resp, err = c.Client.R().SetContext(ctx).
		SetHeader("Referer", "https://ssl.xui.ptlogin2.weiyun.com/").
		SetQueryParams(map[string]string{
			"pt_local_tk": pt_local_tk,
			"r":           RandomT(),
			"callback":    "ptui_getuins_CB",
		}).Get("https://localhost.ptlogin2.weiyun.com:4301/pt_get_uins")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("err http code: %d", resp.StatusCode())
	}
	if !strings.Contains(resp.String(), qq) {
		return nil, fmt.Errorf("local QQ not login")
	}

	// 登录指定qq
	resp, err = c.Client.R().SetContext(ctx).
		SetHeader("Referer", "https://ssl.xui.ptlogin2.weiyun.com/").
		SetQueryParams(map[string]string{
			"clientuin":   qq,
			"pt_local_tk": pt_local_tk,
			"r":           RandomT(),
			"callback":    "__jp0",
		}).Get("https://localhost.ptlogin2.weiyun.com:4301/pt_get_st")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("err http code: %d", resp.StatusCode())
	}

	// 获取cookie
	resp, err = c.Client.R().SetContext(ctx).
		SetQueryParams(map[string]string{
			"clientuin":   qq,
			"u1":          callbackURL,
			"pt_aid":      appid, // appid
			"daid":        daid,
			"pt_local_tk": pt_local_tk,

			"keyindex":   "9",
			"style":      "40",
			"ptopt":      "1",
			"pt_3rd_aid": "0",
		}).
		Get("https://ssl.ptlogin2.weiyun.com/jump")
	if err != nil {
		return nil, err
	}

	ptuiCB := regexp.MustCompile(`ptui_qlogin_CB\((.*)\)`).FindStringSubmatch(resp.String())[1]
	codes := MustSliceConvert(strings.Split(ptuiCB, ","), func(code string) string {
		return strings.Trim(strings.TrimSpace(code), "'")
	})

	switch codes[0] {
	case "0":
		cookies := resp.Cookies()
		// 获取 p_skey
		resp, err = c.Client.R().SetContext(ctx).Get(codes[1])
		if err != nil {
			return nil, err
		}
		cookies = append(cookies, resp.Cookies()...)
		c.SetCookies(cookies)
		return cookies, nil
	default:
		return nil, fmt.Errorf("login fail,msg: %s", resp.String())
	}
}

// 2:wx 1:qq 0:not login
func (c *WeiYunClient) LoginType() int8 {
	cks := c.GetCookies()
	if GetCookieValue("wy_appid", cks) != "" {
		return 2
	} else if GetCookieValue("p_skey", cks) != "" {
		return 1
	}
	return 0
}

func (c *WeiYunClient) ParseTokenInfo() Json {
	// 微信登录
	cks := c.GetCookies()
	if c.LoginType() == 2 {
		return Json{
			"token_type":      1,
			"openid":          GetCookieValue("openid", cks),
			"open_appid":      GetCookieValue("wy_appid", cks),
			"access_token":    GetCookieValue("access_token", cks),
			"login_key_type":  192,
			"login_key_value": GetCookieValue("access_token", cks),
		}
	}
	// qq登录
	return Json{
		"openid": "",

		"token_type":      0,
		"login_key_type":  27,
		"login_key_value": GetCookieValue("p_skey", cks),
	}
}
