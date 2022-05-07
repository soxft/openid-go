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
	userInfo, _ := c.Get("userInfo")
	type UserInfo = userutil.UserInfo
	userLast := userutil.GetUserLast(userInfo.(UserInfo).UserId)
	api.Out(true, "userInfo", gin.H{
		"userId":   userInfo.(UserInfo).UserId,
		"username": userInfo.(UserInfo).Username,
		"email":    userInfo.(UserInfo).Email,
		"lastTime": userLast.LastTime,
	})
}
