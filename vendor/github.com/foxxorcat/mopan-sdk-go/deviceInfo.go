package mopan

import (
	"encoding/base64"
	"encoding/json"
)

type DeviceInfo struct {
	DeviceNo string `json:"deviceNo"` //（PC：1104a897925070c638d，Android：140fe1da9fbf4b719f1）

	MpRemoteType    string `json:"mp_remoteType"`    // android:1,windows:3,mac:5
	MpRemoteChannel string `json:"mp_remoteChannel"` // android:203,PC:100

	MpVersion     string `json:"mp_version"`      // 程序版本
	MpVersionCode int    `json:"mp_version_code"` // 程序版本Code

	MpDeviceSerialNum string `json:"mp_deviceSerialNum"` // 设备编号 uuid len:32
	MpManufcaturer    string `json:"mp_manufcaturer"`    // (Windows端,Mac端,{设备品牌})
	MpModel           string `json:"mp_model"`           // 手机型号

	MpOs         string `json:"mp_os"`         // (Android,Windows,macOS)
	MpOsVersion  string `json:"mp_osVersion"`  // 31
	MpOsVersion2 string `json:"mp_osVersion2"` // 12
}

func (e *DeviceInfo) Encrypt(key string) string {
	data, _ := json.Marshal(e)
	enc, _ := AesEncrypt(data, []byte(key))
	return base64.StdEncoding.EncodeToString(enc)

}

func (e *DeviceInfo) Decrypt(v string, key string) error {
	enc, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err
	}
	dec, err := AesDecrypt(enc, []byte(key))
	if err != nil {
		return err
	}
	return json.Unmarshal(dec, e)
}
