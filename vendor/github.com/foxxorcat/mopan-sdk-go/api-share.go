package mopan

type ShareInfoData struct {
	ID         string `json:"id"`
	AccessCode string `json:"accessCode"`
	Code       string `json:"code"`

	SendName  string `json:"sendName"`
	SendType  int    `json:"sendType"`
	Status    int    `json:"status"`
	URL       string `json:"url"`
	TotleSize Int64  `json:"totleSize"`

	UserID    string `json:"userId"`    // 分享用户ID
	ShareName string `json:"shareName"` // 分享用户名称（不完整）
	UserPic   string `json:"userPic"`

	IsPerpetual Bool `json:"isPerpetual"`
	IsDel       Bool `json:"isDel"`

	ViewCount   int `json:"viewCount"`
	DownCount   int `json:"downCount"`
	ExportCount int `json:"exportCount"`

	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

func (c *MoClient) GetShareInfo(accessCode string, option ...RestyOption) (*ShareInfoData, error) {
	var resp ShareInfoData
	_, err := c.Request(MoPanProxyFamily+"/share/info/accessCode", Json{
		"accessCode": accessCode,
	}, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *MoClient) GetShareFileDownloadUrl(shareId, fileId string, option ...RestyOption) (*GetFileDownloadUrlData, error) {
	var resp GetFileDownloadUrlData
	_, err := c.Request(MoPanProxyFamily+"", Json{
		"shareId": shareId,
		"fileId":  fileId,
	}, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
