package version_one

import (
	"errors"
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/api/version_one/helper"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/apputil"
)

type CodeRequest struct {
	AppId       string `json:"appid" binding:"required"`
	RedirectUri string `json:"redirect_uri" binding:"required"`
}

type CodeResponse struct {
	Token      string `json:"token"`
	RedirectTo string `json:"redirect_to"`
}

// Code
// @description 处理登录xhr请求 获取code并跳转到redirect_uri
// @route GET /v1/code
func Code(c *gin.Context) {
	api := apiutil.New(c)

	var req CodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Fail("Invalid params")
		return
	}

	// 检测 redirect_uri 是否为app所对应的
	var err error
	var redirectUriDomain *url.URL
	if redirectUriDomain, err = url.Parse(req.RedirectUri); err != nil {
		api.Fail("Invalid redirect_uri")
		return
	} else if redirectUriDomain.Host == "" {
		api.Fail("Invalid redirect_uri")
		return
	}

	// get app Info
	var appInfo apputil.AppFullInfoStruct
	if appInfo, err = apputil.GetAppInfo(req.AppId); err != nil {
		if errors.Is(err, apputil.ErrAppNotExist) {
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
	if !apputil.CheckRedirectUriIsMatchUserGateway(redirectUriDomain.Host, appGateWay) {
		api.FailWithData("redirect_uri is not match with appGateWay", gin.H{
			"legal": appGateWay,
			"given": redirectUriDomain.Host,
		})
		return
	}
	token, err := helper.GenerateToken(c, req.AppId, c.GetInt("userId"))
	if err != nil {
		log.Printf("[ERROR] get app info error: %s", err.Error())
		api.Fail("system error")
		return
	}

	api.SuccessWithData("success", CodeResponse{
		Token:      token,
		RedirectTo: req.RedirectUri + "?token=" + token,
	})
}
