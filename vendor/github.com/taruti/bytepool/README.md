# bytepool


Alloc temporary variable sized byte buffers from sync.Pools.

Example:
```
bs := bytepool.Alloc(134)
defer bytepool.Free(bs)
n,e := myfile.Read(bs)
bs = bs[0:n]
```

# License: MIT

# [![GoDoc](https://godoc.org/github.com/taruti/bytepool?status.png)](http://godoc.org/github.com/taruti/bytepool)

# Test code coverage 100%
