package common

import (
	stdpath "path"

	"github.com/friedHDD/alist/v3/internal/conf"
	"github.com/friedHDD/alist/v3/internal/model"
	"github.com/friedHDD/alist/v3/internal/setting"
	"github.com/friedHDD/alist/v3/internal/sign"
)

func Sign(obj model.Obj, parent string, encrypt bool) string {
	if obj.IsDir() || (!encrypt && !setting.GetBool(conf.SignAll)) {
		return ""
	}
	return sign.Sign(stdpath.Join(parent, obj.GetName()))
}
