package weiyunsdkgo

type SafeBoxCheckStatusData struct {
	SafeBoxPdirkey string `json:"safe_box_pdirkey"`
	SafeBoxDirkey  string `json:"safe_box_dirkey"`
	PwdFailLimit   any    `json:"pwd_fail_limit"`
	PwdRetryCnt    int    `json:"pwd_retry_cnt"`
	NeedNotice     any    `json:"need_notice"`
}

func (c *WeiYunClient) SafeBoxCheckStatus(opts ...RestyOption) (*SafeBoxCheckStatusData, error) {
	var resp SafeBoxCheckStatusData
	_, err := c.WeiyunSafeBoxRequest("SafeBoxCheckStatus", 28406, nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SafeBoxPublicKeyGetData struct {
	PublicKey string `json:"public_key"`
}

func (c *WeiYunClient) SafeBoxPublicKeyGet(opts ...RestyOption) (*SafeBoxPublicKeyGetData, error) {
	var resp SafeBoxPublicKeyGetData
	_, err := c.WeiyunSafeBoxRequest("SafeBoxPublicKeyGet", 28408, nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return nil, err
}
