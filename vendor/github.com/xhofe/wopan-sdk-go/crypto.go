package wopan

import "encoding/base64"

type Crypto struct {
	key       []byte
	iv        []byte
	accessKey []byte
}

func NewCrypto() *Crypto {
	c := &Crypto{
		key: []byte(DefaultClientSecret),
		iv:  []byte("wNSOYIB1k1DjY5lA"),
	}
	return c
}

func (c *Crypto) SetAccessToken(token string) error {
	if len(token) < 16 {
		return ErrInvalidAccessToken
	}
	c.accessKey = []byte(token[:16])
	return nil
}

func (c *Crypto) EncryptBytes(bs []byte, channel string) (string, error) {
	key := c.accessKey
	if channel == ChannelAPIUser {
		key = c.key
	}
	res, err := AesEncrypt(bs, key, c.iv)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func (c *Crypto) Encrypt(content string, channel string) (string, error) {
	return c.EncryptBytes([]byte(content), channel)
}

func (c *Crypto) Decrypt(content string, channel string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	key := c.accessKey
	if channel == ChannelAPIUser {
		key = c.key
	}
	res, err := AesDecrypt(bs, key, c.iv)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *Crypto) UserEncrypt(content string) (string, error) {
	return c.Encrypt(content, ChannelAPIUser)
}

func (c *Crypto) UserDecrypt(content string) (string, error) {
	return c.Decrypt(content, ChannelAPIUser)
}

func (c *Crypto) WoHomeEncrypt(content string) (string, error) {
	return c.Encrypt(content, ChannelWoHome)
}

func (c *Crypto) WoHomeDecrypt(content string) (string, error) {
	return c.Decrypt(content, ChannelWoHome)
}
