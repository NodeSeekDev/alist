package weiyunsdkgo

import (
	"encoding"
	"hash"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

/*
*	Cookie
**/

func GetCookieValue(name string, cks []*http.Cookie) string {
	for _, ck := range cks {
		if ck.Name == name && ck.Value != "" {
			return ck.Value
		}
	}
	return ""
}

func SetCookieValue(name string, value string, cks []*http.Cookie) {
	for _, ck := range cks {
		if ck.Name == name {
			ck.Value = value
		}
	}
}

func ParseCookieStr(str string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", str)
	request := http.Request{Header: header}
	return request.Cookies()
}

func CookieToString(cookies []*http.Cookie) string {
	if cookies == nil {
		return ""
	}
	cookieStrings := make([]string, len(cookies))
	for i, cookie := range cookies {
		cookieStrings[i] = cookie.Name + "=" + cookie.Value
	}
	return strings.Join(cookieStrings, ";")
}

func ClearCookie(cookies []*http.Cookie) []*http.Cookie {
	newCookies := make([]*http.Cookie, 0, len(cookies))
S:
	for _, cookie := range cookies {
		// 去空
		if cookie != nil && cookie.Value != "" {
			// 去重复,保留最后一个
			for _, newCookie := range newCookies {
				if newCookie.Name == cookie.Name {
					*newCookie = *cookie
					continue S
				}
			}
			newCookies = append(newCookies, cookie)
		}
	}
	return newCookies
}

/*
*	Slice
**/
func SliceContains[T comparable](arrs []T, v T, contains func(arr, v T) bool) bool {
	for _, vv := range arrs {
		if contains(vv, v) {
			return true
		}
	}
	return false
}

// SliceConvert convert slice to another type slice
func SliceConvert[S any, D any](srcS []S, convert func(src S) (D, error)) ([]D, error) {
	res := make([]D, 0, len(srcS))
	for i := range srcS {
		dst, err := convert(srcS[i])
		if err != nil {
			return nil, err
		}
		res = append(res, dst)
	}
	return res, nil
}

func MustSliceConvert[S any, D any](srcS []S, convert func(src S) D) []D {
	res := make([]D, 0, len(srcS))
	for i := range srcS {
		dst := convert(srcS[i])
		res = append(res, dst)
	}
	return res
}

/*
*	Json
**/
// 时间戳
type TimeStamp time.Time

func (t *TimeStamp) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(strings.Trim(string(b), "\""), 10, 64)
	if err != nil {
		return err
	}
	*t = TimeStamp(time.UnixMilli(i))
	return nil
}

func MustJsonMarshalToString(v any) string {
	s, _ := jsoniter.MarshalToString(v)
	return s
}

/*
*	Other
**/

func GetSha1State(h hash.Hash) []byte {
	v, _ := h.(encoding.BinaryMarshaler).MarshalBinary()
	// 按4byte调整顺序
	new := make([]byte, 0, 20)
	for v = v[4:][:20]; len(v) > 0; v = v[4:] {
		new = append(new, v[3], v[2], v[1], v[0])
	}
	return new
}

func RandomT() string {
	return strconv.FormatFloat(rand.Float64(), 'f', 16, 32)
}

func GetHash33(d string) string {
	var e int
	for _, t := range d {
		e += (e << 5) + int(t)
	}
	return strconv.Itoa(0x7fffffff & e)
}

// func GetGTK(pskey string) string {
// 	e := 5381
// 	for _, t := range pskey {
// 		e += (e << 5) + int(t)
// 	}
// 	return strconv.Itoa(0x7fffffff & e)
// }

// func GetBKN(skey string) string {
// 	return GetGTK(skey)
// }

func GetDirFileIDFormUrl(url string) string {
	i := strings.LastIndexByte(url, '/')
	if i == -1 {
		return ""
	}
	return url[i+1:]
}
