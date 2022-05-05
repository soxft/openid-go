package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"openid/library/mail"
	"openid/queueutil"
)

// RegisterSubmit
// @description send code to email
func RegisterSubmit(c *gin.Context) {
	for i := 0; i < 2; i++ {
		_msg, _ := json.Marshal(mail.MailStruct{
			ToAddress: "code@xcsoft.top",
			Subject:   "123",
			Content:   "你好啊",
		})
		_ = queueutil.Q.Publish("mail", string(_msg), 0)
	}
	c.JSON(200, gin.H{
		"message": "Register",
	})
}

// RegisterVerify
// @description do register
func RegisterVerify(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Register",
	})
}
