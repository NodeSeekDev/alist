package json_types

import (
	"encoding/json"
	"fmt"
)

type Int64OrFloat struct {
	Int64 int64
	Float float64
}

func (i *Int64OrFloat) UnmarshalJSON(data []byte) error {
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		i.Int64 = num
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		i.Float = f
		return nil
	}
	return fmt.Errorf("data is neither int64 nor float64")
}
