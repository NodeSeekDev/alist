package mopan

import (
	"strconv"
)

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

// 查询每页大小
func QueryFileOptionPageSize(size int) ParamOption {
	return func(j Json) {
		j["pageSize"] = strconv.Itoa(size)
	}
}

type IconOption int8

const (
	IconSmall IconOption = (1 << iota)
	IconMedium
	IconLarge
	IconMax600
)

// 返回icon方案
// 关闭: icon 0
func QueryFileOptionIconOption(icon IconOption) ParamOption {
	return func(j Json) {
		j["iconOption"] = icon
	}
}

type OrderBy int8

const (
	S_FileName OrderBy = iota
	S_LastOpTime
	S_FileSize
)

// 数据排序
func QueryFileOptionSort(orderBy OrderBy, desc bool) ParamOption {
	return func(j Json) {
		switch orderBy {
		case S_FileSize:
			j["orderBy"] = "filesize"
		case S_LastOpTime:
			j["orderBy"] = "lastOpTime"
		case S_FileName:
			fallthrough
		default:
			j["orderBy"] = "filename"
		}
		j["descending"] = desc
	}
}

// 查询媒体信息
// 似乎没用
func QueryFileOptionMediaAttr(flag bool) ParamOption {
	return func(j Json) {
		j["mediaAttr"] = Bool2Int(flag)
	}
}

type fileType int16

const (
	AllType fileType = iota
	FileType
	FolderType
	ImageType
	AudioType
	VideoType
	DocumentType
)

func QueryFileOptionFileType(fileType fileType) ParamOption {
	return func(j Json) {
		switch fileType {

		case FileType:
			j["fileType"] = 1
		case FolderType:
			j["fileType"] = 2
		case AllType, ImageType, AudioType, VideoType, DocumentType:
			fallthrough
		default:
			j["fileType"] = 0
		}
		switch fileType {

		case ImageType:
			j["mediaType"] = 1
		case AudioType:
			j["mediaType"] = 2
		case VideoType:
			j["mediaType"] = 3
		case DocumentType:
			j["mediaType"] = 4
		case AllType, FileType, FolderType:
			fallthrough
		default:
			j["mediaType"] = 0
		}
	}
}

// 是否递归搜索
func SearchFilesOptionRecursive(recursive bool) ParamOption {
	return func(j Json) {
		j["recursive"] = Bool2Int(recursive)
	}
}

// 通用: 查询共享空间
// 如果 cloudId 为空则不做任何操作
func ParamOptionShareFile(cloudId string) ParamOption {
	return func(j Json) {
		if cloudId != "" {
			j["source"] = 2
			j["cloudId"] = cloudId
		}
	}
}

type updloadOpertype int

const (
	U_Rename    updloadOpertype = 1
	U_Overwrite updloadOpertype = 3
)

func UpdloadOptionOpertype(opertype updloadOpertype) ParamOption {
	return func(j Json) {
		j["opertype"] = opertype
	}
}

// 通用: 查询相册
// func QureyOptionAlbumFile() ParamOption {
// 	return func(j Json) {
// 		j["type"] = 2
// 	}
// }
