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
		r.POST("/register/code", controller.RegisterCode)
		r.POST("/register", controller.RegisterSubmit)
		// login
		r.POST("/login", controller.Login)

		user := r.Group("/user")
		{
			user.Use(middleware.AuthPermission())
			user.GET("/status", controller.UserStatus)
			user.GET("/info", controller.UserInfo)
			user.PATCH("/password/update", controller.UserPasswordUpdate)
			user.POST("/email/update/code", controller.UserEmailUpdateCode)
			user.PATCH("/email/update", controller.UserEmailUpdate)
		}

		app := r.Group("/app")
		{
			app.Use(middleware.AuthPermission())
			app.GET("/list", controller.AppGetList)
			app.POST("/create", controller.AppCreate)
			app.PUT("/id/:appId", controller.AppEdit)
			app.DELETE("/id/:appId", controller.AppDel)
			app.GET("/id/:appId", controller.AppInfo)
		}

		forget := r.Group("/forget")
		{
			forget.POST("/password/code", controller.ForgetPasswordCode)
			forget.PATCH("/password/update", controller.ForgetPasswordUpdate)
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
