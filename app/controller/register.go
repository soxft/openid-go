package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"openid/library/codeutil"
	"openid/library/mailutil"
	"openid/library/tool"
	"openid/library/userutil"
	"openid/process/mysqlutil"
	"openid/process/queueutil"
	"time"
)

// RegisterSendCode
// @description send code to email
func RegisterSendCode(c *gin.Context) {
	email := c.PostForm("email")

	api := &tool.ApiController{
		Ctx: c,
	}
	// verify email by re
	if !tool.IsEmail(email) {
		api.Out(false, "invalid email", gin.H{})
		return
	}

	// 防止频繁发送验证码
	if beacon, err := mailutil.CheckBeacon(c, email); beacon || err != nil {
		api.Out(false, "code send too frequently", gin.H{})
		return
	}
	// check mail exists
	if exists, err := userutil.CheckEmail(email); err != nil {
		api.Out(false, "server error", gin.H{})
		return
	} else {
		if exists {
			api.Out(false, "email already exists", gin.H{})
			return
		}
	}

	// send Code
	coder := &codeutil.VerifyCode{}
	verifyCode := coder.Create(4)
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   "注册验证码",
		Content:   "您的验证码为: " + verifyCode + ", 有效期10分钟",
		Typ:       "register",
	})

	if err := coder.Save("register", email, verifyCode, 60*10); err != nil {
		api.Out(false, "send code failed", gin.H{})
		return
	}
	if err := queueutil.Q.Publish("mail", string(_msg), 0); err != nil {
		coder.Consume("register", email) // 删除code
		api.Out(false, "send code failed", gin.H{})
		return
	}
	_ = mailutil.CreateBeacon(c, email, 120)

	api.Out(true, "code send success", gin.H{})
}

// RegisterSubmit
// @description do register
func RegisterSubmit(c *gin.Context) {
	email := c.PostForm("email")
	verifyCode := c.PostForm("code")
	username := c.PostForm("username")
	password := c.PostForm("password")

	api := &tool.ApiController{
		Ctx: c,
	}
	// 合法检测
	if !tool.IsEmail(email) {
		api.Out(false, "非法的邮箱", gin.H{})
		return
	}
	if !tool.IsUserName(username) {
		api.Out(false, "非法的用户名", gin.H{})
		return
	}
	if !tool.IsPassword(password) {
		api.Out(false, "密码应在6～128位", gin.H{})
		return
	}

	// 验证码检测
	coder := &codeutil.VerifyCode{}
	if pass, err := coder.Check("register", email, verifyCode); !pass || err != nil {
		api.Out(false, "invalid code", gin.H{})
		return
	}

	// 重复检测
	if success, msg := userutil.RegisterCheck(username, email); !success {
		api.Out(false, msg, gin.H{})
		return
	}

	// 创建用户
	userIp := c.ClientIP()
	timestamp := time.Now().Unix()

	salt := userutil.GenerateSalt()
	pwd := tool.Sha1(salt + password)

	// insert
	_db, err := mysqlutil.D.Prepare("INSERT INTO `account` (`username`,`password`,`salt`,`email`,`regTime`,`regIp`,`lastTime`,`lastIp`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		api.Out(false, "register failed", gin.H{})
		return
	}
	_, err = _db.Query(username, pwd, salt, email, timestamp, userIp, timestamp, userIp)
	if err != nil {
		api.Out(false, "register failed", gin.H{})
		return
	}
	coder.Consume("register", email)
	api.Out(true, "success", gin.H{})
}
