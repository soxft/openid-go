package handler

import "github.com/gin-gonic/gin"

func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"message": "Route not exists",
	})
}
