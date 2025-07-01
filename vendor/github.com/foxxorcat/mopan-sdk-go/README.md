# mopan-sdk-go
 mopan SDK for the Go programming language

## Installation

```bash
go get github.com/foxxorcat/mopan-sdk-go
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/foxxorcat/mopan-sdk-go"
)

func main() {
	m := mopan.NewClient().SetAuthorization("ZR7Lxxx")
	res, err := m.GetUserInfo()
	if err != nil {
		fmt.Printf("GetUserInfo() error = %v", err)
	} else {
		fmt.Printf("GetUserInfo() = %+v", res)
	}
}

```