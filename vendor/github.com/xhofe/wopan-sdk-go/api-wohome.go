package wopan

type FCloudProductOrdListQryCtxData struct {
	FcloudProductOrds []struct {
		ActiviteCode        string `json:"activiteCode"`
		AppStorePackageDesc string `json:"appStorePackageDesc"`
		AppStorePackageFee  string `json:"appStorePackageFee"`
		AppStoreProductId   string `json:"appStoreProductId"`
		ApplyTime           string `json:"applyTime"`
		ApplyTimeFormate    string `json:"applyTimeFormate"`
		CbssOrderId         string `json:"cbssOrderId"`
		City                string `json:"city"`
		ClientId            string `json:"clientId"`
		Days                string `json:"days"`
		DescUrl             string `json:"descUrl"`
		EffectState         string `json:"effectState"`
		EffectiveDays       int    `json:"effectiveDays"`
		ExpireTime          string `json:"expireTime"`
		ExpireTimeFormate   string `json:"expireTimeFormate"`
		Fee                 string `json:"fee"`
		IsAppStorePay       string `json:"isAppStorePay"`
		IsAutoSub           string `json:"isAutoSub"`
		IsExpire            string `json:"isExpire"`
		IsNewPackage        string `json:"isNewPackage"`
		IsOnline            string `json:"isOnline"`
		IsPlus              string `json:"isPlus"`
		IsShowExpireTips    string `json:"isShowExpireTips"`
		OrderId             string `json:"orderId"`
		OrderState          string `json:"orderState"`
		OrderStatus         string `json:"orderStatus"`
		PackageDesc         string `json:"packageDesc"`
		PackageProductCode  string `json:"packageProductCode"`
		PackageProductId    string `json:"packageProductId"`
		PayMethod           string `json:"payMethod"`
		PayTransactionId    string `json:"payTransactionId"`
		PayType             string `json:"payType"`
		Province            string `json:"province"`
		RemainDays          string `json:"remainDays"`
		SignStatus          string `json:"signStatus"`
		Source              string `json:"source"`
		SubTime             string `json:"subTime"`
		SubTimeFormate      string `json:"subTimeFormate"`
		SubType             string `json:"subType"`
		UserId              string `json:"userId"`
		VipDesc             string `json:"vipDesc"`
		VipDescNew          string `json:"vipDescNew"`
		VipExpireTimeLabel  string `json:"vipExpireTimeLabel"`
		VipLevel            string `json:"vipLevel"`
	} `json:"fcloudProductOrds"`
	MaxVipLevel string `json:"maxVipLevel"`
	IsShowInlet string `json:"isShowInlet"`
}

func (w *WoClient) FCloudProductOrdListQry(opts ...RestyOption) (*FCloudProductOrdListQryCtxData, error) {
	var resp FCloudProductOrdListQryCtxData
	_, err := w.RequestWoHome(KeyFCloudProductOrdListQry, nil, Json{
		"qryType":  "1",
		"clientId": DefaultClientID,
	}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type QueryCloudUsageInfoData struct {
	Code      string `json:"code"`
	UsageInfo struct {
		TotalSize     string `json:"totalSize"`
		UsedSize      int64  `json:"usedSize"`
		ImageSize     int64  `json:"imageSize"`
		VideoSize     int64  `json:"videoSize"`
		AudioSize     int64  `json:"audioSize"`
		TextSize      int64  `json:"textSize"`
		OtherSize     int64  `json:"otherSize"`
		ByteUsedSize  int64  `json:"byteUsedSize"`
		ByteTotalSize string `json:"byteTotalSize"`
	} `json:"usageInfo"`
	VipLevel   string `json:"vipLevel"`
	ExpireTime string `json:"expireTime"`
	ApplyTime  string `json:"applyTime"`
	PayType    string `json:"payType"`
	Source     string `json:"source"`
	OrderState string `json:"orderState"`
	Status     string `json:"status"`
}

func (w *WoClient) QueryCloudUsageInfo(opts ...RestyOption) (*QueryCloudUsageInfoData, error) {
	err := w.InitPhone()
	if err != nil {
		return nil, err
	}
	var resp QueryCloudUsageInfoData
	_, err = w.RequestWoHome(KeyQueryCloudUsageInfo, Json{
		"phoneNum": w.Phone,
		"clientId": DefaultClientID,
	}, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FCloudProductPackage is not required

type ClassifyRuleData struct {
	FileTypes map[string]struct {
		SubType string `json:"subType"`
		Ability string `json:"ability"`
		Type    string `json:"type"`
	} `json:"fileTypes"`
	FileIcons struct {
		App struct {
			Zip   string `json:"zip"`
			Image string `json:"image"`
			Other string `json:"other"`
			Xmind string `json:"xmind"`
			Gif   string `json:"gif"`
			Csv   string `json:"csv"`
			Video string `json:"video"`
			Excel string `json:"excel"`
			Pdf   string `json:"pdf"`
			Ppt   string `json:"ppt"`
			Doc   string `json:"doc"`
			Audio string `json:"audio"`
			Text  string `json:"text"`
			Word  string `json:"word"`
		} `json:"app"`
		H5 struct {
		} `json:"h5"`
	} `json:"fileIcons"`
}

func (w *WoClient) ClassifyRule(opts ...RestyOption) (*ClassifyRuleData, error) {
	var resp ClassifyRuleData
	_, err := w.RequestWoHome(KeyClassifyRule, Json{}, Json{
		"key": true,
	}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetZoneInfoData struct {
	Url string `json:"url"`
}

func (w *WoClient) GetZoneInfo(opts ...RestyOption) (*GetZoneInfoData, error) {
	var resp GetZoneInfoData
	_, err := w.RequestWoHome(KeyGetZoneInfo, Json{
		"appId": DefaultAppID,
	}, Json{
		"key": true,
	}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type FamilyUserCurrentEncodeData struct {
	Count           string `json:"count"`
	DefaultHomeId   int    `json:"defaultHomeId"`
	DefaultHomeName string `json:"defaultHomeName"`
	GroupHeadUrl    string `json:"groupHeadUrl"`
	GroupName       string `json:"groupName"`
	Id              int    `json:"id"`
	MemberRole      string `json:"memberRole"`
	OwnerId         string `json:"owner  Id"`
	UnreadFlag      string `json:"unreadFlag"`
}

func (w *WoClient) FamilyUserCurrentEncode(opts ...RestyOption) (*FamilyUserCurrentEncodeData, error) {
	var resp FamilyUserCurrentEncodeData
	_, err := w.RequestWoHome(KeyFamilyUserCurrentEncode, Json{
		"clientId": DefaultClientID,
	}, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
