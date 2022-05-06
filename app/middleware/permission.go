package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"openid/library/userutil"
	"strconv"
	"time"
)

func AuthPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check timestamp
		timestamp := c.GetHeader("X-timestamp")
		if timestamp == "" {
			abortWithStatusJSON(c, "Unauthorized", 0)
			return
		}
		if times, err := strconv.ParseInt(timestamp, 10, 64); err != nil {
			abortWithStatusJSON(c, "Unauthorized", 1)
			return
		} else {
			if math.Abs(float64(times-time.Now().Unix())) > 30 {
				abortWithStatusJSON(c, "Unauthorized", 2)
				return
			}
		}

		// check jwt token
		var token string
		if token = userutil.GetJwtFromAuth(c.GetHeader("Authorization")); token == "" {
			abortWithStatusJSON(c, "Unauthorized", 3)
			return
		}
		if userInfo, err := userutil.CheckJwt(token); err != nil {
			log.Printf("[ERROR] CheckJwt error: %s", err)
			abortWithStatusJSON(c, "Unauthorized", 4)
			return
		} else {
			c.Set("userInfo", userInfo)
		}
		c.Next()
	}
}

func abortWithStatusJSON(c *gin.Context, message string, errorCode int) {
	c.AbortWithStatusJSON(401, gin.H{
		"success": false,
		"message": message,
		"data": gin.H{
			"errorCode": errorCode,
		},
	})
}
