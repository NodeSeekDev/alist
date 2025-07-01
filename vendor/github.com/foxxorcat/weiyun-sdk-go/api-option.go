package weiyunsdkgo

type ParamOption func(Json)

func WarpParamOption(option ...ParamOption) []ParamOption {
	return option
}
func ApplyParamOption(param Json, option ...ParamOption) Json {
	for _, opt := range option {
		opt(param)
	}
	return param
}

func QueryFileOptionCount(size int) ParamOption {
	return func(j Json) {
		j["count"] = size
	}
}

type OrderBy int8

const (
	_ OrderBy = iota
	FileName
	FileMtime
	FileSize
)

func QueryFileOptionSort(orderBy OrderBy, desc bool) ParamOption {
	return func(j Json) {
		j["sort_field"] = orderBy
		j["reverse_order"] = desc
	}
}

type GetType int8

const (
	FileAndDir = iota
	OnlyDir
	OnlyFile
)

func QueryFileOptionGetType(t GetType) ParamOption {
	return func(j Json) {
		j["get_type"] = t
	}
}

func QueryFileOptionOffest(start int64) ParamOption {
	return func(j Json) {
		j["start"] = start
	}
}
