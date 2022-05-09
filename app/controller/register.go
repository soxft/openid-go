package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"openid/config"
	"openid/library/apiutil"
	"openid/library/codeutil"
	"openid/library/mailutil"
	"openid/library/tool"
	"openid/library/userutil"
	"openid/process/mysqlutil"
	"openid/process/queueutil"
	"time"
)

// RegisterCode
// @description send code to email
// @route POST /register/code
func RegisterCode(c *gin.Context) {
	email := c.PostForm("email")

	api := apiutil.New(c)
	// verify email by re
	if !tool.IsEmail(email) {
		api.Fail("invalid email")
		return
	}

	// 防止频繁发送验证码
	if beacon, err := mailutil.CheckBeacon(c, email); beacon || err != nil {
		api.Fail("code send too frequently")
		return
	}
	// check mail exists
	if exists, err := userutil.CheckEmailExists(email); err != nil {
		api.Fail("server error")
		return
	} else {
		if exists {
			api.Fail("email already exists")
			return
		}
	}

	// send Code
	coder := codeutil.New()
	verifyCode := coder.Create(4)
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   verifyCode + " 为您的验证码",
		Content:   "您正在注册 " + config.C.Server.Title + ". 您的验证码为: " + verifyCode + ", 有效期10分钟.",
		Typ:       "register",
	})

	if err := coder.Save("register", email, verifyCode, 60*10); err != nil {
		api.Out(false, "send code failed", gin.H{})
		return
	}
	if err := queueutil.Q.Publish("mail", string(_msg), 0); err != nil {
		coder.Consume("register", email) // 删除code
		api.Fail("send code failed")
		return
	}
	_ = mailutil.CreateBeacon(c, email, 120)

	api.Success("code send success")
}

// RegisterSubmit
// @description do register
// @route POST /register/
func RegisterSubmit(c *gin.Context) {
	email := c.PostForm("email")
	verifyCode := c.PostForm("code")
	username := c.PostForm("username")
	password := c.PostForm("password")

	api := apiutil.New(c)
	// 合法检测
	if !tool.IsEmail(email) {
		api.Fail("非法的邮箱")
		return
	}
	if !tool.IsUserName(username) {
		api.Fail("非法的用户名")
		return
	}
	if !tool.IsPassword(password) {
		api.Fail("密码应在8～64位")
		return
	}

	// 验证码检测
	coder := codeutil.New()
	if pass, err := coder.Check("register", email, verifyCode); !pass || err != nil {
		api.Fail("invalid code")
		return
	}

	// 重复检测
	if success, msg := userutil.RegisterCheck(username, email); !success {
		api.Fail(msg)
		return
	}

	// 创建用户
	userIp := c.ClientIP()
	timestamp := time.Now().Unix()

	salt := userutil.GenerateSalt()
	pwd := tool.Sha1(password + salt)

	// insert
	_db, err := mysqlutil.D.Prepare("INSERT INTO `account` (`username`,`password`,`salt`,`email`,`regTime`,`regIp`,`lastTime`,`lastIp`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] RegisterSubmit %s", err.Error())
		api.Fail("register failed")
		return
	}
	_, err = _db.Query(username, pwd, salt, email, timestamp, userIp, timestamp, userIp)
	if err != nil {
		log.Printf("[ERROR] RegisterSubmit %s", err.Error())
		api.Fail("register failed")
		return
	}
	coder.Consume("register", email)
	api.Success("success")
}
