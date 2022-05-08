package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
	"openid/library/userutil"
)

// UserStatus
// @description 判断用户登录状态
// @router GET /user/status
func UserStatus(c *gin.Context) {
	api := apiutil.New(c)
	// 中间件中已经处理, 直接输出
	api.Success("logon")
}

// UserInfo
// @description 获取用户信息
// @router GET /user/info
func UserInfo(c *gin.Context) {
	api := apiutil.New(c)

	userId := c.GetInt("userId")
	userLast := userutil.GetUserLast(userId)
	api.SuccessWithData("userInfo", gin.H{
		"userId":   userId,
		"username": c.GetString("username"),
		"email":    c.GetString("email"),
		"lastTime": userLast.LastTime,
	})
}

// UserPasswordUpdate
// @description 修改用户密码
// @router PATCH /user/password/update
func UserPasswordUpdate(c *gin.Context) {

}

// UserEmailUpdateCode
// @description 发送邮箱验证码
// @router POST /user/email/update/code
func UserEmailUpdateCode(c *gin.Context) {

}

// UserEmailUpdate
// @description 修改用户邮箱
// @router PATCH /user/email/update
func UserEmailUpdate(c *gin.Context) {

}
