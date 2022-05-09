package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"openid/library/apiutil"
	"openid/library/codeutil"
	"openid/library/mailutil"
	"openid/library/tool"
	"openid/library/userutil"
	"openid/process/mysqlutil"
	"openid/process/queueutil"
	"time"
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

	// 防止频繁发送验证码
	if beacon, err := mailutil.CheckBeacon(c, email); beacon || err != nil {
		api.Fail("code send too frequently")
		return
	}

	// send mail
	coder := codeutil.New()
	verifyCode := coder.Create(4)
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   verifyCode + " 为您的验证码",
		Content:   "您正在申请找回密码, 您的验证码为: " + verifyCode + ", 有效期10分钟",
		Typ:       "emailChange",
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
	if !tool.IsEmail(email) {
		api.Fail("非法的邮箱格式")
		return
	}

	if !tool.IsPassword(newPassword) {
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
	passwordDb := tool.Sha1(newPassword + salt)
	if res, err := mysqlutil.D.Exec("UPDATE `account` SET `password` = ?, `salt` = ? WHERE `email` = ?", passwordDb, salt, email); err != nil {
		log.Printf("[ERROR] UserPasswordUpdate %v", err)
		api.Fail("system error")
		return
	} else if rows, _ := res.RowsAffected(); rows == 0 {
		api.Fail("用户不存在")
		return
	}

	// send safe notify email
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: c.GetString("email"),
		Subject:   "您的密码已修改",
		Content:   "您的密码已于" + time.Now().Format("2006-01-02 15:04:05") + "修改, 如果不是您本人操作, 请及时联系管理员",
		Typ:       "passwordChangeNotify",
	})
	_ = queueutil.Q.Publish("mail", string(_msg), 5)

	api.Success("修改成功!")

}
