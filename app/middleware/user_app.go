package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/apputil"
)

// UserApp 用来检测是否为用户APP
func UserApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		api := apiutil.New(c)

		appID := c.Param("appid")
		if appID == "" {
			api.Abort200("app id is empty", "middleware.user_app.app_id_empty")
			return
		}

		var userID int
		if userID = c.GetInt("userId"); userID == 0 {
			api.Abort401("Unauthorized", "middleware.user_app.user_id_empty")
			return
		}

		if i, err := apputil.CheckIfUserApp(appID, userID); err != nil {
			//log.Printf("check if user app error: %v", err)

			api.Abort401("Unauthorized", "middleware.user_app.error.not_user_app")
			return
		} else if !i {
			api.Abort200("Unauthorized", "middleware.user_app.not_user_app")
			return
		}

		c.Next()
	}
}
