package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid/library/apiutil"
	"time"
)

func Ping(c *gin.Context) {
	api := apiutil.New(c)
	api.SuccessWithData("pong", time.Now().UnixNano())
}
