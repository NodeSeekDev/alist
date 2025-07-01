# weiyun-sdk-go
 weiyun SDK for the Go programming language

## Installation

```bash
go get github.com/foxxorcat/weiyun-sdk-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	weiyunsdkgo "github.com/foxxorcat/weiyun-sdk-go"
)

func main() {
	client := weiyunsdkgo.NewWeiYunClient()
	_, err := client.QQQRLogin(context.TODO(), func(b []byte) error {
		os.WriteFile("qr.png", b, 0777)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	// client.SetCookiesStr("cookies")
	// client.QQFastLogin(context.Background(), "qq number")
	userInfo, err := client.DiskUserInfoGet()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userInfo)
}

```