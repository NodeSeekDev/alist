package json_types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Int64OrStr int64

func (i *Int64OrStr) UnmarshalJSON(data []byte) error {
	// 尝试解析为int64
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*i = Int64OrStr(num)
		return nil
	}

	// 尝试解析为字符串，再转换
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("field must be a number or string")
	}

	// 去除可能的引号并转换
	parsed, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number string: %s", str)
	}

	*i = Int64OrStr(parsed)
	return nil
}
