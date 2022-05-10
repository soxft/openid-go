package version_one

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
	"openid/library/apputil"
	"strconv"
)

// AppInfo
// @description 获取应用信息
// @route GET /v1/app/info/:appId
func AppInfo(c *gin.Context) {
	appId := c.Param("appid")
	api := apiutil.New(c)

	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		api.Fail("appId is not a valid number")
		return
	}
	// get app info
	if appInfo, err := apputil.GetAppInfo(appIdInt); err != nil {
		if err == apputil.ErrAppNotExist {
			api.Fail("应用不存在")
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
