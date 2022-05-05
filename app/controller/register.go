package controller

import "github.com/gin-gonic/gin"

// RegisterSubmit
// @description send code to email
func RegisterSubmit(c *gin.Context) {
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
