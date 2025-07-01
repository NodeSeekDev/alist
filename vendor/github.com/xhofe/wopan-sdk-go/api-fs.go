package wopan

type File struct {
	FamilyId     int    `json:"familyId"`
	Fid          string `json:"fid"`
	Creator      string `json:"creator"`
	Size         int64  `json:"size"`
	CreateTime   string `json:"createTime"`
	Name         string `json:"name"`
	ShootingTime string `json:"shootingTime"`
	Id           string `json:"id"`
	Type         int    `json:"type"`
	ThumbUrl     string `json:"thumbUrl"`
	FileType     string `json:"fileType"`
}

type QueryAllFilesData struct {
	Files []File `json:"files"`
}

func (w *WoClient) QueryAllFiles(spaceType, parentDirectoryId string, pageNum, pageSize int, sortRule int, familyId string, opts ...RestyOption) (*QueryAllFilesData, error) {
	var resp QueryAllFilesData
	param := Json{
		"spaceType":         spaceType,
		"parentDirectoryId": parentDirectoryId,
		"pageNum":           pageNum,
		"pageSize":          pageSize,
		"sortRule":          sortRule,
		"clientId":          DefaultClientID,
	}
	if spaceType == SpaceTypeFamily {
		param["familyId"] = familyId
	}
	_, err := w.RequestWoHome(KeyQueryAllFiles, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (w *WoClient) QueryAllFilesPersonal(parentDirectoryId string, pageNum, pageSize int, sortRule int, opts ...RestyOption) (*QueryAllFilesData, error) {
	return w.QueryAllFiles(SpaceTypePersonal, parentDirectoryId, pageNum, pageSize, sortRule, "", opts...)
}

func (w *WoClient) QueryAllFilesFamily(parentDirectoryId string, pageNum, pageSize int, sortRule int, familyId string, opts ...RestyOption) (*QueryAllFilesData, error) {
	return w.QueryAllFiles(SpaceTypeFamily, parentDirectoryId, pageNum, pageSize, sortRule, familyId, opts...)
}

// GetSearchDirectory??

type GetDownloadUrlV2Data struct {
	Type int `json:"type"`
	List []struct {
		Fid         string `json:"fid"`
		DownloadUrl string `json:"downloadUrl"`
	} `json:"list"`
}

func (w *WoClient) GetDownloadUrlV2(fidList []string, opts ...RestyOption) (*GetDownloadUrlV2Data, error) {
	var resp GetDownloadUrlV2Data
	param := Json{
		"type":     "1",
		"fidList":  fidList,
		"clientId": DefaultClientID,
	}
	_, err := w.RequestWoHome(KeyGetDownloadUrlV2, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetDownloadUrlData struct {
	Fid         string `json:"fid"`
	DownloadUrl string `json:"downloadUrl"`
}

func (w *WoClient) GetDownloadUrl(spaceType string, fidList []string, opts ...RestyOption) ([]GetDownloadUrlData, error) {
	var resp []GetDownloadUrlData
	param := Json{
		"fidList":   fidList,
		"clientId":  DefaultClientID, // ???
		"spaceType": spaceType,
	}
	_, err := w.RequestWoHome(KeyGetDownloadUrl, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type CreateDirectoryData struct {
	Id string `json:"id"`
}

func (w *WoClient) CreateDirectory(spaceType, parentDirectoryId string, directoryName, familyId string, opts ...RestyOption) (*CreateDirectoryData, error) {
	var resp CreateDirectoryData
	param := Json{
		"spaceType":         spaceType,
		"familyId":          familyId,
		"parentDirectoryId": parentDirectoryId,
		"directoryName":     directoryName,
		"clientId":          DefaultClientID,
	}
	_, err := w.RequestWoHome(KeyCreateDirectory, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// RenameFileOrDirectory
// _type: 1: file, 0: directory
func (w *WoClient) RenameFileOrDirectory(spaceType string, _type int, id string, name string, familyId string, opts ...RestyOption) error {
	fileType := "0"
	if _type != 0 {
		fileType = w.GetFileType(name)
	}
	param := Json{
		"spaceType": spaceType,
		"type":      _type,
		"fileType":  fileType,
		"id":        id,
		"name":      name,
		"clientId":  DefaultClientID,
	}
	if spaceType == SpaceTypeFamily {
		param["familyId"] = familyId
	}
	_, err := w.RequestWoHome(KeyRenameFileOrDirectory, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) RenameFileOrDirectoryPersonal(_type int, id string, name string, opts ...RestyOption) error {
	return w.RenameFileOrDirectory(SpaceTypePersonal, _type, id, name, "", opts...)
}

func (w *WoClient) RenameFileOrDirectoryFamily(_type int, id string, name string, familyId string, opts ...RestyOption) error {
	return w.RenameFileOrDirectory(SpaceTypeFamily, _type, id, name, familyId, opts...)
}

func (w *WoClient) MoveFile(dirList, fileList []string, targetDirId string, sourceType, targetType string, fromFamilyId, targetFamilyId string, opts ...RestyOption) error {
	param := Json{
		"targetDirId": targetDirId,
		"sourceType":  sourceType,
		"targetType":  targetType,
		"dirList":     dirList,
		"fileList":    fileList,
		"secret":      false,
		"clientId":    DefaultClientID,
	}
	if sourceType == SpaceTypeFamily {
		param["fromFamilyId"] = fromFamilyId
	}
	if targetType == SpaceTypeFamily {
		param["familyId"] = targetFamilyId
	}
	_, err := w.RequestWoHome(KeyMoveFile, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) CopyFile(dirList, fileList []string, targetDirId string, sourceType, targetType string, fromFamilyId, targetFamilyId string, opts ...RestyOption) error {
	param := Json{
		"targetDirId": targetDirId,
		"sourceType":  sourceType,
		"targetType":  targetType,
		"dirList":     dirList,
		"fileList":    fileList,
		"secret":      false,
		"clientId":    DefaultClientID,
	}
	if sourceType == SpaceTypeFamily {
		param["fromFamilyId"] = fromFamilyId
	}
	if targetType == SpaceTypeFamily {
		param["familyId"] = targetFamilyId
	}
	_, err := w.RequestWoHome(KeyCopyFile, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) DeleteFile(spaceType string, dirList, fileList []string, opts ...RestyOption) error {
	param := Json{
		"spaceType": spaceType,
		"vipLevel":  "0",
		"dirList":   dirList,
		"fileList":  fileList,
		"clientId":  DefaultClientID,
	}
	_, err := w.RequestWoHome(KeyDeleteFile, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) EmptyRecycleData(opts ...RestyOption) error {
	param := Json{
		"clientId": DefaultClientID,
	}
	_, err := w.RequestWoHome(KeyEmptyRecycleData, param, JsonSecret, nil, opts...)
	return err
}
