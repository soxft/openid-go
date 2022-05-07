package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
	"openid/library/userutil"
)

// UserStatus
// @description 判断用户登录状态
func UserStatus(c *gin.Context) {
	api := apiutil.Api{
		Ctx: c,
	}
	// 中间件中已经处理, 直接输出
	api.Out(true, "logon", gin.H{})
}

// UserInfo
// @description 获取用户信息
func UserInfo(c *gin.Context) {
	api := apiutil.Api{
		Ctx: c,
	}
	userId := c.GetInt("userId")
	userLast := userutil.GetUserLast(userId)
	api.Out(true, "userInfo", gin.H{
		"userId":   userId,
		"username": c.GetString("username"),
		"email":    c.GetString("email"),
		"lastTime": userLast.LastTime,
	})
}
