package tool

import "github.com/gin-gonic/gin"

type Api interface {
	Out(router *gin.Engine)
}

type ApiController struct {
	Ctx *gin.Context
}

func (c *ApiController) Out(success bool, msg string, data interface{}) {
	c.Ctx.JSON(200, gin.H{
		"success": success,
		"message": msg,
		"data":    data,
	})
}
