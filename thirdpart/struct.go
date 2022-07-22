package thirdpart

import "github.com/gin-gonic/gin"

type ThirdPart interface {
	Handler(c *gin.Context) (string, error)
}
