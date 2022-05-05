package handler

import (
	"github.com/gin-gonic/gin"
	"time"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"success":   true,
		"message":   "pong",
		"timestamp": time.Now().Unix(),
	})
}
