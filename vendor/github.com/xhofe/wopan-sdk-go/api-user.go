package wopan

// PcWebLoginData no encrypt
type PcWebLoginData struct {
	NeedSmsCode string `json:"needSmsCode"`
}

func (w *WoClient) PcWebLogin(phone, password string, opts ...RestyOption) (*PcWebLoginData, error) {
	var resp PcWebLoginData
	_, err := w.RequestApiUser(KeyPcWebLogin, Json{
		"phone":        phone,
		"password":     password,
		"uuid":         "",
		"verifyCode":   "",
		"clientSecret": DefaultClientSecret,
	}, JsonClientIDSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// PcLoginVerifyCodeData no encrypt
type PcLoginVerifyCodeData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func (w *WoClient) PcLoginVerifyCode(phone, password, messageCode string, opts ...RestyOption) (*PcLoginVerifyCodeData, error) {
	var resp PcLoginVerifyCodeData
	_, err := w.RequestApiUser(KeyPcLoginVerifyCode, Json{
		"phone":        phone,
		"messageCode":  messageCode,
		"verifyCode":   nil,
		"uuid":         nil,
		"clientSecret": DefaultClientSecret,
		"password":     password,
	}, JsonClientIDSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AppQueryUserData struct {
	UserId        string `json:"userId"`
	HeadUrl       string `json:"headUrl"`
	UserName      string `json:"userName"`
	Sex           string `json:"sex"`
	Birthday      string `json:"birthday"`
	IsModify      string `json:"isModify"`
	IsHeadModify  string `json:"isHeadModify"`
	IsSetPassword string `json:"isSetPassword"`
	RegisterTime  string `json:"registerTime"`
}

func (w *WoClient) AppQueryUser(opts ...RestyOption) (*AppQueryUserData, error) {
	var resp AppQueryUserData
	_, err := w.RequestApiUser(KeyAppQueryUser, Json{
		"accessToken": w.accessToken,
	}, JsonClientIDSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppRefreshTokenData no encrypt
type AppRefreshTokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (w *WoClient) AppRefreshToken(opts ...RestyOption) (*AppRefreshTokenData, error) {
	var resp AppRefreshTokenData
	_, err := w.RequestApiUser(KeyAppRefreshToken, Json{
		"refreshToken": w.refreshToken,
		"clientSecret": DefaultClientSecret,
	}, JsonClientIDSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (w *WoClient) AppLogout(opts ...RestyOption) error {
	_, err := w.RequestApiUser(KeyAppLogout, Json{
		"accessToken": w.accessToken,
	}, JsonClientIDSecret, nil, opts...)
	return err
}
