package mopan

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 用户基本信息
// type UserKeyInfo struct {
// 	ID         string `json:"id"`     // 数据库ID
// 	UserID     string `json:"userId"` // 用户ID
// 	Passwd     string `json:"passwd"` // 用户密码(完整) !!
// 	FailTime   string `json:"failTime"`
// 	CreateTime string `json:"createTime"` // 创建时间
// 	UpdateTime string `json:"updateTime"` // 更新时间
// 	IsDel      Bool   `json:"isDel"`      // 软删除
// }

// 用户详细信息
type UserInfo struct {
	UserID   string `json:"userId"`   // 用户ID
	KickName string `json:"kickName"` // 昵称
	Name     string `json:"name"`     // 手机号(省略)
	Phone    string `json:"phone"`    // 手机号(完整)

	RemoteChannel string `json:"remoteChannel"` // 与 DeviceInfo 相同

	IsTag   Bool `json:"isTag"`
	Isp     Bool `json:"isp"`
	Status  int  `json:"status"`
	Version int  `json:"version"`

	ArrearsFlag Bool `json:"arrearsFlag"` // 欠费标识
	FreezeFlag  int  `json:"freezeFlag"`  // 冻结标识
	PackageFlag int  `json:"packageFlag"`
	PrimaryFlag int  `json:"primaryFlag"`
	VipFlag     int  `json:"vipFlag"` // vip标识
	VirtualFlag int  `json:"virtualFlag"`

	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
	IsDel      Bool   `json:"isDel"`      // 软删除
}

