package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/config"
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

	// 先创建 beacon 再说
	_ = mailutil.CreateBeacon(c, email, 120)

	// check mail exists
	if exists, err := userutil.CheckEmailExists(email); err != nil {
		go mailutil.DeleteBeacon(c, email) // 删除信标

		api.Fail("server error")
		return
	} else if exists {
		go mailutil.DeleteBeacon(c, email) // 删除信标

		api.Fail("email already exists")
		return
	}

	// send Code
	coder := codeutil.New(c)
	verifyCode := coder.Create(4)

	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   verifyCode + " 为您的验证码",
		Content:   "您正在注册 " + config.Server.Title + ". 您的验证码为: " + verifyCode + ", 有效期10分钟.",
		Typ:       "register",
	})

	if err := coder.Save("register", email, verifyCode, 60*time.Minute); err != nil {
		go mailutil.DeleteBeacon(c, email) // 删除信标

		api.Fail("send code failed")
		return
	}

	if err := queueutil.Q.Publish("mail", string(_msg), 0); err != nil {
		go coder.Consume("register", email) // 删除code
		go mailutil.DeleteBeacon(c, email)  // 删除信标

		api.Fail("send code failed")
		return
	}

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
	coder := codeutil.New(c)
	if pass, err := coder.Check("register", email, verifyCode); !pass || err != nil {
		api.Fail("invalid code")
		return
	}

	// 重复检测
	if err := userutil.RegisterCheck(username, email); err != nil {
		if errors.Is(err, userutil.ErrUsernameExists) {
			api.Fail("用户名已存在")
			return
		} else if errors.Is(err, userutil.ErrEmailExists) {
			api.Fail("邮箱已存在")
			return
		}
		api.Fail("server error")
		return
	}

	// 创建用户
	userIp := c.ClientIP()
	timestamp := time.Now().Unix()

	var err error
	var pwd string
	if pwd, err = userutil.GeneratePwd(password); err != nil {
		log.Println(err)
		api.Fail("pwd generate error")
		return
	}

	// insert to Database
	newUser := model.Account{
		Username: username,
		Password: pwd,
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
