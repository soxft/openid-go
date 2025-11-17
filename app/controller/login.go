package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/app/dto"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/userutil"
)

func Login(c *gin.Context) {
	var req dto.LoginRequest
	api := apiutil.New(c)
	
	if err := dto.BindJSON(c, &req); err != nil {
		api.Fail("请求参数错误")
		return
	}

	// check username and password
	if userId, err := userutil.CheckPassword(req.Username, req.Password); err != nil {
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
