package controller

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/app/dto"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/apputil"
	"github.com/soxft/openid-go/process/dbutil"
)

// AppCreate
// @description: 创建应用
// @route: POST /app/create
func AppCreate(c *gin.Context) {
	var req dto.AppCreateRequest
	api := apiutil.New(c)
	
	if err := dto.BindJSON(c, &req); err != nil {
		api.Fail("请求参数错误")
		return
	}
	
	if !apputil.CheckName(req.AppName) {
		api.Fail("应用名称不合法")
		return
	}
	// 创建应用
	if success, err := apputil.CreateApp(c.GetInt("userId"), req.AppName); !success {
		api.Fail(err.Error())
		return
	}

	api.Success("创建应用成功")
}

// AppEdit
// @description: 编辑应用
// @route: PUT /app/:id
func AppEdit(c *gin.Context) {
	var req dto.AppEditRequest
	appId := c.Param("appid")
	api := apiutil.New(c)
	
	if err := dto.BindJSON(c, &req); err != nil {
		api.Fail("请求参数错误")
		return
	}
	
	// 参数合法性检测
	if !apputil.CheckName(req.AppName) {
		api.Fail("应用名称不合法")
		return
	}

	if len(req.AppGateway) == 0 {
		api.Fail("网关不能为空")
		return
	} else if len(req.AppGateway) > 200 {
		api.Fail("网关长度不能超过 200字符")
		return
	}

	var gateWayCount int
	// 检测网关是否合法
	var gateways []string

	for _, gateway := range strings.Split(req.AppGateway, "\n") {
		gateway = strings.TrimSpace(gateway)
		if gateway == "" {
			continue
		}
		if !apputil.CheckGateway(gateway) {
			apiutil.New(c).Fail(fmt.Sprintf("网关 %s 不合法", gateway))
			return
		}
		gateways = append(gateways, gateway)

		if gateWayCount++; gateWayCount > 10 {
			api.Fail("网关数量不能超过 10 个")
			return
		}
	}

	// Do Update
	err := dbutil.D.Model(model.App{}).
		Where(model.App{
			AppId: appId,
		}).
		Updates(model.App{
			AppName:    req.AppName,
			AppGateway: strings.Join(gateways, ","),
		}).Error
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

	// delete
	if success, err := apputil.DeleteUserApp(appId); !success {
		api.Fail(err.Error())
	} else if err != nil {
		api.Fail("system error")
	} else {
		api.Success("删除成功")
	}
}

// AppReGenerateSecret
// @description: 重新生成secret
func AppReGenerateSecret(c *gin.Context) {
	appId := c.Param("appid")

	api := apiutil.New(c)

	// re generate secret
	if newToken, err := apputil.ReGenerateSecret(appId); err != nil {
		log.Printf("[ERROR] ReGenerateSecret error: %s", err)
		api.Fail("re generate secret failed, try again later")
	} else {
		api.SuccessWithData("重置 AppSecret 成功!", gin.H{
			"secret": newToken,
		})
	}
}

// AppInfo
// @description: 获取app详细信息
// GET /app/:id
func AppInfo(c *gin.Context) {
	appId := c.Param("appid")
	api := apiutil.New(c)

	// get app info
	if appInfo, err := apputil.GetAppInfo(appId); err != nil {
		if errors.Is(err, apputil.ErrAppNotExist) {
			api.Fail("应用不存在")
			return
		}
		api.Fail("system error")
		return
	} else {
		appInfo.AppGateway = strings.ReplaceAll(appInfo.AppGateway, ",", "\n")
		api.SuccessWithData("success", appInfo)
	}
}

// AppGetList
// @desc 获取用户app列表
// @route GET /app/list
func AppGetList(c *gin.Context) {
	var req dto.AppListRequest
	api := apiutil.New(c)
	
	// 支持从查询参数获取
	pageTmp := c.DefaultQuery("page", "1")
	limitTmp := c.DefaultQuery("per_page", "10")
	
	var err error
	var page, limit int
	if page, err = strconv.Atoi(pageTmp); err != nil || page < 1 {
		page = 1
	}
	if limit, err = strconv.Atoi(limitTmp); err != nil || limit < 1 || limit > 100 {
		limit = 10
	}
	
	// 尝试从JSON获取，如果有的话会覆盖query参数
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err == nil {
			if req.Page > 0 {
				page = req.Page
			}
			if req.PerPage > 0 && req.PerPage <= 100 {
				limit = req.PerPage
			}
		}
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