// 获取用户详细信息
func (c *MoClient) GetUserInfo(option ...RestyOption) (*UserInfo, error) {
	var resp UserInfo
	_, err := c.Request(MoPanProxyFamily+"/user/info/getUserInfo", Json{}, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type (
	TokenData struct {
		AccessToken string `json:"accessToken"` // 可能用于文档预览
	}

	LoginResp struct {
		Token string `json:"token"`

		ExpiresIn string `json:"expiresIn"`

		AutoToken TokenData `json:"autoToken"`

		UserID   string `json:"userId"`   // 用户ID
		NickName string `json:"nickName"` // 昵称
		// UserInfo    UserInfo    `json:"userInfo"`    // 用户信息
		// UserKeyInfo UserKeyInfo `json:"userKeyInfo"` // 用户关键信息
		// HomeSubject any         `json:"homeSubject"`

		UserCloudStorageRelations []UserCloudStorageRelation `json:"userCloudStorageRelations"` // 用户空间信息
	}

	UserCloudStorageRelation struct {
		CatalogueType int    `json:"catalogueType"`
		CloudType     int    `json:"cloudType"`
		CreateTime    string `json:"createTime"`
		FileOrAlbum   int    `json:"fileOrAlbum"`
		FolderID      string `json:"folderId"`
		ID            string `json:"id"`
		IsDel         int    `json:"isDel"`
		Name          string `json:"name"`
		Path          string `json:"path"`
		UpdateTime    string `json:"updateTime"`
		UserID        string `json:"userId"`
	}
)

// 通过手机号密码登录
func (c *MoClient) Login(phone, password string, option ...RestyOption) (*LoginResp, error) {
	timestamp := time.Now().UnixMilli()

	random := strconv.FormatFloat(rand.Float64(), 'f', 20, 32)
	sign := Md5Hex(fmt.Sprintf("%s%s%d%s", phone, password, timestamp, random))
	var resp LoginResp
	_, err := c.Request(MoPanProxyAuthUrl+"/through/login/pwd", Json{
		"phone":    phone,
		"password": password,
		"random":   random,
		"time":     timestamp,
		"sign":     sign,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 短信登陆
func (c *MoClient) LoginBySms(phone string, option ...RestyOption) (func(code string) (*LoginResp, error), error) {
	_, err := c.Request(MoPanProxyAuthUrl+"/through/auth/getSmsCode", Json{
		"phone": phone,
	}, nil, option...)
	if err != nil {
		return nil, err
	}
	return func(smsCode string) (*LoginResp, error) {
		return c.LoginBySmsStep2(phone, smsCode, option...)
	}, nil
}

// 短信登陆步骤2
func (c *MoClient) LoginBySmsStep2(phone, smsCode string, option ...RestyOption) (*LoginResp, error) {
	var resp LoginResp
	_, err := c.request(MoPanProxyAuthUrl+"/through/login/sms", Json{
		"phone":   phone,
		"smsCode": smsCode,
	}, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 忘记密码
func (c *MoClient) ForgotPassword(phone string, option ...RestyOption) (func(code, newPassword string) error, error) {
	_, err := c.Request(MoPanProxyFamily+"/user/sendUpdateUserPasswordSMS", Json{
		"addressee":  phone,
		"isEmail":    "",
		"templateID": "557917",
	}, nil, option...)
	if err != nil {
		return nil, err
	}
	return func(code, newPassword string) error {
		return c.ForgotPasswordStep2(phone, code, newPassword, option...)
	}, nil
}

// 忘记密码步骤2
func (c *MoClient) ForgotPasswordStep2(phone, code, newPassword string, option ...RestyOption) error {
	_, err := c.Request(MoPanProxyFamily+"/user/updateUserPasswordCheckCode", Json{
		"code":      code,
		"phone":     phone,
		"newPasswd": newPassword,
		"email":     phone + "@189.cn",
	}, nil, option...)
	return err
}

type OnlineData struct {
	// 登录信息

	DeviceNo     string `json:"deviceNo"` // 设备编号
	IP           string `json:"ip"`
	LoginAccount string `json:"loginAccount"`
	LoginDevice  string `json:"loginDevice"`
	LoginOs      string `json:"loginOs"`
	LoginPlat    string `json:"loginPlat"`
	LoginTime    Time3  `json:"loginTime"`

	Random string `json:"random"`

	// 登录用户信息

	UserID      string `json:"userId"`
	Username    string `json:"username"`
	Phone       string `json:"phone"`
	NickName    string `json:"nickName"`
	Authorities []any  `json:"authorities"`
	Roles       []any  `json:"roles"`

	AccountNonExpired     bool  `json:"accountNonExpired"`     // 帐户未过期
	AccountNonLocked      bool  `json:"accountNonLocked"`      // 帐户未锁定
	CredentialsNonExpired bool  `json:"credentialsNonExpired"` // 证书未过期
	ExpireTime            Time3 `json:"expireTime"`            // 过期时间

	Token   string `json:"token,omitempty"` // 当前设备返回Token
	Enabled bool   `json:"enabled"`
}

// 获取所有在线设备
func (c *MoClient) Online(option ...RestyOption) ([]OnlineData, error) {
	var resp []OnlineData
	_, err := c.Request(MoPanProxyAuthUrl+"/api/userAccount/online", nil, &resp, option...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 登出在线设备(错误的参数可能导致当前账号登出)
func (c *MoClient) LogoutOnline(random, remoteType string, option ...RestyOption) error {
	param := Json{
		"random":     random,
		"remoteType": remoteType,
	}
	_, err := c.Request(MoPanProxyAuthUrl+"/api/userAccount/logout", param, nil, option...)
	if err != nil {
		return err
	}
	return nil
}

// 登出当前账号
func (c *MoClient) Logout(option ...RestyOption) error {
	_, err := c.Request(MoPanProxyAuthUrl+"/api/userAccount/logout", Json{}, nil, option...)
	if err != nil {
		return err
	}
	return nil
}

type UsedSpaceData struct {
	Capacity          Int64 `json:"capacity"`
	Used              Int64 `json:"used"`
	UserCloudSizeList []struct {
		UserID string `json:"userId"`
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Type   int    `json:"type"`
		Used   Int64  `json:"used"`
	} `json:"userCloudSizeList"`
}

// 使用空间查询
func (c *MoClient) UsedSpace(option ...RestyOption) (*UsedSpaceData, error) {
	var resp UsedSpaceData
	_, err := c.Request(MoPanProxyFamily+"/my/usedSpace", nil, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
