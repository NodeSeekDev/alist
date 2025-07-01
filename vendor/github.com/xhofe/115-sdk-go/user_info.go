package sdk

import (
	"context"
	"net/http"

	"github.com/xhofe/115-sdk-go/json_types"
)

type UserInfoResp_Size struct {
	Size       json_types.Int64OrFloat `json:"size"`
	SizeFormat string                  `json:"size_format"`
}

type UserInfoResp struct {
	UserID      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	UserFaceS   string `json:"user_face_s"`
	UserFaceM   string `json:"user_face_m"`
	UserFaceL   string `json:"user_face_l"`
	RtSpaceInfo struct {
		AllTotal  UserInfoResp_Size `json:"all_total"`
		AllRemain UserInfoResp_Size `json:"all_remain"`
		AllUse    UserInfoResp_Size `json:"all_use"`
	} `json:"rt_space_info"`
	VipInfo struct {
		LevelName string `json:"level_name"`
		Expire    int64  `json:"expire"`
	} `json:"vip_info"`
}

// UserInfo: https://www.yuque.com/115yun/open/ot1litggzxa1czww
func (c *Client) UserInfo(ctx context.Context) (*UserInfoResp, error) {
	var resp UserInfoResp
	_, err := c.AuthRequest(ctx, ApiUserInfo, http.MethodGet, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}
