package route

import (
	"github.com/gin-gonic/gin"
	"openid/api/version_one"
	"openid/app/controller"
	"openid/app/middleware"
	"openid/config"
)

func Init(r *gin.Engine) {
	r.Use(gin.Recovery())
	if config.C.Server.Log {
		r.Use(gin.Logger())
	}
	r.Use(middleware.Cors())
	{
		// ping
		{
			r.HEAD("/ping", controller.Ping)
			r.GET("/ping", controller.Ping)
		}

		// register
		reg := r.Group("/register")
		{
			reg.POST("/sendCode", controller.RegisterSendCode)
			reg.POST("/submit", controller.RegisterSubmit)
		}
		// login
		r.POST("/login", controller.Login)

		user := r.Group("/user")
		{
			user.Use(middleware.AuthPermission())
			user.POST("/status", controller.UserStatus)
			user.POST("/info", controller.UserInfo)
		}

		v1 := r.Group("/v1")
		{
			v1.GET("/login", version_one.Login)
			v1.POST("/info", version_one.Info)
			v1.POST("/login", version_one.LoginHandler)
		}

		r.NoRoute(noRoute)
	}
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"message": "Route not exists",
	})
}
