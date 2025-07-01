package weiyunsdkgo

import (
	"encoding/json"
	"fmt"
)

type RespHeader struct {
	Retcode int    `json:"retcode"`
	Retmsg  string `json:"retmsg"`

	Cmd     int    `json:"cmd"`
	CmdName string `json:"cmdName"`

	Uin  int64 `json:"uin"`
	Type int   `json:"type"`

	// flag
	ZipFlag int `json:"zip_flag"`
	Encrypt int `json:"encrypt"`

	Seq int64 `json:"seq"`
}

type Resp struct {
	Code int    `json:"ret"`
	Msg  string `json:"msg"`

	Data *struct {
		RspHeader RespHeader `json:"rsp_header"`
		RspBody   struct {
			RspMsgBody json.RawMessage `json:"RspMsg_body"`
		} `json:"rsp_body"`
	} `json:"data"`
}

func (r *Resp) HasError() bool {
	return r.Data.RspHeader.Retcode != 0 || r.Code != 0
}

func (r *Resp) Error() string {
	if r.Msg != "" {
		return fmt.Sprintf("errcode:%d, errmsg:%s", r.Code, r.Msg)
	}
	return fmt.Sprintf("cmd:%d, cmdName:%s, seq:%d, errcode:%d, errmsg:%s",
		r.Data.RspHeader.Cmd,
		r.Data.RspHeader.CmdName,
		r.Data.RspHeader.Seq,
		r.Data.RspHeader.Retcode,
		r.Data.RspHeader.Retmsg)
}

func (r *Resp) GetHeader() RespHeader {
	return r.Data.RspHeader
}

func (r *Resp) GetBody() json.RawMessage {
	return r.Data.RspBody.RspMsgBody
}

type OpenApiErron struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func (e *OpenApiErron) IsError() bool {
	return e.ErrCode != 0
}

func (e *OpenApiErron) Error() string {
	return fmt.Sprintf("code:%d, msg:%s", e.ErrCode, e.ErrMsg)
}

type OpenApiToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}
