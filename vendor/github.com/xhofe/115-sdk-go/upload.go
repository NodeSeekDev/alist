package sdk

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xhofe/115-sdk-go/json_types"
)

type UploadGetTokenResp struct {
	Endpoint        string `json:"endpoint"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"expiration"`
	AccessKeyId     string `json:"AccessKeyId"`
}

// UploadGetToken: https://www.yuque.com/115yun/open/kzacvzl0g7aiyyn4
func (c *Client) UploadGetToken(ctx context.Context) (*UploadGetTokenResp, error) {
	var resp UploadGetTokenResp
	_, err := c.AuthRequest(ctx, ApiFsUploadGetToken, http.MethodGet, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type UploadInitReq struct {
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	Target    string `json:"target"`
	FileID    string `json:"fileid"` // 文件Sha1
	PreID     string `json:"preid"`  // 文件前128k的sha1
	PickCode  string `json:"pick_code"`
	TopUpload string `json:"topupload"`
	SignKey   string `json:"sign_key"`
	SignVal   string `json:"sign_val"`
}

type UploadInitResp struct {
	PickCode  string `json:"pick_code"`
	Status    int    `json:"status"`
	SignKey   string `json:"sign_key"`
	SignCheck string `json:"sign_check"`
	FileID    string `json:"file_id"`
	Target    string `json:"target"`
	Bucket    string `json:"bucket"`
	Object    string `json:"object"`
	Callback  json_types.StructOrArray[struct {
		Callback    string `json:"callback"`
		CallbackVar string `json:"callback_var"`
	}] `json:"callback"`
}

// UploadInit: https://www.yuque.com/115yun/open/ul4mrauo5i2uza0q
func (c *Client) UploadInit(ctx context.Context, req *UploadInitReq) (*UploadInitResp, error) {
	var resp UploadInitResp
	_, err := c.AuthRequest(ctx, ApiFsUploadInit, http.MethodPost, &resp, ReqWithForm(Form{
		"file_name": req.FileName,
		"file_size": strconv.FormatInt(req.FileSize, 10),
		"target":    fmt.Sprintf("U_1_%s", req.Target),
		"fileid":    req.FileID,
		"preid":     req.PreID,
		"pick_code": req.PickCode,
		"topupload": req.TopUpload,
		"sign_key":  req.SignKey,
		"sign_val":  req.SignVal,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type UploadResumeReq struct {
	FileSize int64  `json:"file_size"`
	Target   string `json:"target"`
	FileID   string `json:"fileid"` // 文件Sha1
	PickCode string `json:"pick_code"`
}

type UploadResumeResp struct {
	PickCode string `json:"pick_code"`
	Target   string `json:"target"`
	Version  string `json:"version"`
	Bucket   string `json:"bucket"`
	Object   string `json:"object"`
	Callback json_types.StructOrArray[struct {
		Callback    string `json:"callback"`
		CallbackVar string `json:"callback_var"`
	}] `json:"callback"`
}

// UploadResume: https://www.yuque.com/115yun/open/tzvi9sbcg59msddz
func (c *Client) UploadResume(ctx context.Context, req *UploadResumeReq) (*UploadResumeResp, error) {
	var resp UploadResumeResp
	_, err := c.AuthRequest(ctx, ApiFsUploadResume, http.MethodPost, &resp, ReqWithForm(Form{
		"file_size": strconv.FormatInt(req.FileSize, 10),
		"target":    fmt.Sprintf("U_1_%s", req.Target),
		"fileid":    req.FileID,
		"pick_code": req.PickCode,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}
