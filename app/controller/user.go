package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/tool"
	"openid/library/userutil"
)

func UserStatus(c *gin.Context) {
	api := tool.ApiController{
		Ctx: c,
	}
	userInfo, _ := c.Get("userInfo")
	api.Out(true, "logon", gin.H{
		"username": userInfo.(userutil.UserInfo).Username,
		"userId":   userInfo.(userutil.UserInfo).UserId,
	})
}
