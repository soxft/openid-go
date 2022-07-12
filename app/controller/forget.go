package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"openid/library/apiutil"
	"openid/library/codeutil"
	"openid/library/mailutil"
	"openid/library/toolutil"
	"openid/library/userutil"
	"openid/process/dbutil"
	"openid/process/queueutil"
	"time"
)

// ForgetPasswordCode
// @description 忘记密码发送邮件
// @router POST /forget/password/code
func ForgetPasswordCode(c *gin.Context) {
	email := c.PostForm("email")
	api := apiutil.New(c)
	if !toolutil.IsEmail(email) {
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

	// 防止频繁发送验证码
	if beacon, err := mailutil.CheckBeacon(c, email); beacon || err != nil {
		api.Fail("code send too frequently")
		return
	}

	// send mail
	coder := codeutil.New()
	verifyCode := coder.Create(6)
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   verifyCode + " 为您的验证码",
		Content:   "您正在申请找回密码, 您的验证码为: " + verifyCode + ", 有效期10分钟",
		Typ:       "forgetPwd",
	})

	if err := coder.Save("forgetPwd", email, verifyCode, 60*10); err != nil {
		api.Out(false, "send code failed", gin.H{})
		return
	}
	if err := queueutil.Q.Publish("mail", string(_msg), 0); err != nil {
		coder.Consume("forgetPwd", email) // 删除code
		api.Fail("send code failed")
		return
	}
	_ = mailutil.CreateBeacon(c, email, 120)

	api.Success("success")
}

// ForgetPasswordUpdate
// @description 忘记密码重置
// @router PATCH /forget/password/update
func ForgetPasswordUpdate(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")
	newPassword := c.PostForm("new_password")

	api := apiutil.New(c)
	if !toolutil.IsEmail(email) {
		api.Fail("非法的邮箱格式")
		return
	}

	if !toolutil.IsPassword(newPassword) {
		api.Fail("密码应在8-64位之间")
		return
	}

	// verify code
	coder := codeutil.New()
	if pass, err := coder.Check("forgetPwd", email, code); !pass || err != nil {
		api.Fail("验证码错误或已过期")
		return
	}

	// update password
	salt := userutil.GenerateSalt()
	passwordDb := toolutil.Sha1(newPassword + salt)

	// get Username by email
	var username string
	err := dbutil.D.Model(dbutil.Account{}).Where(dbutil.Account{Email: email}).Select("username").Take(&username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 系统中不存在该邮箱
		api.Fail("验证码错误或已过期")
		return
	} else if err != nil {
		log.Printf("[ERROR] GetUsername SQL: %s", err)
		api.Fail("system error")
		return
	}

	err = dbutil.D.Model(dbutil.Account{}).Where(dbutil.Account{Email: email}).Updates(&dbutil.Account{Password: passwordDb, Salt: salt}).Error
	if err != nil {
		log.Printf("[ERROR] UserPasswordUpdate %v", err)
		api.Fail("system error")
		return
	}
	coder.Consume("forgetPwd", email)

	// 修改密码后续安全操作
	_ = userutil.SetUserJwtExpire(username, time.Now().Unix())
	userutil.PasswordChangeNotify(email, time.Now())

	api.Success("修改成功!")
}
