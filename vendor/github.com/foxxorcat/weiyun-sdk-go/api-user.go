package weiyunsdkgo

type DiskUserInfoGetData struct {
	Uin        int64     `json:"uin"`
	NickName   string    `json:"nick_name"`
	HeadImgUrl string    `json:"head_img_url"`
	UserCtime  TimeStamp `json:"user_ctime"`
	UserMtime  TimeStamp `json:"user_mtime"`

	UsedSpace  int64 `json:"used_space"`
	TotalSpace int64 `json:"total_space"`

	DirTotal  int64 `json:"dir_total"`
	FileTotal int64 `json:"file_total"`

	RootDirKey string `json:"root_dir_key"`
	MainDirKey string `json:"main_dir_key"`
}

// 获取用户信息
func (c *WeiYunClient) DiskUserInfoGet(opts ...RestyOption) (*DiskUserInfoGetData, error) {
	param := Json{
		"is_get_upload_flow_flag":     false,
		"is_get_high_speed_flow_info": false,
		"is_get_weiyun_flag":          false,
		"is_get_space_clean_info":     false,
	}
	var resp DiskUserInfoGetData
	_, err := c.WeiyunQdiskClientRequest("DiskUserInfoGet", 2201, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
