package route

import (
	"github.com/gin-gonic/gin"
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
			user.GET("/status", controller.UserStatus)
			user.GET("/info", controller.UserInfo)
		}

		app := r.Group("/app")
		{
			app.Use(middleware.AuthPermission())
			app.GET("/list", controller.AppGetList)
			app.POST("/create", controller.AppCreate)
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
