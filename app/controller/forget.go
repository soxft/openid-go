package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
	"openid/library/tool"
	"openid/library/userutil"
)

// ForgetPasswordCode
// @description 忘记密码发送邮件
// @router POST /forget/password/code
func ForgetPasswordCode(c *gin.Context) {
	email := c.PostForm("email")
	api := apiutil.New(c)
	if !tool.IsEmail(email) {
		api.Fail("非法的邮箱格式")
		return
	}
	if exists, err := userutil.CheckEmailExists(email); err != nil {
		api.Fail("server error")
		return
	} else if !exists {
		api.Fail("邮箱不存在")
		return
	}

	api.Success("success")
}

// ForgetPasswordUpdate
// @description 忘记密码重置
// @router PATCH /forget/password/update
func ForgetPasswordUpdate(c *gin.Context) {

}
