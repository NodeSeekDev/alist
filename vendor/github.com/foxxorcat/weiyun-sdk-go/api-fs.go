package weiyunsdkgo

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/go-resty/resty/v2"
)

type File struct {
	FileID    string    `json:"file_id"`
	FileName  string    `json:"filename"`
	FileSize  int64     `json:"file_size"`
	FileSha   string    `json:"file_sha"`
	FileCtime TimeStamp `json:"file_ctime"`
	FileMtime TimeStamp `json:"file_mtime"`

	ExtInfo struct {
		ThumbURL string `json:"thumb_url"`
	} `json:"ext_info"`
}

type Folder struct {
	DirKey   string    `json:"dir_key"`
	DirName  string    `json:"dir_name"`
	DirCtime TimeStamp `json:"dir_ctime"`
	DirMtime TimeStamp `json:"dir_mtime"`
}

type FolderPath struct {
	PdirKey string `json:"pdir_key"`
	Folder
}

// 查询文件夹完整路径
func (c *WeiYunClient) LibDirPathGet(dirKey string, opts ...RestyOption) ([]FolderPath, error) {
	var resp struct {
		Items []FolderPath `json:"items"`
	}
	_, err := c.WeiyunFileLibClientRequest("LibDirPathGet", 26150, Json{"dir_key": dirKey}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

type DiskListData struct {
	DirList        []Folder `json:"dir_list"`
	FileList       []File   `json:"file_list"`
	PdirKey        string   `json:"pdir_key"`
	FinishFlag     bool     `json:"finish_flag"`
	TotalDirCount  int64    `json:"total_dir_count"`
	TotalFileCount int64    `json:"total_file_count"`
	TotalSpace     int64    `json:"total_space"`
	HideDirCount   int64    `json:"hide_dir_count"`
	HideFileCount  int64    `json:"hide_file_count"`
}

// 查询文件、文件夹
// 数量限制 500
func (c *WeiYunClient) DiskDirFileList(dirKey string, paramOption []ParamOption, opts ...RestyOption) (*DiskListData, error) {
	param := Json{
		//"pdir_key": pdirKey,
		"dir_key": dirKey,

		"start": 0,
		"count": 500,

		"sort_field":    2,
		"reverse_order": false,

		"get_type": 0,

		"get_abstract_url":    false,
		"get_dir_detail_info": false,
	}
	ApplyParamOption(param, paramOption...)

	var resp DiskListData
	_, err := c.WeiyunQdiskRequest("DiskDirList", 2208, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type DiskDirFileBatchListParam struct {
	DirKey      string
	ParamOption []ParamOption
}

// 批量查询文件、文件夹
func (c *WeiYunClient) DiskDirFileBatchList(batchParam []DiskDirFileBatchListParam, commonParamOption []ParamOption, opts ...RestyOption) ([]DiskListData, error) {
	param := Json{
		//"pdir_key": pdirKey,
		"dir_list": MustSliceConvert(batchParam, func(b DiskDirFileBatchListParam) Json {
			dParam := Json{
				"dir_key": b.DirKey,

				"start": 0,
				"count": 500,

				"sort_field":    2,
				"reverse_order": false,

				"get_type": 0,

				"get_abstract_url":    false,
				"get_dir_detail_info": false,
			}
			ApplyParamOption(dParam, commonParamOption...)
			ApplyParamOption(dParam, b.ParamOption...)
			return dParam
		}),
	}

	var resp struct {
		DirList []DiskListData `json:"dir_list"`
	}

	_, err := c.WeiyunQdiskRequest("DiskDirBatchList", 2209, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.DirList, nil
}

type FolderParam struct {
	PPdirKey string `json:"ppdir_key,omitempty"` // 父父目录ID(打包下载忽略)(移动dst忽略)
	PdirKey  string `json:"pdir_key,omitempty"`  // 父目录ID(打包下载忽略)
	DirKey   string `json:"dir_key,omitempty"`   // 目录ID(创建忽略)
	DirName  string `json:"dir_name,omitempty"`  // 目录名称
}

// 文件夹重命名
func (c *WeiYunClient) DiskDirAttrModify(dParam FolderParam, newDirName string, opts ...RestyOption) error {
	param := Json{
		"ppdir_key":    dParam.PPdirKey,
		"pdir_key":     dParam.PdirKey,
		"dir_key":      dParam.DirKey,
		"src_dir_name": dParam.DirName,
		"dst_dir_name": newDirName,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirAttrModify", 2615, param, nil, opts...)
	return err
}

// 文件夹删除
func (c *WeiYunClient) DiskDirDelete(dParam FolderParam, opts ...RestyOption) error {
	param := Json{
		"dir_list": []FolderParam{dParam},
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchDeleteEx", 2509, param, nil, opts...)
	return err
}

// 文件夹移动
func (c *WeiYunClient) DiskDirMove(srcParam FolderParam, dstParam FolderParam, opts ...RestyOption) error {
	param := Json{
		"src_ppdir_key": srcParam.PPdirKey,
		"src_pdir_key":  srcParam.PdirKey,
		"dir_list":      []FolderParam{srcParam},
		"dst_ppdir_key": dstParam.PdirKey,
		"dst_pdir_key":  dstParam.DirKey,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchMove", 2618, param, nil, opts...)
	return err
}

// 文件夹创建
func (c *WeiYunClient) DiskDirCreate(dParam FolderParam, opts ...RestyOption) (*Folder, error) {
	param := Json{
		"ppdir_key":         dParam.PPdirKey,
		"pdir_key":          dParam.PdirKey,
		"dir_name":          dParam.DirName,
		"file_exist_option": 2,
		"create_type":       1,
	}
	var folder Folder
	_, err := c.WeiyunQdiskClientRequest("DiskDirCreate", 2614, param, &folder, opts...)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

type FileParam struct {
	PPdirKey string `json:"ppdir_key,omitempty"` // 父父目录ID(下载忽略)
	PdirKey  string `json:"pdir_key,omitempty"`  // 父目录ID
	FileID   string `json:"file_id,omitempty"`   // 文件ID
	FileName string `json:"filename,omitempty"`  // 文件名称(下载忽略)
}

// 文件重命名
func (c *WeiYunClient) DiskFileRename(fParam FileParam, newFileName string, opts ...RestyOption) error {
	param := Json{
		"ppdir_key":    fParam.PPdirKey,
		"pdir_key":     fParam.PdirKey,
		"file_id":      fParam.FileID,
		"src_filename": fParam.FileName,
		"filename":     newFileName,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskFileRename", 2605, param, nil, opts...)
	return err
}

// 文件删除
func (c *WeiYunClient) DiskFileDelete(fParam FileParam, opts ...RestyOption) error {
	param := Json{
		"file_list": []FileParam{fParam},
	}

	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchDeleteEx", 2509, param, nil, opts...)
	return err
}

// 文件移动
func (c *WeiYunClient) DiskFileMove(srcParam FileParam, dstParam FolderParam, opts ...RestyOption) error {
	param := Json{
		"src_ppdir_key": srcParam.PPdirKey,
		"src_pdir_key":  srcParam.PdirKey,
		"file_list":     []FileParam{srcParam},
		"dst_ppdir_key": dstParam.PdirKey,
		"dst_pdir_key":  dstParam.DirKey,
	}
	_, err := c.WeiyunQdiskClientRequest("DiskDirFileBatchMove", 2618, param, nil, opts...)
	return err
}

type DiskFileDownloadData struct {
	Retcode int    `json:"retcode"`
	Retmsg  string `json:"retmsg"`

	CookieName  string `json:"cookie_name"`
	CookieValue string `json:"cookie_value"`

	DownloadUrl string `json:"download_url"`
}

// 文件下载
func (c *WeiYunClient) DiskFileDownload(fParam FileParam, opts ...RestyOption) (*DiskFileDownloadData, error) {
	resp, err := c.DiskFileBatchDownload([]FileParam{fParam}, opts...)
	if err != nil {
		return nil, err
	}
	return &resp[0], nil
}

// 批量下载
func (c *WeiYunClient) DiskFileBatchDownload(fParam []FileParam, opts ...RestyOption) ([]DiskFileDownloadData, error) {
	param := Json{
		"file_list":     fParam,
		"download_type": 0,
	}
	var resp struct {
		Data []DiskFileDownloadData `json:"file_list"`
	}
	_, err := c.WeiyunQdiskClientRequest("DiskFileBatchDownload", 2402, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type DiskFilePackageDownloadParam struct {
	PdirKey  string `json:"pdir_key"`
	PdirList any    `json:"pdir_list"`
}

// 文件打包下载
func (c *WeiYunClient) DiskFilePackageDownload(param []DiskFilePackageDownloadParam, zipFilename string, opts ...RestyOption) (*DiskFileDownloadData, error) {
	param_ := Json{
		"pdir_list": MustSliceConvert(param, func(p DiskFilePackageDownloadParam) Json {
			list := batchParamConvert(p.PdirList)
			list["pdir_key"] = p.PdirKey
			return list
		}),
		"zip_filename": zipFilename,
	}

	var resp DiskFileDownloadData
	_, err := c.WeiyunQdiskRequest("DiskFilePackageDownload", 2403, param_, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func batchParamConvert(param any) Json {
	param_ := Json{}
	switch v := param.(type) {
	case FileParam, *FileParam:
		param_["file_list"] = []any{param}
	case FolderParam, *FolderParam:
		param_["dir_list"] = []any{param}
	case []FileParam, []*FileParam:
		param_["file_list"] = param
	case []FolderParam, []*FolderParam:
		param_["dir_list"] = param
	case []any:
		var fileList []any
		var dirList []any
		for _, vv := range v {
			switch vv.(type) {
			case FileParam, *FileParam:
				fileList = append(fileList, vv)
			case []FolderParam, []*FolderParam:
				dirList = append(dirList, vv)
			}
		}
		param_["file_list"] = fileList
		param_["dir_list"] = dirList
	case Json:
		param_ = v
	}
	return param_
}

type UploadAuthData struct {
	UploadKey string `json:"upload_key"`
	Ex        string `json:"ex"`
}

type UploadChannelData struct {
	ID     int   `json:"id"`
	Offset int64 `json:"offset"`
	Len    int   `json:"len"`
}

type PreUploadData struct {
	FileExist bool `json:"file_exist"` // 文件是否存在
	File      File `json:"common_upload_rsp"`

	UploadScr int `json:"upload_scr"` // 未知

	// 上传授权
	UploadAuthData

	ChannelList []UploadChannelData `json:"channel_list"` // 上传通道

	Speedlimit int `json:"speedlimit"` // 上传速度限制
	FlowState  int `json:"flow_state"`

	UploadState     int `json:"upload_state"`      // 上传状态 1:上传未完成,3:该通道无剩余分片，2:上传完成
	UploadedDataLen int `json:"uploaded_data_len"` // 已经上传的长度
}

type UpdloadFileParam struct {
	PdirKey string // 父父目录ID
	DirKey  string // 父目录ID

	FileName string
	FileSize int64
	File     io.ReadSeeker

	ChannelCount int // 上传通道数量

	// 遇到同名文件操作
	// 1(覆盖) 2 to n(重命名)
	FileExistOption int
}

func (c *WeiYunClient) PreUpload(ctx context.Context, param UpdloadFileParam, opts ...RestyOption) (*PreUploadData, error) {
	const blockSize = 1024 * 1024

	var (
		beforeBlockSize int64            // 之前的块总大小
		lastBlockSize   = param.FileSize // 最后一块大小
		checkBlockSize  int64            // checkData大小
	)

	if param.FileSize > 0 {
		lastBlockSize = (param.FileSize % blockSize)
		if lastBlockSize == 0 {
			lastBlockSize = blockSize
		}
		checkBlockSize = lastBlockSize % 128
		if checkBlockSize == 0 {
			checkBlockSize = 128
		}
		beforeBlockSize = param.FileSize - lastBlockSize
	}

	type BlockInfo struct {
		Sha1   string `json:"sha"`
		Offset int64  `json:"offset"`
		Size   int64  `json:"size"`
	}

	var (
		fileHash      string
		checkSha1     string
		checkData     string
		blockInfoList = make([]BlockInfo, 0, (param.FileSize/blockSize)+1)
		hash          = sha1.New()
	)

	// before
	// 计算除最后一块hash

	for offset := int64(0); offset < beforeBlockSize; offset += blockSize {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if _, err := io.CopyN(hash, param.File, blockSize); err != nil {
			return nil, err
		}
		blockInfoList = append(blockInfoList, BlockInfo{
			Sha1:   hex.EncodeToString(GetSha1State(hash)),
			Offset: offset,
			Size:   blockSize,
		})
	}

	/*
		    校验逻辑
			hash := Sha1(checkSha1)
			hash.Write(checkData)
			hash.Sum() == fileHash
	*/

	// between
	// 得到校验点checkSha1

	if _, err := io.CopyN(hash, param.File, lastBlockSize-checkBlockSize); err != nil {
		return nil, err
	}
	checkSha1 = hex.EncodeToString(GetSha1State(hash))

	// after
	// 得到校验数据checkData
	var buf [128]byte
	_, err := io.ReadFull(io.TeeReader(param.File, hash), buf[:checkBlockSize])
	if err != nil {
		return nil, err
	}
	checkData = base64.StdEncoding.EncodeToString(buf[:checkBlockSize])
	fileHash = hex.EncodeToString(hash.Sum(nil))

	blockInfoList = append(blockInfoList, BlockInfo{
		Sha1:   fileHash,
		Offset: beforeBlockSize,
		Size:   lastBlockSize,
	})

	paramJson := Json{
		"common_upload_req": Json{
			"ppdir_key":         param.PdirKey,
			"pdir_key":          param.DirKey,
			"file_size":         param.FileSize,
			"filename":          param.FileName,
			"file_exist_option": param.FileExistOption,
			"use_mutil_channel": true,
		},
		"upload_scr":      0,
		"channel_count":   param.ChannelCount,
		"block_size":      blockSize,
		"check_sha":       checkSha1,
		"check_data":      checkData,
		"block_info_list": blockInfoList,
	}

	if _, err := param.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var resp PreUploadData
	_, err = c.UploadRequest("PreUpload", 247120, paramJson, &resp, append([]RestyOption{func(request *resty.Request) { request.SetContext(ctx) }}, opts...)...)
	if err != nil {
		return nil, err
	}
	resp.File.FileSha = fileHash
	resp.File.FileSize = param.FileSize
	return &resp, nil
}

type AddChannelData struct {
	// 源通道数量
	OrigChannelCount int `json:"orig_channel_count"`
	// 当前通道数量
	FinalChannelCount int `json:"final_channel_count"`

	// 源通道信息
	OrigChannels []UploadChannelData `json:"orig_channels"`
	// 增加通道信息
	AddChannels []UploadChannelData `json:"channels"`
}

// 增加上传通道
func (c *WeiYunClient) AddUploadChannel(origChannelCount, destChannelCount int, auth UploadAuthData, opts ...RestyOption) (*AddChannelData, error) {
	param := Json{
		"upload_key": auth.UploadKey,
		"ex":         auth.Ex,

		"orig_channel_count": origChannelCount,
		"dest_channel_count": destChannelCount,

		"speed": 4303,
	}

	var resp AddChannelData
	_, err := c.UploadRequest("AddChannel", 247122, param, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type UploadPieceData struct {
	Channel UploadChannelData `json:"channel"` // 下一个分片
	Ex      string            `json:"ex"`

	UploadState int `json:"upload_state"` // 上传状态 1:上传未完成,3:该通道无剩余分片，2:上传完成
	FlowState   int `json:"flow_state"`
}

func (c *WeiYunClient) UploadFile(ctx context.Context, channel UploadChannelData, auth UploadAuthData, r io.Reader, opts ...RestyOption) (*UploadPieceData, error) {
	param := Json{
		"upload_key": auth.UploadKey,
		"ex":         auth.Ex,
		"channel":    channel,

		// 用于传递文件句柄
		"fileReader": r,
	}

	var resp UploadPieceData
	_, err := c.UploadRequest("UploadPiece", 247121, param, &resp, append([]RestyOption{func(request *resty.Request) {
		request.SetContext(ctx)
	}}, opts...)...)
	if err != nil {
		return nil, err
	}

	if resp.Channel.Len == 0 && resp.UploadState == 1 {
		resp.Channel.Len = channel.Len
	}
	return &resp, nil
}
