package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"openid/library/code"
	"openid/library/mail"
	"openid/library/tool"
	"openid/queueutil"
)

// RegisterSendCode
// @description send code to email
func RegisterSendCode(c *gin.Context) {
	email := c.PostForm("email")
	// verify email by re
	if !tool.IsEmail(email) {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid email",
			"data":    gin.H{},
		})
		return
	}

	if beacon, err := mail.CheckBeacon(c, email); beacon || err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "code send too frequently",
			"data":    gin.H{},
		})
		return
	}

	// send Code
	coder := &code.VerifyCode{}
	verifyCode := coder.Create(6)
	_msg, _ := json.Marshal(mail.Mail{
		ToAddress: email,
		Subject:   "注册验证码",
		Content:   "您的验证码为: " + verifyCode + ", 有效期10分钟",
	})
	if err := queueutil.Q.Publish("mail", string(_msg), 0); err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "send code failed",
			"data":    gin.H{},
		})
		return
	}

	if err := coder.Save("register", email, verifyCode, 60*10); err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "send code failed",
			"data":    gin.H{},
		})
	}
	_ = mail.CreateBeacon(c, email, 120)

	c.JSON(200, gin.H{
		"success": true,
		"message": "code send success",
		"data":    gin.H{},
	})
}

// RegisterSubmit
// @description do register
func RegisterSubmit(c *gin.Context) {
	email := c.PostForm("email")
	verifyCode := c.PostForm("code")

	coder := &code.VerifyCode{}
	if pass, err := coder.Check("register", email, verifyCode); !pass || err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "invalid code",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "success",
		"data":    gin.H{},
	})
}
