package version_one

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
	"openid/library/apputil"
	"strconv"
)

// Info
// @description 获取openid 和 uniqueId
// @route POST /v1/info
func Info(c *gin.Context) {
	token := c.PostForm("token")
	appId := c.PostForm("appid")
	appSecret := c.PostForm("app_secret")

	api := apiutil.New(c)

	if token == "" || appId == "" || appSecret == "" {
		api.Fail("Invalid params")
		return
	}

	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		api.Fail("appId is not a valid number")
		return
	}

	// 判断appId与appSecret是否正确
	if err := apputil.CheckAppSecret(appIdInt, appSecret); err != nil {
		api.Fail(err.Error())
		return
	}

}
