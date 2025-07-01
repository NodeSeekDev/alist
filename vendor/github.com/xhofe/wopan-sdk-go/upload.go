package wopan

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Upload2COption struct {
	OnProgress func(current, total int64)
	Ctx        context.Context
}

type Upload2CFile struct {
	Name        string
	Size        int64
	Content     io.Reader
	ContentType string
}

type Upload2CResp struct {
	Code string `json:"code"`
	Data struct {
		Fid string `json:"fid"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func (w *WoClient) Upload2C(spaceType string, file Upload2CFile, targetDirId string, familyId string, opt Upload2COption) (string, error) {
	zoneURL := DefaultZoneURL
	w.zoneURLOnce.Do(func() {
		_ = w.InitZoneURL()
	})
	if w.ZoneURL != "" {
		zoneURL = w.ZoneURL
	}
	batchNo := time.Now().Format("20060102150405")
	fileInfo := Json{
		"spaceType":   spaceType,
		"directoryId": targetDirId,
		"batchNo":     batchNo,
		"fileName":    file.Name,
		"fileSize":    file.Size,
		"fileType":    w.GetFileType(file.Name),
	}
	if spaceType == SpaceTypeFamily {
		fileInfo["familyId"] = familyId
	}
	fileInfoStr, err := w.EncryptParam(ChannelWoHome, fileInfo)
	if err != nil {
		return "", err
	}
	uploadURL := zoneURL + "/openapi/client/" + KeyUpload2C
	totalPart := file.Size / DefaultPartSize
	if totalPart == 0 {
		totalPart = 1
	}
	formData := map[string]string{
		"uniqueId":    strconv.FormatInt(time.Now().UnixMilli(), 10),
		"accessToken": w.accessToken,
		"fileName":    file.Name,
		"psToken":     "undefined",
		"fileSize":    strconv.FormatInt(file.Size, 10),
		"totalPart":   strconv.FormatInt(totalPart, 10),
		"channel":     ChannelWoCloud,
		"directoryId": targetDirId,
		"fileInfo":    fileInfoStr,
		//partSize: 8388608
		//partIndex: 7
		//file: (binary)
	}
	var resp Upload2CResp
	var fid string
	var finishedSize int64 = 0
	for partIndex := int64(1); partIndex <= totalPart; partIndex++ {
		if opt.Ctx != nil {
			select {
			case <-opt.Ctx.Done():
				return "", opt.Ctx.Err()
			default:
			}
		}
		partSize := DefaultPartSize
		if partIndex == totalPart {
			partSize = file.Size - finishedSize
		}
		formData["partSize"] = strconv.FormatInt(partSize, 10)
		formData["partIndex"] = strconv.FormatInt(partIndex, 10)

		req := w.NewRequest().
			SetResult(&resp).
			ForceContentType("application/json;charset=UTF-8").
			SetHeaders(map[string]string{
				"Origin":     "https://pan.wo.cn",
				"Referer":    "https://pan.wo.cn/",
				"User-Agent": w.ua,
			}).
			SetMultipartFormData(formData).
			SetMultipartField("file", file.Name, file.ContentType, io.LimitReader(file.Content, partSize))
		if opt.Ctx != nil {
			req.SetContext(opt.Ctx)
		}
		res, err := req.Post(uploadURL)
		if err != nil {
			return "", err
		}
		if res.IsError() {
			return "", fmt.Errorf("partIndex: %d, failed to upload2C with http status: %d, body: %s", partIndex, res.StatusCode(), res.String())
		}
		if resp.Code != "0000" {
			return "", fmt.Errorf("partIndex: %d, failed to upload2C with code: %s, msg: %s", partIndex, resp.Code, resp.Msg)
		}
		if resp.Data.Fid != "" {
			fid = resp.Data.Fid
		}
		finishedSize += partSize
		if opt.OnProgress != nil {
			opt.OnProgress(finishedSize, file.Size)
		}
	}
	return fid, nil
}

func (w *WoClient) Upload2CPersonal(file Upload2CFile, targetDirId string, opt Upload2COption) (string, error) {
	return w.Upload2C(SpaceTypePersonal, file, targetDirId, "", opt)
}

func (w *WoClient) Upload2CFamily(file Upload2CFile, targetDirId string, familyId string, opt Upload2COption) (string, error) {
	return w.Upload2C(SpaceTypeFamily, file, targetDirId, familyId, opt)
}
