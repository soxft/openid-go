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
	userInfo, _ := userutil.GetUserInfoByContext(c)
	api.Out(true, "logon", gin.H{
		"username": userInfo.Username,
		"userId":   userInfo.UserId,
	})
}
