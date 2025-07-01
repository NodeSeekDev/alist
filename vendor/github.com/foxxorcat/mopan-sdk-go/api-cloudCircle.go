package mopan

// 共享空间
type CloudCircleInfo struct {
	AlbumFolderID string `json:"albumFolderId"` // 相册ID
	AlbumName     string `json:"albumName"`     // 相册名称
	AlbumPath     string `json:"albumPath"`     // 相册路径

	CloudCircleName  string `json:"cloudCircleName"` // 空间名称
	CloudID          string `json:"cloudId"`         // 空间ID
	CloudRelationID  string `json:"cloudRelationId"` // 关联ID
	CreateUser       string `json:"createUser"`      // 创建用户
	DynamicCoverType int    `json:"dynamicCoverType"`

	FileFolderID string `json:"fileFolderId"` // 目录ID
	FileName     string `json:"fileName"`     // 文件名称
	FilePath     string `json:"filePath"`     // 文件路径

	IsDisturb Bool `json:"isDisturb"` // 消息免打扰
	IsTop     Bool `json:"isTop"`     // 是否置顶
	IsTv      Bool `json:"isTv"`      // 是否在TV上展示

	Aliss         string `json:"aliss"`        // 用户昵称
	JoinOrCreate  int    `json:"joinOrCreate"` // (1:creare,2:join)
	JoinType      int    `json:"joinType"`     // 加入权限(1:允许任何人加入,2:管理员审核加入,3:不允许任何人加入)
	PublishFlag   int    `json:"publishFlag"`  // 发布权限(1:仅限管理员,2:所有成员)
	MemberFlag    int    `json:"memberFlag"`
	OperationFlag int    `json:"operationFlag"`
	PackageFlag   int    `json:"packageFlag"`

	UserList []CloudCircleUser `json:"userList"`          // 用户列表
	TopTime  Time3             `json:"topTime,omitempty"` // 置顶操作时间
}

type CloudCircleUser struct {
	UserID           string `json:"userId"`           // 用户ID
	Phone            string `json:"phone"`            // 手机号(完整)
	Aliss            string `json:"aliss"`            // 昵称
	CircleRelationID string `json:"circleRelationId"` // 共享云ID
	Type             int    `json:"type"`
}

// 查询所有共享空间信息
func (c *MoClient) QueryAllCloudCircleApp(option ...RestyOption) ([]CloudCircleInfo, error) {
	param := Json{
		"type": 2,
	}

	var resp []CloudCircleInfo
	_, err := c.Request(MoPanProxyFamily+"/cloudCircle/queryAllCloudCircleApp", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 查询单个空间信息
func (c *MoClient) CloudCircleUserDetail(cloudId string, option ...RestyOption) (*CloudCircleInfo, error) {
	param := Json{
		"cloudId": cloudId,
		"type":    2,
	}

	var resp CloudCircleInfo
	_, err := c.Request(MoPanProxyFamily+"/cloudCircle/cloudCircleUserDetail", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *MoClient) GetCloudToken(cloudID string, option ...RestyOption) (*TokenData, error) {
	param := Json{
		"cloudId": cloudID,
	}

	var resp TokenData
	_, err := c.Request(MoPanProxyFamily+"/user/getTokenByCloudId", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 创建空间
func (c *MoClient) CreateCloudCircle(cloudName string, option ...RestyOption) (*CloudCircleInfo, error) {
	var resp CloudCircleInfo
	_, err := c.Request(MoPanProxyFamily+"/cloudCircle/createCloudCircle", Json{
		"cloudName": cloudName,
	}, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 解散空间
func (c *MoClient) DismissCloudCircle(cloudID string, option ...RestyOption) error {
	param := Json{
		"cloudId": cloudID,
	}
	_, err := c.Request(MoPanProxyFamily+"/cloudCircle/dismissCloudCircle", param, nil, option...)
	return err
}

// 加入空间
func (c *MoClient) InviteCloudCircleUser(code string, option ...RestyOption) error {
	_, err := c.Request(MoPanProxyFamily+"/cloudCircle/inviteCloudCircleUser", Json{
		"code": code,
	}, nil, option...)
	if err != nil {
		return err
	}
	return nil
}

// 退出空间
func (c *MoClient) CloudCircleUserQuit(cloudId string, option ...RestyOption) error {
	_, err := c.Request(MoPanProxyFamily+"/cloudCircle/cloudCircleUserQuit", Json{
		"cloudId": cloudId,
	}, nil, option...)
	if err != nil {
		return err
	}
	return nil
}
