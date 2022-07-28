package version_one

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/apputil"
)

// AppInfo
// @description 获取应用信息
// @route GET /v1/app/info/:appId
func AppInfo(c *gin.Context) {
	appId := c.Param("appid")
	api := apiutil.New(c)

	// get app info
	if appInfo, err := apputil.GetAppInfo(appId); err != nil {
		if err == apputil.ErrAppNotExist {
			api.Fail("app not exist")
			return
		}
		api.Fail("system error")
		return
	} else {
		api.SuccessWithData("success", gin.H{
			"id":      appInfo.Id,
			"name":    appInfo.AppName,
			"gateway": appInfo.AppGateway,
		})
	}
}
