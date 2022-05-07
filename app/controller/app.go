package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
)

func AppCreate(c *gin.Context) {

}

func AppEdit(c *gin.Context) {

}

func AppDel(c *gin.Context) {

}

// AppGetList
// @desc 获取用户app列表
func AppGetList(c *gin.Context) {
	api := apiutil.New(c)
	api.Success("success", gin.H{
		"userId": c.GetInt("userId"),
	})
}
