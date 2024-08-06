package apiutil

import "github.com/gin-gonic/gin"

func New(ctx *gin.Context) *Api {
	return &Api{
		Ctx: ctx,
	}
}

func (c *Api) Out(httpCode int, success bool, msg string, data interface{}) {
	c.Ctx.JSON(httpCode, gin.H{
		"success": success,
		"message": msg,
		"data":    data,
	})
}

func (c *Api) Success(msg string) {
	c.Out(200, true, msg, gin.H{})
}

func (c *Api) SuccessWithData(msg string, data interface{}) {
	c.Out(200, true, msg, data)
}

func (c *Api) Fail(msg string) {
	c.Out(200, false, msg, gin.H{})
}

func (c *Api) FailWithData(msg string, data interface{}) {
	c.Out(200, false, msg, data)
}

// httpCode

func (c *Api) FailWithHttpCode(httpCode int, msg string) {
	c.Out(httpCode, false, msg, gin.H{})
}

// Abort

func (c *Api) Abort(httpCode int, msg string, errors string) {
	c.Ctx.AbortWithStatusJSON(httpCode, gin.H{
		"success": false,
		"message": msg,
		"data": gin.H{
			"error": errors,
		},
	})
}

func (c *Api) Abort401(msg string, errors string) {
	c.Abort(401, msg, errors)
}

func (c *Api) Abort200(msg string, errors string) {
	c.Abort(200, msg, errors)
}
