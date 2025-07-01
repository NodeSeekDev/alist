# go-cache
 
Modified from <https://github.com/fanjindong/go-cache> with Generics in Go 1.18+.

## Install

`go get -u github.com/Xhofe/go-cache`

## Basic Usage

```go
import "github.com/Xhofe/go-cache"

func main() {
    c := cache.NewMemCache[int]()
    c.Set("a", 1)
    c.Set("b", 1, cache.WithEx[int](1*time.Second))
    time.sleep(1*time.Second)
    c.Get("a") // 1, true
    c.Get("b") // nil, false
}
```