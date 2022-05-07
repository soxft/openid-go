package controller

import (
	"github.com/gin-gonic/gin"
	"openid/library/apiutil"
	"openid/library/apputil"
	"strconv"
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
	pageTmp := c.DefaultQuery("page", "1")
	limitTmp := c.DefaultQuery("limit", "10")
	api := apiutil.New(c)

	var err error
	var page, limit int
	if page, err = strconv.Atoi(pageTmp); err != nil || page < 1 {
		api.Fail("参数错误")
		return
	}
	if limit, err = strconv.Atoi(limitTmp); err != nil || limit < 1 {
		api.Fail("参数错误")
		return
	}
	offset := (page - 1) * limit
	// 获取用户Id
	userId := c.GetInt("userId")

	// 获取用户app数量
	var appCounts int
	if appCounts, err = apputil.GetUserAppCount(userId); err != nil {
		api.Fail("server error")
		return
	}
	if appCounts == 0 {
		api.Fail("没有数据")
		return
	}

	// 获取用户app列表
	var appList []apputil.AppBaseStruct
	if appList, err = apputil.GetUserAppList(userId, limit, offset); err != nil {
		api.FailWithData("获取失败", gin.H{
			"err": err.Error(),
		})
		return
	}
	if len(appList) == 0 {
		api.Fail("当页无数据")
		return
	}

	api.Success("success", gin.H{
		"counts": appCounts,
		"list":   appList,
	})
}
