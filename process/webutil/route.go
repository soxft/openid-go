package webutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/api/version_one"
	"github.com/soxft/openid-go/app/controller"
	"github.com/soxft/openid-go/app/middleware"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/apiutil"
)

func initRoute(r *gin.Engine) {
	r.Use(gin.Recovery())
	if config.Server.Log {
		r.Use(gin.Logger())
	}
	r.Use(middleware.Cors())
	{
		{
			// ping
			r.HEAD("/ping", controller.Ping)
			r.GET("/ping", controller.Ping)

			// register
			r.POST("/register/code", controller.RegisterCode)
			r.POST("/register", controller.RegisterSubmit)

			// login
			r.POST("/login", controller.Login)
		}

		user := r.Group("/user")
		{
			user.Use(middleware.AuthPermission())
			user.GET("/status", controller.UserStatus)
			user.GET("/info", controller.UserInfo)
			user.POST("/logout", controller.UserLogout)
			user.PATCH("/password/update", controller.UserPasswordUpdate)
			user.POST("/email/update/code", controller.UserEmailUpdateCode)
			user.PATCH("/email/update", controller.UserEmailUpdate)
		}

		pass := r.Group("/passkey")
		{
			// Login endpoints (no auth required)
			pass.POST("/login/options", controller.PasskeyLoginOptions)
			pass.POST("/login", controller.PasskeyLoginFinish)

			pass.Use(middleware.AuthPermission())
			// Registration endpoints (auth required)
			pass.POST("/register/options", controller.PasskeyRegistrationOptions)
			pass.POST("/register", controller.PasskeyRegistrationFinish)

			// Management endpoints (auth required)
			pass.GET("", controller.PasskeyList)
			pass.DELETE(":id", controller.PasskeyDelete)
		}

		app := r.Group("/app")
		{
			app.Use(middleware.AuthPermission())
			app.GET("/list", controller.AppGetList)
			app.POST("/create", controller.AppCreate)

			// 判断 App 归属中间件
			app.Use(middleware.UserApp())
			app.PUT("/id/:appid", controller.AppEdit)
			app.DELETE("/id/:appid", controller.AppDel)
			app.GET("/id/:appid", controller.AppInfo)

			app.PUT("/id/:appid/secret", controller.AppReGenerateSecret)
		}

		forget := r.Group("/forget")
		{
			forget.POST("/password/code", controller.ForgetPasswordCode)
			forget.PATCH("/password/update", controller.ForgetPasswordUpdate)
		}

		v1 := r.Group("/v1")
		{
			v1.GET("/login", version_one.Login)
			v1.POST("/code", middleware.AuthPermission(), version_one.Code)
			v1.POST("/info", version_one.Info)
			v1.GET("/app/info/:appid", version_one.AppInfo)
		}

		r.NoRoute(noRoute)
	}
}

func noRoute(c *gin.Context) {
	api := apiutil.New(c)

	api.FailWithHttpCode(404, "Route not exists")
}
