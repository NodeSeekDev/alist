package mopan

import (
	"strings"

	"github.com/google/uuid"
)

const (
	HomeUrl       = "https://mopan.sc.189.cn"
	MoPanProxyUrl = HomeUrl + "/mopanproxy"
	// EnterPriseUrl = HomeUrl + "/enterprise"

	MoPanProxyUpdload = MoPanProxyUrl + "/fileupload"
	MoPanProxyFamily  = MoPanProxyUrl + "/family"
	MoPanProxyAuthUrl = MoPanProxyUrl + "/auth"
)

const PublicKeyV2 = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgRb5UBbJFi3DFnxMxqqWZ0waO5a+dXLih6g47tT8H0ie+uYT3L7
nte63Cm04KX7HRovmN6zHI60m/gg6gukvqYqUtf3R2tP0i8T3KtWSZFjGqcDFLF2yvj3ntZUwU0/O3wZT3CbxOz2YoA6YXz
c0MlAjc8tu/YpBxN5CsO9auiaVSODiCNiUCFqEGBiHvQiRsX08bTOfSaTPw3SEavO24tknjAUahP/++uz2JOgLTN+zY0nmh
RZD3ArrPM84dtrldByc7g2kCwxSU3OsCpYuBZ8Po/Q/09p+Xpz2YP9dBGNnFR3sHIQcNA2Fj/nyNLRNw7FnWAcRwOvQhl8h
NqC40wIDAQAB
-----END PUBLIC KEY-----
`

// `
// -----BEGIN PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzDzuXrkkYjcb0N2/oiZXxNNaYM0PKPZE0aXeCJWo/VBSFE3+Q/1
// m5u1hI7y6U6OP8Kj5w6BStHU4lqs3/OOWG8GYKqGD6pz+b8Vp6m22lk8/7lL9375w2siAz+xSWwovAIKTfbMRwUsmJWoGI2
// vxrwok6jJoWacP6GcsI335cD7fNsHSOFYTb7SCjKWvAowsHhAWu7W8oP7bB3HE3Xth6Wy/gbZl/4Hp9rJU8w44/1Hc6O+uz
// fw4ZNtE0E4cIsK40XifW5SSokpCQkIlPNKHRuzjIuGQRbCjvl682M/DixSouc4whOcOB6Rf102p2XaKvrmmT1OXCA4dFkRa
// 5rjApQIDAQAB
// -----END PUBLIC KEY-----
// `

const PublicKeyV1 = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzDzuXrkkYjcb0N2/oiZX
xNNaYM0PKPZE0aXeCJWo/VBSFE3+Q/1m5u1hI7y6U6OP8Kj5w6BStHU4lqs3/OOW
G8GYKqGD6pz+b8Vp6m22lk8/7lL9375w2siAz+xSWwovAIKTfbMRwUsmJWoGI2vx
rwok6jJoWacP6GcsI335cD7fNsHSOFYTb7SCjKWvAowsHhAWu7W8oP7bB3HE3Xth
6Wy/gbZl/4Hp9rJU8w44/1Hc6O+uzfw4ZNtE0E4cIsK40XifW5SSokpCQkIlPNKH
RuzjIuGQRbCjvl682M/DixSouc4whOcOB6Rf102p2XaKvrmmT1OXCA4dFkRa5rjA
pQIDAQAB
-----END PUBLIC KEY-----
`

// `
// -----BEGIN PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx2bxNtzCYin2B6qWRIwPNLP6E4arzRlpPrlb2UIjWq91pJe8yN2
// QwGigxFDTfkcrv32gqlTJX2xvMr1O7RzA+oIPfA1xfJTzfQHf2HPZ8w67A4WNZHxQmWqcDdUcy6JCKzks1TLGsAH5v17dK/
// AazM2u6n5OvFrqQMnXr/raZhJRVUg3YVXW6Ppbw7fewX2n1DosC+xLU19fpyHSb/YW/9dlDMJ4tvTHrxTxpT8OOM5/bdl5q
// eUN8bBsZht1l97Iyp1Od0oFDbBaorFUsyVEnVa7r5fuFlYSoLgLiCXnMNTLpJF4GbSvEG2vXAmTLrlJ+qYWXBL7O1AJU6tZ
// KchY4wIDAQAB
// -----END PUBLIC KEY-----
// `

// 默认设备信息
// 根据windwos信息构建
var DefaultDeviceInfo = DeviceInfo{
	DeviceNo: "1104a897925070c638d",

	MpRemoteType:    "3",
	MpRemoteChannel: "100",

	MpVersion:     "1.1.202",
	MpVersionCode: 145,

	MpDeviceSerialNum: strings.ReplaceAll(uuid.NewString(), "-", ""),
	MpManufcaturer:    "Windows端",
	MpModel:           "",

	MpOs:         "windows",
	MpOsVersion:  "31",
	MpOsVersion2: "12",
}
