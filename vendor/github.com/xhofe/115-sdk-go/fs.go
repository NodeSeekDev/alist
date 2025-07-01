// https://www.yuque.com/115yun/open/qur839kyx9cgxpxi
package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/xhofe/115-sdk-go/json_types"
)

type MkdirResp struct {
	FileName string `json:"file_name"`
	FileID   string `json:"file_id"`
}

func (c *Client) Mkdir(ctx context.Context, pid, filename string) (*MkdirResp, error) {
	var resp MkdirResp
	_, err := c.AuthRequest(ctx, ApiFsMkdir, http.MethodPost, &resp, ReqWithForm(Form{
		"pid":       pid,
		"file_name": filename,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type GetFilesReq struct {
	CID         string
	Type        string
	Limit       int64
	Offset      int64
	Suffix      string
	ASC         bool   // 1: asc, 0: desc
	O           string // order by file_name|file_size|user_utime|file_type
	CustomOrder int
	Stdir       int  // 筛选文件时， 是否显示文件夹；1:要展示文件夹 0不展示
	Star        bool // 1: filter star files
	Cur         int  // 是否只显示当前文件夹内文件
	ShowDir     bool // 是否显示目录；0 或 1，默认为0
}

type GetFilesResp_File_FL struct {
	Id         string `json:"id"`          // 文件标签ID
	Name       string `json:"name"`        // 文件标签名称
	Sort       string `json:"sort"`        // 文件标签排序
	Color      string `json:"color"`       // 文件标签颜色
	IsDefault  int    `json:"is_default"`  // 文件标签类型；0：最近使用；1：非最近使用；2：为默认标签
	UpdateTime int64  `json:"update_time"` // 文件标签更新时间
	CreateTime int64  `json:"create_time"` // 文件标签创建时间
}

type GetFilesResp_File struct {
	Fid       string                  `json:"fid"`  // 文件ID
	Aid       string                  `json:"aid"`  // 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除
	Pid       string                  `json:"pid"`  // 父文件夹ID
	Fc        string                  `json:"fc"`   // 文件分类 0 文件夹 1 文件
	Fn        string                  `json:"fn"`   // 文件名
	Fco       string                  `json:"fco"`  // 文件夹封面
	Ism       string                  `json:"ism"`  // 是否星标，1：星标
	Isp       int                     `json:"isp"`  // 是否加密；1：加密
	Pc        string                  `json:"pc"`   // 文件提取码
	Upt       int64                   `json:"upt"`  // 修改时间
	Uet       int64                   `json:"uet"`  // 修改时间
	UpPt      int64                   `json:"uppt"` // 上传时间
	Cm        int64                   `json:"cm"`
	FDesc     string                  `json:"fdesc"`     // 文件备注
	IsPl      int64                   `json:"ispl"`      // 是否统计文件夹下视频时长开关
	Fl        []GetFilesResp_File_FL  `json:"fl"`        // 文件标签
	Sha1      string                  `json:"sha1"`      // 文件sha1
	FS        int64                   `json:"fs"`        // 文件大小
	Fta       string                  `json:"fta"`       // 文件状态 0/2 未上传完成，1 已上传完成
	Ico       string                  `json:"ico"`       // 文件后缀名
	Fatr      string                  `json:"fatr"`      // 音频长度
	IsV       int64                   `json:"isv"`       // 是否视频文件
	Def       int64                   `json:"def"`       // 视频清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画
	Def2      int64                   `json:"def2"`      // 视频清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画
	PlayLong  json_types.Int64OrFloat `json:"play_long"` // 音视频时长
	VImg      string                  `json:"v_img"`
	Thumbnail string                  `json:"thumb"` // 图片缩略图
	Uo        string                  `json:"uo"`    // 原图地址
}

type GetFilesResp struct {
	Resp[[]GetFilesResp_File]
	Count          int64                 `json:"count"`     // 当前目录文件数量
	SysCount       int64                 `json:"sys_count"` // 系统文件夹数量
	Offset         int64                 `json:"offset"`    // 偏移量
	Limit          json_types.Int64OrStr `json:"limit"`     // 分页量
	Aid            int                   `json:"aid"`       // 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除
	Cid            int64                 `json:"cid"`       // 父目录ID
	IsAsc          int                   `json:"is_asc"`    // 1: asc, 0: desc
	MinSize        int64                 `json:"min_size"`
	MaxSize        int64                 `json:"max_size"`
	SysDir         string                `json:"sys_dir"`
	HideData       string                `json:"hide_data"`        //是否返回文件数据
	RecordOpenTime string                `json:"record_open_time"` //是否记录文件夹的打开时间
	Star           int                   `json:"star"`             //是否星标；1：星标；0：未星标
	Type           int                   `json:"type"`             //一级筛选大分类，1：文档，2：图片，3：音乐，4：视频，5：压缩包，6：应用
	Suffix         string                `json:"suffix"`           //一级筛选选其他时填写的后缀名
	Path           []struct {
		Name string                `json:"name"` //父目录名
		Aid  json_types.Int64OrStr `json:"aid"`
		Cid  json_types.Int64OrStr `json:"cid"`
		Pid  json_types.Int64OrStr `json:"pid"`
		Isp  json_types.Int64OrStr `json:"isp"`
		PCid string                `json:"p_cid"`
		Fv   string                `json:"fv"`
	} `json:"path"` //父目录树
	Cur    int64  `json:"cur"`
	StdDir int    `json:"stdir"`
	Fields string `json:"fields"`
	Order  string `json:"order"`
}

// GetFiles: https://www.yuque.com/115yun/open/kz9ft9a7s57ep868
func (c *Client) GetFiles(ctx context.Context, req *GetFilesReq) (*GetFilesResp, error) {
	var resp GetFilesResp
	_, err := c.AuthRequestRaw(ctx, ApiFsGetFiles, http.MethodGet, &resp, ReqWithQuery(Form{
		"cid":          req.CID,
		"type":         req.Type,
		"limit":        strconv.FormatInt(req.Limit, 10),
		"offset":       strconv.FormatInt(req.Offset, 10),
		"suffix":       req.Suffix,
		"asc":          Ternary(req.ASC, "1", "0"),
		"o":            req.O,
		"custom_order": strconv.Itoa(req.CustomOrder),
		"stdir":        strconv.Itoa(req.Stdir),
		"star":         Ternary(req.Star, "1", "0"),
		"cur":          strconv.Itoa(req.Cur),
		"show_dir":     Ternary(req.ShowDir, "1", "0"),
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type GetFolderInfoResp struct {
	Count        string `json:"count"`
	Size         string `json:"size"`
	FolderCount  int64  `json:"folder_count"`
	PlayLong     int64  `json:"play_long"`
	ShowPlayLong int64  `json:"show_play_long"`
	PTime        string `json:"ptime"`
	UTime        string `json:"utime"`
	FileName     string `json:"file_name"`
	PickCode     string `json:"pick_code"`
	Sha1         string `json:"sha1"`
	FileID       string `json:"file_id"`
	IsMark       string `json:"is_mark"`
	OpenTime     int64  `json:"open_time"`
	FileCategory string `json:"file_category"`
	Paths        []struct {
		FileID   int64  `json:"file_id"`
		FileName string `json:"file_name"`
	} `json:"paths"`
}

// GetFolderInfo: https://www.yuque.com/115yun/open/rl8zrhe2nag21dfw
func (c *Client) GetFolderInfo(ctx context.Context, fileID string) (*GetFolderInfoResp, error) {
	var resp GetFolderInfoResp
	_, err := c.AuthRequest(ctx, ApiFsGetFolderInfo, http.MethodGet, &resp, ReqWithQuery(Form{
		"file_id": fileID,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type SearchFilesReq struct {
	SearchValue string //查找关键字
	Limit       int64  //单页记录数，默认20，offset+limit最大不超过10000
	Offset      int64  //数据显示偏移量
	FileLabel   string //支持文件标签搜索
	CID         string //目标目录cid=-1时，表示不返回列表任何内容
	GteDay      string //搜索结果匹配的开始时间；格式：2020-11-19
	LteDay      string //搜索结果匹配的结束时间；格式：2020-11-20
	FC          string //只显示文件或文件夹。1 只显示文件夹，2 只显示文件
	Type        string //一级筛选大分类，1：文档，2：图片，3：音乐，4：视频，5：压缩包，6：应用
	Suffix      string //一级筛选选其他时填写的后缀名
}

type SearchFilesResp struct {
	Resp[[]struct {
		FileID       string `json:"file_id"`       // 文件ID
		UserID       string `json:"user_id"`       // 用户ID
		Sha1         string `json:"sha1"`          // 文件sha1值
		FileName     string `json:"file_name"`     // 文件名称
		FileSize     string `json:"file_size"`     // 文件大小
		UserPtime    string `json:"user_ptime"`    // 上传时间
		UserUtime    string `json:"user_utime"`    // 更新时间
		PickCode     string `json:"pick_code"`     // 文件提取码
		ParentID     string `json:"parent_id"`     // 父目录ID
		AreaID       string `json:"area_id"`       // 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除
		IsPrivate    int    `json:"is_private"`    // 文件是否隐藏。0 未隐藏，1 已隐藏
		FileCategory string `json:"file_category"` // 1：文件；0；文件夹
		Ico          string `json:"ico"`           // 文件后缀名
	}]
	Count  int64 `json:"count"`  // 搜索符合条件的文件(夹)总数
	Limit  int64 `json:"limit"`  // 单页记录数
	Offset int64 `json:"offset"` // 数据显示偏移量
}

// SearchFiles: https://www.yuque.com/115yun/open/ft2yelxzopusus38
func (c *Client) SearchFiles(ctx context.Context, req *SearchFilesReq) (*SearchFilesResp, error) {
	var resp SearchFilesResp
	_, err := c.AuthRequestRaw(ctx, ApiFsSearchFiles, http.MethodGet, &resp, ReqWithQuery(Form{
		"search_value": req.SearchValue,
		"limit":        strconv.FormatInt(req.Limit, 10),
		"offset":       strconv.FormatInt(req.Offset, 10),
		"file_label":   req.FileLabel,
		"cid":          req.CID,
		"gte_day":      req.GteDay,
		"lte_day":      req.LteDay,
		"fc":           req.FC,
		"type":         req.Type,
		"suffix":       req.Suffix,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type CopyReq struct {
	PID     string `json:"pid"`      // 目标目录，即所需移动到的目录
	FileID  string `json:"file_id"`  // 所复制的文件和目录ID，多个文件和目录请以 , 隔开
	NoDupli string `json:"no_dupli"` // 复制的文件在目标目录是否允许重复，默认0：0：可以；1：不可以
}

// Copy: https://www.yuque.com/115yun/open/lvas49ar94n47bbk
func (c *Client) Copy(ctx context.Context, req *CopyReq) (any, error) {
	var resp any
	_, err := c.AuthRequest(ctx, ApiFsCopy, http.MethodPost, &resp, ReqWithForm(Form{
		"pid":      req.PID,
		"file_id":  req.FileID,
		"no_dupli": req.NoDupli,
	}))
	return resp, err
}

type MoveReq struct {
	FileIDs string `json:"file_ids"` // 需要移动的文件(夹)ID
	ToCid   string `json:"to_cid"`   // 要移动所在的目录ID，根目录为0
}

// Move: https://www.yuque.com/115yun/open/vc6fhi2mrkenmav2
func (c *Client) Move(ctx context.Context, req *MoveReq) (any, error) {
	var resp any
	_, err := c.AuthRequest(ctx, ApiFsMove, http.MethodPost, &resp, ReqWithForm(Form{
		"file_ids": req.FileIDs,
		"to_cid":   req.ToCid,
	}))
	return resp, err
}

type DownURLResp = map[string]struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	PickCode string `json:"pick_code"`
	Sha1     string `json:"sha1"`
	URL      struct {
		URL string `json:"url"`
	} `json:"url"`
}

// DownURL: https://www.yuque.com/115yun/open/um8whr91bxb5997o
func (c *Client) DownURL(ctx context.Context, pickCode string, ua string) (DownURLResp, error) {
	var resp DownURLResp
	_, err := c.AuthRequest(ctx, ApiFsDownURL, http.MethodPost, &resp, ReqWithForm(Form{
		"pick_code": pickCode,
	}), ReqWithUA(ua))
	if err != nil {
		return nil, err
	}
	return resp, err
}

type UpdateFileReq struct {
	FileID  string `json:"file_id"`   // 需要更改名字的文件(夹)ID
	FileNma string `json:"file_name"` // 新的文件(夹)名字(文件夹名称限制255字节)
	Star    string `json:"star"`      // 是否星标；1：星标；0；取消星标
}

type UpdateFileResp struct {
	FileName string `json:"file_name"`
	Star     string `json:"star"`
}

// UpdateFile: https://www.yuque.com/115yun/open/gyrpw5a0zc4sengm
func (c *Client) UpdateFile(ctx context.Context, req *UpdateFileReq) (*UpdateFileResp, error) {
	var resp UpdateFileResp
	_, err := c.AuthRequest(ctx, ApiFsUpdate, http.MethodPost, &resp, ReqWithForm(Form{
		"file_id":   req.FileID,
		"file_name": req.FileNma,
		"star":      req.Star,
	}))
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type DelFileReq struct {
	FileIDs  string `json:"file_ids"`  // 需要删除的文件(夹)ID
	ParentID string `json:"parent_id"` // 删除的文件(夹)ID所在的父目录ID
}

// DelFile: https://www.yuque.com/115yun/open/kt04fu8vcchd2fnb
func (c *Client) DelFile(ctx context.Context, req *DelFileReq) ([]string, error) {
	var resp []string
	_, err := c.AuthRequest(ctx, ApiFsDelete, http.MethodPost, &resp, ReqWithForm(Form{
		"file_ids":  req.FileIDs,
		"parent_id": req.ParentID,
	}))
	return resp, err
}

type RbListResp_FileInfo struct {
	ID         string `json:"id"`
	FileName   string `json:"file_name"`
	Type       string `json:"type"`
	FileSize   string `json:"file_size"`
	Dtime      string `json:"dtime"`
	ThumbURL   string `json:"thumb_url"`
	Status     string `json:"status"`
	CID        string `json:"cid"`
	ParentName string `json:"parent_name"`
	PickCode   string `json:"pick_code"`
	IsV        int    `json:"isv"`
	Ico        string `json:"ico"`
	Muc        string `json:"muc"`
	DImg       string `json:"d_img"`
	SHA1       string `json:"sha1"`
}

type RbListResp struct {
	Offset int64                          `json:"offset"`  // 偏移量
	Limit  int64                          `json:"limit"`   // 分页量
	Count  string                         `json:"count"`   // 回收站文件总数
	RbPass int                            `json:"rb_pass"` // 是否设置回收站密码
	Files  map[string]RbListResp_FileInfo `json:"-"`
}

// RbList: https://www.yuque.com/115yun/open/bg7l4328t98fwgex
func (c *Client) RbList(ctx context.Context, limit, offset int64) (*RbListResp, error) {
	var rawBytes json.RawMessage
	_, err := c.AuthRequest(ctx, ApiFsRbList, http.MethodGet, &rawBytes, ReqWithQuery(Form{
		"offset": strconv.FormatInt(offset, 10),
		"limit":  strconv.FormatInt(limit, 10),
	}))
	if err != nil {
		return nil, err
	}
	var resp RbListResp
	err = json.Unmarshal(rawBytes, &resp)
	if err != nil {
		return nil, err
	}
	var rawFiles map[string]json.RawMessage
	err = json.Unmarshal(rawBytes, &rawFiles)
	if err != nil {
		return nil, err
	}
	for key, value := range rawFiles {
		if SliceContains([]string{"offset", "limit", "count", "rb_pass"}, key) {
			continue
		}
		var fileInfo RbListResp_FileInfo
		err = json.Unmarshal(value, &fileInfo)
		if err != nil {
			return nil, err
		}
		resp.Files[key] = fileInfo
	}
	return &resp, err
}

type RbRevertResp map[string]struct {
	State bool   `json:"state"`
	Error string `json:"error"`
	ErrNo int    `json:"errno"`
}

// RbRevert: https://www.yuque.com/115yun/open/gq293z80a3kmxbaq
func (c *Client) RbRevert(ctx context.Context, tid string) (RbRevertResp, error) {
	var resp RbRevertResp
	_, err := c.AuthRequest(ctx, ApiFsRbRevert, http.MethodPost, &resp, ReqWithForm(Form{
		"tid": tid,
	}))
	return resp, err
}

// RbDel: https://www.yuque.com/115yun/open/gwtof85nmboulrce
func (c *Client) RbDelete(ctx context.Context, tid string) ([]string, error) {
	var resp []string
	_, err := c.AuthRequest(ctx, ApiFsRbDelete, http.MethodPost, &resp, ReqWithForm(Form{
		"tid": tid,
	}))
	return resp, err
}
