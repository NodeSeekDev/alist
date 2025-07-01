package wopan

import "errors"

var (
	JsonClientIDSecret = Json{
		"clientId": DefaultClientID,
		"secret":   true,
	}
	JsonSecret = Json{
		"secret": true,
	}
)

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
)
