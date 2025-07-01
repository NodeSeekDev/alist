// docs: https://www.yuque.com/115yun/open/shtpzfhewv5nag11
package sdk

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

type AuthDeviceCodeResp struct {
	UID    string `json:"uid"`
	Time   int64  `json:"time"`
	QrCode string `json:"qrcode"`
	Sign   string `json:"sign"`
}

// $code_challenge = base64_encode(sha256($code_verifier));
func calCodeChanllenge(codeVerifier string) string {
	sha := sha256.New()
	sha.Write([]byte(codeVerifier))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}

func (c *Client) AuthDeviceCode(ctx context.Context, clientID string, codeVerifier string) (*AuthDeviceCodeResp, error) {
	var resp AuthDeviceCodeResp
	_, err := c.passportRequest(ctx, ApiAuthDeviceCode, http.MethodPost, &resp, ReqWithForm(Form{
		"client_id":             clientID,
		"code_challenge":        calCodeChanllenge(codeVerifier),
		"code_challenge_method": "sha256",
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type QrCodeStatusResp struct {
	Msg     string `json:"msg"`
	Status  int    `json:"status"`
	Version string `json:"version"`
}

func (c *Client) QrCodeStatus(ctx context.Context, uid, time, sign string) (*QrCodeStatusResp, error) {
	var resp QrCodeStatusResp
	_, err := c.passportRequest(ctx, ApiQrCodeStatus, http.MethodGet, &resp, ReqWithQuery(Form{
		"uid":  uid,
		"time": time,
		"sign": sign,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type CodeToTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (c *Client) CodeToToken(ctx context.Context, uid, codeVerifier string) (*CodeToTokenResp, error) {
	var resp CodeToTokenResp
	_, err := c.passportRequest(ctx, ApiCodeToToken, http.MethodPost, &resp, ReqWithForm(Form{
		"uid":           uid,
		"code_verifier": codeVerifier,
	}))
	if err != nil {
		return nil, err
	}
	// TODO: set token?
	c.SetAccessToken(resp.AccessToken)
	c.SetRefreshToken(resp.RefreshToken)
	return &resp, err
}

type RefreshTokenResp CodeToTokenResp

func (c *Client) RefreshToken(ctx context.Context) (*RefreshTokenResp, error) {
	var resp RefreshTokenResp
	_, err := c.passportRequest(ctx, ApiRefreshToken, http.MethodPost, &resp, ReqWithForm(Form{
		"refresh_token": c.refreshToken,
	}))
	if err != nil {
		return nil, err
	}
	c.SetAccessToken(resp.AccessToken)
	c.SetRefreshToken(resp.RefreshToken)
	if c.onRefreshToken != nil {
		c.onRefreshToken(resp.AccessToken, resp.RefreshToken)
	}
	return &resp, err
}
