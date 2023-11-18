package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/userutil"
)

func AuthPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		api := apiutil.New(c)

		// check jwt token
		var token string
		if token = userutil.GetJwtFromAuth(c.GetHeader("Authorization")); token == "" {
			api.Abort401("Unauthorized", 0)
			return
		}
		if userInfo, err := userutil.CheckPermission(c, token); err != nil {
			api.Abort401("Unauthorized", 1)
			return
		} else {
			c.Set("userId", userInfo.UserId)
			c.Set("username", userInfo.Username)
			c.Set("email", userInfo.Email)
			c.Set("lastTime", userInfo.LastTime)
			c.Set("token", token)
		}
		c.Next()
	}
}
