package sdk

const (
	ApiBaseURL = "https://proapi.115.com"
	ApiAuthURL = "https://passportapi.115.com"
)

// Auth API
const (
	ApiAuthDeviceCode = ApiAuthURL + "/open/authDeviceCode"
	ApiQrCodeStatus   = "https://qrcodeapi.115.com/get/status/"
	ApiCodeToToken    = ApiAuthURL + "/open/deviceCodeToToken"
	ApiRefreshToken   = ApiAuthURL + "/open/refreshToken"
)

// User API
const (
	ApiUserInfo = ApiBaseURL + "/open/user/info"
)

// File API
const (
	ApiFsUploadGetToken = ApiBaseURL + "/open/upload/get_token"
	ApiFsUploadInit     = ApiBaseURL + "/open/upload/init"
	ApiFsUploadResume   = ApiBaseURL + "/open/upload/resume"
	ApiFsMkdir          = ApiBaseURL + "/open/folder/add"
	ApiFsGetFiles       = ApiBaseURL + "/open/ufile/files"
	ApiFsGetFolderInfo  = ApiBaseURL + "/open/folder/get_info"
	ApiFsSearchFiles    = ApiBaseURL + "/open/ufile/search"
	ApiFsCopy           = ApiBaseURL + "/open/ufile/copy"
	ApiFsMove           = ApiBaseURL + "/open/ufile/move"
	ApiFsDownURL        = ApiBaseURL + "/open/ufile/downurl"
	ApiFsUpdate         = ApiBaseURL + "/open/ufile/update"
	ApiFsDelete         = ApiBaseURL + "/open/ufile/delete"
	ApiFsRbList         = ApiBaseURL + "/open/rb/list"
	ApiFsRbRevert       = ApiBaseURL + "/open/rb/revert"
	ApiFsRbDelete       = ApiBaseURL + "/open/rb/del"
)
