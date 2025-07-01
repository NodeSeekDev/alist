package mopan

import (
	"strconv"
)

type (
	IconBase struct {
		Icon struct {
			SmallURL  string `json:"smallUrl,omitempty"`  // icon:1|3|5|7|9|11|13|15
			MediumUrl string `json:"mediumUrl,omitempty"` // icon:2|3|6|7|10|11|14|15
			LargeUrl  string `json:"largeUrl,omitempty"`  // icon:4|5|6|7|12|13|14|15
			Max600    string `json:"max600,omitempty"`    // icon:8|9|10|11|12|13|14|15
		} `json:"icon,omitempty"`
	}
	Folder struct {
		ID       String `json:"id"`       // 文件夹ID
		ParentID String `json:"parentId"` // 文件夹父ID
		Name     string `json:"name"`     // 名称

		FileListSize Int `json:"fileListSize"` // 列表大小

		IsCollect bool `json:"isCollect"` // 是否收藏

		Rev        Time2 `json:"rev"`        // 查询时间
		CreateDate Time1 `json:"createDate"` // 创建时间
		LastOpTime Time1 `json:"lastOpTime"` // 修改时间

		// 未知作用

		StarLabel Int `json:"starLabel"`
		FileCata  Int `json:"fileCata"`
	}

	File struct {
		ID   String `json:"id"`
		Name string `json:"name"` // 文件名
		Size Int64  `json:"size"` // 文件大小
		Md5  string `json:"md5"`

		IconBase

		MediaType int  `json:"mediaType"` // 媒体类型
		IsCollect bool `json:"isCollect"` // 是否收藏
		IsJuBao   bool `json:"isJuBao"`   // 被举报??

		Rev        Time2 `json:"rev"` // 查询时间
		CreateDate Time1 `json:"createDate"`
		LastOpTime Time1 `json:"lastOpTime"`

		// 未知作用

		StarLabel   int `json:"starLabel"`
		FileCata    int `json:"fileCata"`
		Orientation int `json:"orientation,omitempty"`
	}

	FolderInfo struct {
		FileID   String `json:"fileId"`   // 文件夹ID
		ParentID String `json:"parentId"` // 父ID
		FileName string `json:"fileName"` // 文件夹名称

		FilePath string `json:"filePath"` // 文件夹路径

		Rev           Time2 `json:"rev"`           // 查询时间
		CreateTime    Time3 `json:"createTime"`    // 创建时间
		CreateDate    Time1 `json:"createDate"`    // 创建时间
		LastOpTime    Time3 `json:"lastOpTime"`    // 修改时间
		LastOpTimeStr Time1 `json:"lastOpTimeStr"` // 修改时间

		IsCollect bool `json:"isCollect"` // 是否收藏
	}

	FileInfo struct {
		ID       String `json:"id"`       // 文件ID
		ParentID String `json:"parentId"` // 父ID
		Name     string `json:"name"`     // 文件名
		Size     Int64  `json:"size"`     // 文件大小
		Md5      string `json:"md5"`

		NickName string `json:"nickName"` // 用户名称
		FilePath string `json:"filePath"` // 文件路径

		IconBase

		FileDownloadURL string `json:"fileDownloadUrl"` // 下载链接

		MediaType int  `json:"mediaType"` // 媒体类型
		IsCollect bool `json:"isCollect"` // 是否收藏
		IsJuBao   bool `json:"isJuBao"`   // 被举报??

		Rev           Time2 `json:"rev"`                 // 查询时间
		CreateDate    Time1 `json:"createDate"`          // 创建时间
		LastOpTime    Time3 `json:"lastOpTime"`          // 修改时间
		LastOpTimeStr Time1 `json:"lastOpTimeStr"`       // 修改时间
		ShootTime     Time1 `json:"shootTime,omitempty"` // 拍摄时间
	}
)

type UserCloudStorageData struct {
	ID string `json:"id"` // 数据库ID

	FolderID string `json:"folderId"` // 目录ID
	Name     string `json:"name"`     // 名称

	CatalogueType int    `json:"catalogueType"` // 分类 (1:文件|私密空间，2:相册,4:视频,5:音频,6:文档,9:同步盘)
	CloudType     int    `json:"cloudType"`     // (1:普通文件，2:私密空间)
	FileOrAlbum   int    `json:"fileOrAlbum"`   // 文件或相册(1:文件,2:相册)
	Path          string `json:"path"`          // 路径

	UserID     string `json:"userId"`     // 用户ID
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
	IsDel      Bool   `json:"isDel"`      // 软删除
}

