package version_one

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"openid/api/version_one/helper"
	"openid/library/apiutil"
	"openid/library/apputil"
	"strconv"
)

// Code
// @description 处理登录xhr请求 获取code并跳转到redirect_uri
// @route GET /v1/code
func Code(c *gin.Context) {
	appId := c.PostForm("appid")
	redirectUri := c.PostForm("redirect_uri")

	api := apiutil.New(c)
	if appId == "" || redirectUri == "" {
		api.Fail("appid or redirect_uri is empty")
		return
	}

	appIdInt, err := strconv.Atoi(appId)
	if err != nil {
		api.Fail("appId is not a valid number")
		return
	}

	// 检测 redirect_uri 是否为app所对应的
	var redirectUriDomain *url.URL
	if redirectUriDomain, err = url.Parse(redirectUri); err != nil {
		api.Fail("redirect_uri is invalid")
		return
	} else if redirectUriDomain.Host == "" {
		api.Fail("redirect_uri is invalid")
		return
	}

	// get app Info
	var appInfo apputil.AppFullInfoStruct
	if appInfo, err = apputil.GetAppInfo(appIdInt); err != nil {
		if err == apputil.ErrAppNotExist {
			api.Fail("应用不存在")
			return
		}
		api.Fail("system error")
		return
	}

	// 获取 app gateway
	appGateWay := appInfo.AppGateway
	if appGateWay == "" {
		api.Fail("appGateWay is empty, setting it first")
		return
	}

	// 判断是否一致
	if appGateWay != redirectUriDomain.Host {
		api.FailWithData("redirect_uri is not match with appGateWay", gin.H{
			"legal": appGateWay,
			"given": redirectUriDomain.Host,
		})
		return
	}
	token, err := helper.GenerateToken(appIdInt, c.GetInt("userId"))
	if err != nil {
		log.Printf("[ERROR] get app info error: %s", err.Error())
		api.Fail("system error")
		return
	}

	api.SuccessWithData("success", gin.H{
		"token":       token,
		"redirect_to": redirectUri + "?token=" + token,
	})
}
