package wopan

func (w *WoClient) InitPhone() error {
	if w.Phone == "" {
		_resp, err := w.AppQueryUser()
		if err != nil {
			return err
		}
		w.Phone = _resp.UserId
	}
	return nil
}

func (w *WoClient) InitClassifyRule() error {
	if w.ClassifyRuleData != nil {
		return nil
	}
	data, err := w.ClassifyRule()
	if err != nil {
		return err
	}
	w.ClassifyRuleData = data
	return nil
}

func (w *WoClient) InitZoneURL() error {
	if w.ZoneURL != "" {
		return nil
	}
	data, err := w.GetZoneInfo()
	if err != nil {
		return err
	}
	w.ZoneURL = data.Url
	return nil
}

func (w *WoClient) RefreshToken() error {
	resp, err := w.AppRefreshToken()
	if err != nil {
		return err
	}
	w.onRefreshToken(resp.AccessToken, resp.RefreshToken)
	w.SetAccessToken(resp.AccessToken)
	w.SetRefreshToken(resp.RefreshToken)
	return nil
}

func (w *WoClient) InitData() error {
	if w.accessToken == "" && w.refreshToken != "" {
		if err := w.RefreshToken(); err != nil {
			return err
		}
	}
	if err := w.InitPhone(); err != nil {
		return err
	}
	if err := w.InitClassifyRule(); err != nil {
		return err
	}
	if err := w.InitZoneURL(); err != nil {
		return err
	}
	return nil
}
