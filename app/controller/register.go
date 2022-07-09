package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"openid/config"
	"openid/library/apiutil"
	"openid/library/codeutil"
	"openid/library/mailutil"
	"openid/library/toolutil"
	"openid/library/userutil"
	"openid/process/dbutil"
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
	if !toolutil.IsEmail(email) {
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
	} else if exists {
		api.Fail("email already exists")
		return
	}

	// send Code
	coder := codeutil.New()
	verifyCode := coder.Create(4)
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   verifyCode + " 为您的验证码",
		Content:   "您正在注册 " + config.Server.Title + ". 您的验证码为: " + verifyCode + ", 有效期10分钟.",
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
// @route POST /register
func RegisterSubmit(c *gin.Context) {
	email := c.PostForm("email")
	verifyCode := c.PostForm("code")
	username := c.PostForm("username")
	password := c.PostForm("password")

	api := apiutil.New(c)
	// 合法检测
	if !toolutil.IsEmail(email) {
		api.Fail("非法的邮箱")
		return
	}
	if !toolutil.IsUserName(username) {
		api.Fail("非法的用户名")
		return
	}
	if !toolutil.IsPassword(password) {
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
	if err := userutil.RegisterCheck(username, email); err != nil {
		if err == userutil.ErrUsernameExists {
			api.Fail("用户名已存在")
			return
		} else if err == userutil.ErrEmailExists {
			api.Fail("邮箱已存在")
			return
		}
		api.Fail("server error")
		return
	}

	// 创建用户
	userIp := c.ClientIP()
	timestamp := time.Now().Unix()

	salt := userutil.GenerateSalt()
	pwd := toolutil.Sha1(password + salt)

	// insert to Database
	newUser := dbutil.Account{
		Username: username,
		Password: pwd,
		Salt:     salt,
		Email:    email,
		RegTime:  timestamp,
		RegIp:    userIp,
		LastTime: timestamp,
		LastIp:   userIp,
	}
	result := dbutil.D.Create(&newUser)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Printf("[ERROR] RegisterSubmit %s", result.Error.Error())
		api.Fail("register failed")
		return
	}

	// 消费验证码
	coder.Consume("register", email)
	api.Success("success")
}
