package route

import (
	"github.com/gin-gonic/gin"
	"openid/app/controller"
	"openid/app/controller/handler"
	"openid/config"
)

func Init(r *gin.Engine) {
	r.Use(gin.Recovery())
	if config.C.Server.Log {
		r.Use(gin.Logger())
	}
	r.NoRoute(handler.NoRoute)

	// ping
	r.HEAD("/ping", handler.Ping)
	r.GET("/ping", handler.Ping)

	// register
	reg := r.Group("/register")
	{
		reg.POST("/submit", controller.RegisterSubmit)
		reg.POST("/verify", controller.RegisterVerify)
	}

}
