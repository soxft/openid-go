package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/codeutil"
	"github.com/soxft/openid-go/library/mailutil"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/library/userutil"
	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/queueutil"
	"log"
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
	api.SuccessWithData("success", gin.H{
		"userId":   userId,
		"username": c.GetString("username"),
		"email":    c.GetString("email"),
		"lastTime": c.GetInt64("lastTime"),
	})
}

// UserLogout
// @description 用户退出
func UserLogout(c *gin.Context) {
	api := apiutil.New(c)
	_ = userutil.SetJwtExpire(c.GetString("token"))
	api.Success("success")
}

// UserPasswordUpdate
// @description 修改用户密码
// @router PATCH /user/password/update
func UserPasswordUpdate(c *gin.Context) {
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	userId := c.GetInt("userId")
	api := apiutil.New(c)

	if !toolutil.IsPassword(newPassword) {
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
	passwordDb := toolutil.Sha1(newPassword + salt)

	result := dbutil.D.Model(model.Account{}).Where(&model.Account{ID: userId}).Updates(&model.Account{Password: passwordDb, Salt: salt})

	if result.Error != nil {
		log.Printf("[ERROR] UserPasswordUpdate %v", result.Error)
		api.Fail("system error")
		return
	} else if result.RowsAffected == 0 {
		api.Fail("用户不存在")
		return
	}

	// make jwt token expire
	_ = userutil.SetJwtExpire(c.GetString("token"))

	// send safe notify email
	userutil.PasswordChangeNotify(c.GetString("email"), time.Now())

	api.Success("修改成功, 请重新登录")
}

// UserEmailUpdateCode
// @description 发送邮箱验证码
// @router POST /user/email/update/code
func UserEmailUpdateCode(c *gin.Context) {
	password := c.PostForm("password")
	newEmail := c.PostForm("new_email")
	userId := c.GetInt("userId")
	api := apiutil.New(c)
	if !toolutil.IsEmail(newEmail) {
		api.Fail("非法的邮箱格式")
		return
	}

	// verify old password
	if right, err := userutil.CheckPasswordByUserId(userId, password); err != nil {
		api.Fail(err.Error())
		return
	} else if !right {
		api.Fail("旧密码错误")
		return
	}

	if exist, err := userutil.CheckEmailExists(newEmail); err != nil {
		api.Fail("system error")
		return
	} else if exist {
		api.Fail("邮箱已存在")
		return
	}

	// 防止频繁发送验证码
	if beacon, err := mailutil.CheckBeacon(c, newEmail); beacon || err != nil {
		api.Fail("code send too frequently")
		return
	}

	// send mail
	coder := codeutil.New()
	verifyCode := coder.Create(4)
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: newEmail,
		Subject:   verifyCode + " 为您的验证码",
		Content:   "您正在申请修改邮箱, 您的验证码为: " + verifyCode + ", 有效期10分钟",
		Typ:       "emailChange",
	})

	if err := coder.Save("emailChange", newEmail, verifyCode, 60*10); err != nil {
		api.Out(false, "send code failed", gin.H{})
		return
	}
	if err := queueutil.Q.Publish("mail", string(_msg), 0); err != nil {
		coder.Consume("emailChange", newEmail) // 删除code
		api.Fail("send code failed")
		return
	}
	_ = mailutil.CreateBeacon(c, newEmail, 120)

	api.Success("发送成功")
}

// UserEmailUpdate
// @description 修改用户邮箱
// @router PATCH /user/email/update
func UserEmailUpdate(c *gin.Context) {
	newEmail := c.PostForm("new_email")
	code := c.PostForm("code")
	api := apiutil.New(c)
	if !toolutil.IsEmail(newEmail) {
		api.Fail("非法的邮箱格式")
		return
	}

	// verify code
	coder := codeutil.New()
	if pass, err := coder.Check("emailChange", newEmail, code); !pass || err != nil {
		api.Fail("验证码错误或已过期")
		return
	}

	// update email
	userId := c.GetInt("userId") // get userid from middleware
	result := dbutil.D.Model(&model.Account{}).Where(&model.Account{ID: userId}).Update("email", newEmail)
	if result.Error != nil {
		log.Printf("[ERROR] UserEmailUpdate %v", result.Error)
		api.Fail("system error")
		return
	} else if result.RowsAffected == 0 {
		api.Fail("用户不存在")
		return
	}

	coder.Consume("emailChange", newEmail)
	userutil.EmailChangeNotify(c.GetString("email"), time.Now())
	_ = userutil.SetJwtExpire(c.GetString("token"))
	api.Success("修改成功, 请重新登录")
}
