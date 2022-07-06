package version_one

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"openid/library/apiutil"
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
	c.Redirect(302, "/v1/"+appid+"?redirect_uri="+url.QueryEscape(redirectUri))
}
