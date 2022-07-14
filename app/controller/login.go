package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid/library/apiutil"
	"github.com/soxft/openid/library/userutil"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	api := apiutil.New(c)
	if len(username) == 0 || len(password) == 0 {
		api.Fail("用户名或密码不能为空")
		return
	}

	// check username and password
	if userId, err := userutil.CheckPassword(username, password); err != nil {
		api.Fail(err.Error())
		return
	} else {
		// get token
		if token, err := userutil.GenerateJwt(userId, c.ClientIP()); err != nil {
			api.Fail("system error")
		} else {
			api.SuccessWithData("登录成功", gin.H{
				"token": token,
			})
		}
	}
}
