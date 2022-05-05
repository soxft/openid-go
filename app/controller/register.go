package controller

import (
	"github.com/gin-gonic/gin"
	"openid/queueutil"
)

// RegisterSubmit
// @description send code to email
func RegisterSubmit(c *gin.Context) {
	for i := 0; i < 2; i++ {
		_ = queueutil.Q.Publish("mail", "hello", 0)
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
