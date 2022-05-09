package middleware

import (
	"github.com/gin-gonic/gin"
	"math"
	"openid/library/apiutil"
	"openid/library/userutil"
	"strconv"
	"time"
)

func AuthPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check timestamp
		timestamp := c.GetHeader("X-timestamp")
		api := apiutil.Api{
			Ctx: c,
		}
		if timestamp == "" {
			api.Abort401("Unauthorized", 0)
			return
		}
		if times, err := strconv.ParseInt(timestamp, 10, 64); err != nil {
			api.Abort401("Unauthorized", 1)
			return
		} else {
			if math.Abs(float64(times-time.Now().Unix())) > 30 {
				api.Abort401("Unauthorized", 2)
				return
			}
		}

		// check jwt token
		var token string
		if token = userutil.GetJwtFromAuth(c.GetHeader("Authorization")); token == "" {
			api.Abort401("Unauthorized", 3)
			return
		}
		if userInfo, err := userutil.CheckJwt(token); err != nil {
			api.Abort401("Unauthorized", 4)
			return
		} else {
			c.Set("userId", userInfo.UserId)
			c.Set("username", userInfo.Username)
			c.Set("email", userInfo.Email)
		}
		c.Next()
	}
}
