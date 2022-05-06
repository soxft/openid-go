package middleware

import (
	"github.com/gin-gonic/gin"
)

func UserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
