package version_one

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/apiutil"
	"net/url"
)

// Login
// @description v1 登录
// @route GET /v1/login
func Login(c *gin.Context) {
	api := apiutil.New(c)
	appid := c.DefaultQuery("appid", "")
	redirectUri := c.DefaultQuery("redirect_uri", "")
	if appid == "" || redirectUri == "" {
		api.Fail("Invalid params")
		return
	}

	c.Redirect(302, fmt.Sprintf("%s/v1/%s?redirect_uri=%s", config.Server.FrontUrl, appid, url.QueryEscape(redirectUri)))
}
