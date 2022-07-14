package version_one

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid/api/version_one/helper"
	"github.com/soxft/openid/library/apiutil"
	"github.com/soxft/openid/library/apputil"
	"log"
	"net/url"
)

// Code
// @description 处理登录xhr请求 获取code并跳转到redirect_uri
// @route GET /v1/code
func Code(c *gin.Context) {
	appId := c.PostForm("appid")
	redirectUri := c.PostForm("redirect_uri")

	api := apiutil.New(c)
	if appId == "" || redirectUri == "" {
		api.Fail("Invalid params")
		return
	}

	// 检测 redirect_uri 是否为app所对应的
	var err error
	var redirectUriDomain *url.URL
	if redirectUriDomain, err = url.Parse(redirectUri); err != nil {
		api.Fail("Invalid redirect_uri")
		return
	} else if redirectUriDomain.Host == "" {
		api.Fail("Invalid redirect_uri")
		return
	}

	// get app Info
	var appInfo apputil.AppFullInfoStruct
	if appInfo, err = apputil.GetAppInfo(appId); err != nil {
		if err == apputil.ErrAppNotExist {
			api.Fail("app not exist")
			return
		}
		api.Fail("system error")
		return
	}

	// 获取 app gateway
	appGateWay := appInfo.AppGateway
	if appGateWay == "" {
		api.Fail("Invalid appGateWay, setting it first")
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
	token, err := helper.GenerateToken(appId, c.GetInt("userId"))
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
