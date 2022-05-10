package version_one

import (
	"github.com/gin-gonic/gin"
	"openid/api/version_one/helper"
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

	// 检测token是否正确 并获取userId
	userId, err := helper.GetUserIdByToken(appIdInt, token)
	if err != nil {
		if err == helper.ErrTokenNotExists {
			api.Fail("Token not exists")
			return
		}
		api.Fail(err.Error())
		return
	}
	userIds, err := helper.GetUserIds(appIdInt, userId)
	if err != nil {
		api.Fail(err.Error())
		return
	}
	// delete token
	_ = helper.DeleteToken(appIdInt, token)
	api.SuccessWithData("success", gin.H{
		"openid":   userIds.OpenId,
		"uniqueId": userIds.UniqueId,
	})

}
