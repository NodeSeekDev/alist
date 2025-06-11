package middlewares

import (
	"github.com/friedHDD/alist/v3/internal/conf"
	"github.com/friedHDD/alist/v3/internal/errs"
	"github.com/friedHDD/alist/v3/internal/setting"
	"github.com/friedHDD/alist/v3/server/common"
	"github.com/gin-gonic/gin"
)

func SearchIndex(c *gin.Context) {
	mode := setting.GetStr(conf.SearchIndex)
	if mode == "none" {
		common.ErrorResp(c, errs.SearchNotAvailable, 500)
		c.Abort()
	} else {
		c.Next()
	}
}
