package route

import (
	"github.com/gin-gonic/gin"
	"openid/app/controller"
	"openid/config"
)

func Init(r *gin.Engine) {
	r.Use(gin.Recovery())
	if config.C.Server.Log {
		r.Use(gin.Logger())
	}
	r.NoRoute(noRoute)

	// ping
	r.HEAD("/ping", controller.Ping)
	r.GET("/ping", controller.Ping)

	// register
	reg := r.Group("/register")
	{
		reg.POST("/sendCode", controller.RegisterSendCode)
		reg.POST("/submit", controller.RegisterSubmit)
	}
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"message": "Route not exists",
	})
}