// 查询所有用户空间
func (c *MoClient) QueryUserCloudStorage(option ...RestyOption) ([]UserCloudStorageData, error) {
	var resp []UserCloudStorageData
	_, err := c.Request(MoPanProxyFamily+"/user/cloudStorage/getByUserId", Json{
		"type": 2,
	}, &resp, option...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type SearchFileData struct {
	Count      int    `json:"count"`
	FolderList Folder `json:"folderList"`
	FileList   File   `json:"fileList"`
}

// 搜索文件、文件夹
func (c *MoClient) SearchFiles(filename, folderId string, pageNum int, paramOption []ParamOption, option ...RestyOption) (*SearchFileData, error) {
	param := Json{
		"filename":  filename,
		"folderId":  folderId,
		"recursive": 1, // 递归查找
		"pageNum":   pageNum,

		"source": 1, // (个人:1,分享:2)
		"type":   1,
	}
	QueryFileOptionFileType(AllType)(param)
	QueryFileOptionSort(S_FileName, true)(param)
	QueryFileOptionPageSize(1024)(param)

	QueryFileOptionMediaAttr(true)(param)
	QueryFileOptionIconOption(1)(param)

	ApplyParamOption(param, paramOption...)

	var resp SearchFileData
	_, err := c.Request(MoPanProxyFamily+"/file/searchFiles", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type QueryFilesData struct {
	LastRev    Time2 `json:"lastRev"` // 查询时间
	FileListAO struct {
		FileListSize int      `json:"fileListSize"` // 返回文件(非文件夹)数量
		Count        int      `json:"count"`        // 所有文件/文件夹数量
		FolderList   []Folder `json:"folderList"`
		FileList     []File   `json:"fileList"`
	} `json:"fileListAO"`
}

// 查询文件、文件夹
// 当folderId == -11时列出所有文件夹
func (c *MoClient) QueryFiles(folderId string, pageNum int, paramOption []ParamOption, option ...RestyOption) (*QueryFilesData, error) {
	param := Json{
		"folderId": folderId,              // 目录ID
		"pageNum":  strconv.Itoa(pageNum), // 页数

		"source": 1, // (个人:1,分享:2)
		"type":   1,

		"remark": 60,
	}

	QueryFileOptionFileType(AllType)(param)
	QueryFileOptionSort(S_FileName, true)(param)
	QueryFileOptionPageSize(1024)(param)

	QueryFileOptionMediaAttr(true)(param)
	QueryFileOptionIconOption(1)(param)

	ApplyParamOption(param, paramOption...)

	var resp QueryFilesData
	_, err := c.Request(MoPanProxyFamily+"/file/listFiles", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 获取文件详细信息
// 可用选项 QueryFileOptionIconOption QueryFileOptionMediaAttr QueryFileOptionShareFile
func (c *MoClient) GetFileInfo(fileId string, paramOption []ParamOption, option ...RestyOption) (*FileInfo, error) {
	param := Json{
		"fileId": fileId,

		"source": 1, // (个人:1,共享:2)
		"type":   1,
	}

	QueryFileOptionMediaAttr(true)(param)
	QueryFileOptionIconOption(1)(param)

	ApplyParamOption(param, paramOption...)

	var file FileInfo
	_, err := c.Request(MoPanProxyFamily+"/file/getFileInfo", param, &file, option...)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// 获取文件详细信息
func (c *MoClient) GetFolderInfo(folderId string, paramOption []ParamOption, option ...RestyOption) (*FolderInfo, error) {
	param := Json{
		"folderId": folderId,

		"source": 1, // (个人:1,共享:2)
		"type":   1,
	}
	ApplyParamOption(param, paramOption...)

	var folder FolderInfo
	_, err := c.Request(MoPanProxyFamily+"/file/getFolderInfo", param, &folder, option...)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// 创建文件夹
func (c *MoClient) CreateFolder(folderName, parentFolderId string, paramOption []ParamOption, option ...RestyOption) (*Folder, error) {
	param := Json{
		"folderName":     folderName,
		"parentFolderId": parentFolderId,
		"relativePath":   "/", // 相当于MkdirAll(/{relativePath}/{folderName})
		"source":         1,
	}
	ApplyParamOption(param, paramOption...)

	var resp Folder
	_, err := c.Request(MoPanProxyFamily+"/file/createFolder", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 创建文件夹（不检查是否存在）
func (c *MoClient) CreateFolderNotCheck(folderName, parentFolderId string, paramOption []ParamOption, option ...RestyOption) (*Folder, error) {
	param := Json{
		"folderName":     folderName,
		"parentFolderId": parentFolderId,
		"relativePath":   "/", // 相当于MkdirAll(/{relativePath}/{folderName})
		"source":         1,
	}
	ApplyParamOption(param, paramOption...)

	var resp Folder
	_, err := c.Request(MoPanProxyFamily+"/file/createFolderNotCheck", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RenameFileData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Md5  string `json:"md5"`
	Size Int64  `json:"size"`

	FileCata  int `json:"fileCata"`
	MediaType int `json:"mediaType"`

	Rev        Time2 `json:"rev"`        // 查询时间
	CreateDate Time1 `json:"createDate"` // 创建时间
	LastOpTime Time1 `json:"lastOpTime"` // 修改时间
}

// 重命名文件
func (c *MoClient) RenameFile(fileID, fileName string, paramOption []ParamOption, option ...RestyOption) (*RenameFileData, error) {
	param := Json{
		"fileId":       fileID,
		"destFileName": fileName,
		"source":       1,
	}

	ApplyParamOption(param, paramOption...)

	var resp RenameFileData
	_, err := c.Request(MoPanProxyFamily+"/file/renameFile", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RenameFolderData struct {
	ID       string `json:"id"`
	ParentID string `json:"parentId"`
	Name     string `json:"name"`

	FileListSize int `json:"fileListSize"`
	FileCata     int `json:"fileCata"`

	Rev        Time2 `json:"rev"`        // 查询时间
	CreateDate Time1 `json:"createDate"` // 创建时间
	LastOpTime Time1 `json:"lastOpTime"` // 修改时间
}

// 重命名文件夹
func (c *MoClient) RenameFolder(fileID, fileName string, paramOption []ParamOption, option ...RestyOption) (*RenameFolderData, error) {
	param := Json{
		"folderId":       fileID,
		"destFolderName": fileName,
		"source":         1,
	}

	ApplyParamOption(param, paramOption...)

	var resp RenameFolderData
	_, err := c.Request(MoPanProxyFamily+"/file/renameFolder", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type DeleteTaskData struct {
	NotDeleteFileIds []string `json:"notDeleteFileIds"`
	TaskBaseInfo
}

// 删除文件到回收站
func (c *MoClient) DeleteToRecycle(file []TaskFileParam, paramOption []ParamOption, option ...RestyOption) (*DeleteTaskData, error) {
	param := Json{
		"source":    1,
		"type":      1,
		"taskInfos": file,
	}
	ApplyParamOption(param, paramOption...)

	var resp DeleteTaskData
	_, err := c.Request(MoPanProxyFamily+"/recycle/deleteToRecycle", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetFileDownloadUrlData struct {
	DownloadUrl string `json:"downloadUrl"`
	FileID      string `json:"fileId"`
}

// 获取下载链接
func (c *MoClient) GetFileDownloadUrl(fileID string, paramOption []ParamOption, option ...RestyOption) (*GetFileDownloadUrlData, error) {
	param := Json{
		"fileId":    fileID,
		"forcedGet": 0, // 强制获取
		"ifShort":   false,
		"limitRate": "10485760",
		"source":    1,
		//"shareId":   "",
	}

	ApplyParamOption(param, paramOption...)

	var resp GetFileDownloadUrlData
	_, err := c.Request(MoPanProxyFamily+"/file/getFileDownloadUrl", param, &resp, option...)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type (
	TaskType string

	TaskFileParam struct {
		FileID   string `json:"fileId"`
		IsFolder bool   `json:"isFolder"`
		FileName string `json:"fileName"`
	}

	// 批量任务信息
	TaskStatus struct {
		TaskID              string            `json:"taskId"`                       // 任务ID
		TaskStatus          int               `json:"taskStatus"`                   // 任务状态(1-4,4:完成,2:存在冲突)
		SubTaskCount        int               `json:"subTaskCount"`                 // 任务数量
		SuccessedCount      int               `json:"successedCount"`               // 成功数量
		FailedCount         int               `json:"failedCount"`                  // 失败数量
		SkipCount           int               `json:"skipCount"`                    // 跳过数量
		SuccessedFileIDList []string          `json:"successedFileIdList"`          // 成功文件ID
		SuccessedFileIdMap  map[string]string `json:"successedFileIdMap,omitempty"` // only copy
	}

	TaskParam struct {
		UserOrCloudID           string          `json:"userOrCloudId"`       // 用户ID或共享空间ID
		Source                  int             `json:"source"`              // (1:用户,2:共享)
		TaskType                TaskType        `json:"taskType"`            // 操作
		TargetUserOrCloudID     string          `json:"targetUserOrCloudId"` // 目标用户ID或共享空间ID
		TargetSource            int             `json:"targetSource"`        // 目标 (1:用户,2:共享)
		TargetType              int             `json:"targetType"`          // 1
		TargetFolderID          string          `json:"targetFolderId"`      // 目标文件夹
		TaskStatusDetailDTOList []TaskFileParam `json:"taskStatusDetailDTOList"`
	}

	TaskCheckParam struct {
		TaskId              string   `json:"taskId"`
		TaskType            TaskType `json:"taskType"`            // 操作
		TargetType          int      `json:"targetType"`          // 1
		TargetFolderID      string   `json:"targetFolderId"`      // 目标文件夹
		TargetSource        int      `json:"targetSource"`        // 目标 (1:用户,2:共享)
		TargetUserOrCloudID string   `json:"targetUserOrCloudId"` // 目标用户ID或共享空间ID
	}

	TaskManageFileParam struct {
		DealWay int `json:"dealWay"` // 处理方法（1：跳过，2：同时存在，3：覆盖）
		TaskFileParam
	}

	TaskManageParam struct {
		TaskID         string                `json:"taskId"`
		TaskType       TaskType              `json:"taskType"`       // 操作
		TargetFolderID string                `json:"targetFolderId"` // 目标文件夹
		Type           int                   `json:"type"`           // 1
		TaskInfos      []TaskManageFileParam `json:"taskInfos"`
	}

	TaskBaseInfo struct {
		TaskIDList []string `json:"taskIdList"`
		TaskType   TaskType `json:"taskType"`
	}

	TaskData struct {
		TargetFolderID string `json:"targetFolderId"`
		TaskBaseInfo
	}
)

const (
	TASK_COPY TaskType = "COPY"
	TASK_MOVE TaskType = "MOVE"

	TASK_SHARE_SAVE TaskType = "SHARE_SAVE"

	TASK_DELETE TaskType = "DELETE"
)

// 检查批量任务状态
func (c *MoClient) CheckBatchTask(task TaskCheckParam, option ...RestyOption) (*TaskStatus, error) {
	var param Json
	data, err := c.Client.JSONMarshal(task)
	if err != nil {
		return nil, err
	}
	c.Client.JSONUnmarshal(data, &param)

	var resp TaskStatus
	_, err = c.Request(MoPanProxyFamily+"/task/status/checkBatchTask", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 增加批量任务
func (c *MoClient) AddBatchTask(task TaskParam, option ...RestyOption) (*TaskData, error) {
	var param Json
	data, err := c.Client.JSONMarshal(task)
	if err != nil {
		return nil, err
	}
	c.Client.JSONUnmarshal(data, &param)

	var resp TaskData
	_, err = c.Request(MoPanProxyFamily+"/task/status/addTask", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 取消批量任务
func (c *MoClient) CancelBatchTask(taskId string, taskType TaskType, option ...RestyOption) error {
	_, err := c.Request(MoPanProxyFamily+"/task/status/cancelBatchTask", Json{
		"taskId":   taskId,
		"taskType": taskType,
	}, nil, option...)
	if err != nil {
		return err
	}
	return nil
}

type GetConflictTaskInfoData struct {
	TaskID    string `json:"taskId"`
	TaskInfos string `json:"taskInfos"`
}

// 获取任务冲突信息
func (c *MoClient) GetConflictTaskInfo(taskID string, taskType TaskType, option ...RestyOption) (*GetConflictTaskInfoData, error) {
	param := Json{
		"taskId":   taskID,
		"taskType": taskType,
	}

	var resp GetConflictTaskInfoData
	_, err := c.Request(MoPanProxyFamily+"/task/status/getConflictTaskInfo", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// 处理冲突任务
func (c *MoClient) ManageBatchTask(task TaskManageParam, option ...RestyOption) error {
	var param Json
	data, err := c.Client.JSONMarshal(task)
	if err != nil {
		return err
	}
	c.Client.JSONUnmarshal(data, &param)
	_, err = c.Request(MoPanProxyFamily+"/task/status/manageBatchTask", param, nil, option...)
	return err
}
