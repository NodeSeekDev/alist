# wopan-sdk-go
 Wopan SDK for the Go programming language

## Installation

```bash
go get github.com/xhofe/wopan-sdk-go
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/xhofe/wopan-sdk-go"
)

func main() {
	w := wopan.DefaultWithRefreshToken("91d4b946-xxxx-4909-bac1-d9914e45f2de")
	res, err := w.AppQueryUser()
	if err != nil {
		fmt.Printf("AppQueryUser() error = %v", err)
	} else {
		fmt.Printf("AppQueryUser() = %+v", res)
	}
}

```