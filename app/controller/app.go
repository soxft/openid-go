package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"openid/library/apiutil"
	"openid/library/apputil"
	"openid/process/dbutil"
	"strconv"
	"strings"
)

// AppCreate
// @description: 创建应用
// @route: POST /app/create
func AppCreate(c *gin.Context) {
	appName := c.PostForm("app_name")
	api := apiutil.New(c)
	if !apputil.CheckName(appName) {
		api.Fail("应用名称不合法")
		return
	}
	// 创建应用
	if success, err := apputil.CreateApp(c.GetInt("userId"), appName); !success {
		api.Fail(err.Error())
		return
	}

	api.Success("创建应用成功")
}

// AppEdit
// @description: 编辑应用
// @route: PUT /app/:id
func AppEdit(c *gin.Context) {
	appId := c.Param("appid")
	appName := c.PostForm("app_name")
	appGateway := c.PostForm("app_gateway")

	appGateway = strings.TrimSpace(appGateway)
	api := apiutil.New(c)

	// 参数合法性检测
	if !apputil.CheckName(appName) {
		api.Fail("应用名称不合法")
		return
	}
	if !apputil.CheckGateway(appGateway) {
		api.Fail("应用网关不合法")
		return
	}
	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		api.Fail("非法访问")
		return
	}

	// 判断是否为 该用户的app
	if i, err := apputil.CheckIfUserApp(appIdInt, c.GetInt("userId")); err != nil {
		api.Fail("system error")
		return
	} else if !i {
		api.Fail("没有权限")
		return
	}

	// do change
	db, err := dbutil.D.Prepare("UPDATE `app` SET `appName` = ?, `appGateway` = ? WHERE `appId` = ?")
	if err != nil {
		log.Printf("[ERROR] mysqlutil.D.Prepare err: %v", err)
		api.Fail("system error")
		return
	}
	_, err = db.Exec(appName, appGateway, appIdInt)
	if err != nil {
		log.Printf("[ERROR] db.Exec err: %v", err)
		api.Fail("system error")
		return
	}
	api.Success("修改成功")
}

// AppDel
// @description: 删除app
// @route: DELETE /app/:id
func AppDel(c *gin.Context) {
	appId := c.Param("appid")
	api := apiutil.New(c)
	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		api.Fail("非法访问 ")
	}

	// 判断是否为 该用户的app
	if i, err := apputil.CheckIfUserApp(appIdInt, c.GetInt("userId")); err != nil {
		api.Fail("system error")
		return
	} else if !i {
		api.Fail("没有权限")
		return
	}

	// delete
	if success, err := apputil.DeleteUserApp(appIdInt); !success {
		api.Fail(err.Error())
	} else if err != nil {
		api.Fail("system error")
	} else {
		api.Success("删除成功")
	}
}

// AppInfo
// @description: 获取app详细信息
// GET /app/:id
func AppInfo(c *gin.Context) {
	appId := c.Param("appid")
	api := apiutil.New(c)
	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		api.Fail("非法访问")
		return
	}
	// 判断是否为 该用户的app
	if i, err := apputil.CheckIfUserApp(appIdInt, c.GetInt("userId")); err != nil {
		api.Fail("system error")
		return
	} else if !i {
		api.Fail("没有权限")
		return
	}

	// get app info
	if appInfo, err := apputil.GetAppInfo(appIdInt); err != nil {
		if err == apputil.ErrAppNotExist {
			api.Fail("应用不存在")
			return
		}
		api.Fail("system error")
		return
	} else {
		api.SuccessWithData("success", appInfo)
	}
}

// AppGetList
// @desc 获取用户app列表
// @route GET /app/list
func AppGetList(c *gin.Context) {
	pageTmp := c.DefaultQuery("page", "1")
	limitTmp := c.DefaultQuery("per_page", "10")
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
		api.SuccessWithData("success", gin.H{
			"total": 0,
			"list":  []gin.H{},
		})
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

	api.SuccessWithData("success", gin.H{
		"total": appCounts,
		"list":  appList,
	})
}
