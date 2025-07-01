package sdk

import (
	"strconv"
	"strings"
)

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Is401Started(code int64) bool {
	codeStr := strconv.FormatInt(code, 10)
	return strings.HasPrefix(codeStr, "401")
}
