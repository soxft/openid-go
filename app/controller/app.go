package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
)

func CreateApp(c *gin.Context) {

}

func EditApp(c *gin.Context) {

}

func DelApp(c *gin.Context) {

}

// AppGetList
// @desc 获取用户app列表
func AppGetList(c *gin.Context) {
	userInfo, _ := c.Get("userInfo")
	api := apiutil.New(c)
	api.Success("success", userInfo)

}
