package middleware

import (
	"github.com/gin-gonic/gin"
	"openid/library/userutil"
)

func UserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if token = userutil.GetJwtFromAuth(c.GetHeader("Authorization")); token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Unauthorized",
				"data":    gin.H{},
			})
			return
		}
		if userInfo, err := userutil.CheckJwt(token); err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Unauthorized",
				"data": gin.H{
					"error": err.Error(),
				},
			})
			return
		} else {
			c.Set("userInfo", userInfo)
		}
		c.Next()
	}
}
