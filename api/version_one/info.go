package version_one

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid/api/version_one/helper"
	"github.com/soxft/openid/library/apiutil"
	"github.com/soxft/openid/library/apputil"
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

	// 判断appId与appSecret是否正确
	if err := apputil.CheckAppSecret(appId, appSecret); err != nil {
		api.Fail(err.Error())
		return
	}

	// 检测token是否正确 并获取userId
	userId, err := helper.GetUserIdByToken(appId, token)
	if err != nil {
		if err == helper.ErrTokenNotExists {
			api.Fail("Token not exists")
			return
		}
		api.Fail(err.Error())
		return
	}
	userIds, err := helper.GetUserIds(appId, userId)
	if err != nil {
		api.Fail(err.Error())
		return
	}
	// delete token
	_ = helper.DeleteToken(appId, token)
	api.SuccessWithData("success", gin.H{
		"openId":   userIds.OpenId,
		"uniqueId": userIds.UniqueId,
	})

}
