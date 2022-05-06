package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/tool"
	"openid/library/userutil"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	api := &tool.ApiController{
		Ctx: c,
	}
	if len(username) == 0 || len(password) == 0 {
		api.Out(false, "用户名或密码不能为空", gin.H{})
		return
	}

	// check username and password
	if success, userId, msg := userutil.CheckPassword(username, password); !success {
		api.Out(false, msg, gin.H{})
		return
	} else {
		// get token
		if token, err := userutil.GenerateJwt(userId); err != nil {
			api.Out(false, "system error", gin.H{})
		} else {
			api.Out(true, "登录成功", gin.H{
				"token": token,
			})
		}
	}
}
