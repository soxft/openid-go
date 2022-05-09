package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"openid/library/apiutil"
	"openid/library/tool"
	"openid/library/userutil"
	"openid/process/mysqlutil"
	"time"
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
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	userId := c.GetInt("userId")
	api := apiutil.New(c)

	if !tool.IsPassword(newPassword) {
		api.Fail("密码应在8～64位")
		return
	}
	// verify old password
	if right, err := userutil.CheckPasswordByUserId(userId, oldPassword); err != nil {
		api.Fail(err.Error())
		return
	} else if !right {
		api.Fail("旧密码错误")
		return
	}

	// change password
	salt := userutil.GenerateSalt()
	passwordDb := tool.Sha1(newPassword + salt)
	if res, err := mysqlutil.D.Exec("UPDATE `account` SET `password` = ?, `salt` = ? WHERE `id` = ?", passwordDb, salt, userId); err != nil {
		log.Printf("[ERROR] UserPasswordUpdate %v", err)
		api.Fail("system error")
		return
	} else if rows, _ := res.RowsAffected(); rows == 0 {
		api.Fail("用户不存在")
		return
	}
	// make jwt token expire
	_ = userutil.SetUserJwtExpire(c.GetString("username"), time.Now().Unix())
	api.Success("修改成功")
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
