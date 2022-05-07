package apiutil

import "github.com/gin-gonic/gin"

type Apier interface {
	Out(success bool, msg string, data interface{})
	Success(msg string, data interface{})
	Abort(httpCode int, msg string, errorCode int)
	Abort401(msg string, errCode int)
}

type Api struct {
	Ctx *gin.Context
}

func (c *Api) Out(success bool, msg string, data interface{}) {
	c.Ctx.JSON(200, gin.H{
		"success": success,
		"message": msg,
		"data":    data,
	})
}

func (c *Api) Success(msg string, data interface{}) {
	c.Out(true, msg, data)
}

func (c *Api) Abort(httpCode int, msg string, errorCode int) {
	c.Ctx.AbortWithStatusJSON(httpCode, gin.H{
		"success": false,
		"message": msg,
		"data": gin.H{
			"errorCode": errorCode,
		},
	})
}

func (c *Api) Abort401(msg string, errorCode int) {
	c.Abort(401, msg, errorCode)
}
