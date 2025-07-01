package mopan

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetSecretKey() string {
	return uuid.NewString()[:16]
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// AesEncrypt 加密
func AesEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptBytes := pkcs7Padding(data, block.BlockSize())
	encrypted := make([]byte, len(encryptBytes))
	NewECBEncrypter(block).CryptBlocks(encrypted, encryptBytes)
	return encrypted, nil
}

// AesDecrypt 解密
func AesDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	NewECBDecrypter(block).CryptBlocks(decrypted, data)
	return pkcs7UnPadding(decrypted)
}

// RsaEncrypt 加密
func RsaEncrypt(data, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data)
}

// RsaEncryptBase64 使用 rsa 加密和 base64 编码数据
func RsaEncryptBase64(data, key []byte) ([]byte, error) {
	encrypted, err := RsaEncrypt(data, key)
	if err != nil {
		return nil, err
	}
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(encrypted)))
	enc.Encode(buf, encrypted)
	return buf, err
}

// MustRsaEncryptBase64Str 使用 rsa 加密和 base64 编码数据
func MustRsaEncryptBase64Str(data, key string) string {
	v, _ := RsaEncryptBase64([]byte(data), []byte(key))
	return string(v)
}

// AesEncryptBase64 使用 ase 加密和 base64 编码数据
func AesEncryptBase64(data, key []byte) ([]byte, error) {
	encrypted, err := AesEncrypt(data, key)
	if err != nil {
		return nil, err
	}
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(encrypted)))
	enc.Encode(buf, encrypted)
	return buf, nil
}

// MustAesEncryptBase64Str 使用 ase 加密和 base64 编码数据
func MustAesEncryptBase64Str(data, key string) string {
	v, _ := AesEncryptBase64([]byte(data), []byte(key))
	return string(v)
}

// AesDecryptBase64 使用 base64 解码和 ase 解密数据
func AesDecryptBase64(data, key []byte) ([]byte, error) {
	enc := base64.StdEncoding
	encrypted := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(encrypted, data)
	if err != nil {
		return nil, err
	}
	decrypted, err := AesDecrypt(encrypted[:n], key)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// MustAesDecryptBase64Str 使用 base64 解码和 ase 解密数据
func MustAesDecryptBase64Str(data, key string) string {
	v, _ := AesDecryptBase64([]byte(data), []byte(key))
	return string(v)
}

// 计算 md5 并 hex 编码
func Md5Hex(v string) string {
	d := md5.Sum([]byte(v))
	return hex.EncodeToString(d[:])
}

// 返回[0,1]
func Bool2Int(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// 仅为1时返回true
func Int2Bool(i int8) bool {
	return i-1 == 0
}

// 2006-01-02 15:04:05
type Time1 time.Time

func (t *Time1) UnmarshalJSON(b []byte) error {
	v, err := time.ParseInLocation("2006-01-02 15:04:05 -07", strings.Trim(string(b), "\"")+" +08", time.Local)
	if err != nil {
		return err
	}
	*t = Time1(v)
	return nil
}

// 20060102150405
// 时间戳
type Time2 time.Time

func (t *Time2) UnmarshalJSON(b []byte) error {
	v, err := time.ParseInLocation("20060102150405-07", strings.Trim(string(b), "\"")+"+08", time.Local)
	if err != nil {
		i, err := strconv.ParseInt(strings.Trim(string(b), "\""), 10, 64)
		if err != nil {
			return err
		}
		v = time.UnixMilli(i)
	}
	*t = Time2(v)
	return nil
}

// 时间戳
type Time3 time.Time

func (t *Time3) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(strings.Trim(string(b), "\""), 10, 64)
	if err != nil {
		return err
	}
	*t = Time3(time.UnixMilli(i))
	return nil
}

type Time4 time.Time

func (t *Time4) UnmarshalJSON(b []byte) error {
	v, err := time.ParseInLocation("Jan 2, 2006 15:04:05 PM -07", strings.Trim(string(b), "\"")+" +08", time.Local)
	if err != nil {
		return err
	}
	*t = Time4(v)
	return nil
}

type String string

func (s *String) UnmarshalJSON(b []byte) error {
	*s = String(bytes.Trim(b, "\""))
	return nil
}

type Int int

func (i *Int) UnmarshalJSON(b []byte) error {
	v, err := strconv.Atoi(string(bytes.Trim(b, "\"")))
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}

type Int64 int64

func (i *Int64) UnmarshalJSON(b []byte) error {
	v, err := strconv.ParseInt(string(bytes.Trim(b, "\"")), 10, 64)
	if err != nil {
		return err
	}
	*i = Int64(v)
	return nil
}

type Bool bool

func (i *Bool) UnmarshalJSON(b []byte) error {
	v, err := strconv.ParseBool(string(bytes.Trim(b, "\"")))
	if err != nil {
		return err
	}
	*i = Bool(v)
	return nil
}
