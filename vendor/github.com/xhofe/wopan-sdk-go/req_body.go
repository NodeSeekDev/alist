package wopan

func (w *WoClient) EncryptParam(channel string, param Json) (string, error) {
	if param == nil {
		return "", nil
	}
	jsonBytes, err := w.jsonMarshalFunc(param)
	if err != nil {
		return "", err
	}
	encrypted, err := w.crypto.EncryptBytes(jsonBytes, channel)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func (w *WoClient) NewBody(channel string, param, other Json) (Json, error) {
	if param == nil {
		return other, nil
	}
	encrypted, err := w.EncryptParam(channel, param)
	if err != nil {
		return nil, err
	}
	// copy other to avoid modifying the original and concurrent map writes
	_other := copyJson(other)
	_other["param"] = encrypted
	return _other, nil
}
