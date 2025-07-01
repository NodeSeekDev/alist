package wopan

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func calHeader(channel, key string) Header {
	resTime := time.Now().UnixMilli()
	reqSeq := rand.Int31n(8999) + 1e5
	version := ""
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("%s%d%d%s%s", key, resTime, reqSeq, channel, version)))
	sign := hex.EncodeToString(m.Sum(nil))
	return Header{
		Key:     key,
		ResTime: resTime,
		ReqSeq:  int(reqSeq),
		Channel: channel,
		Sign:    sign,
		Version: version,
	}
}
