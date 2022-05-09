package version_one

import (
	"github.com/gin-gonic/gin"
)

// Login
// @description v1 登录
// @route GET /v1/login
func Login(c *gin.Context) {
	c.Redirect(302, "/login?"+c.Request.URL.Query().Encode())
}

// Code
// @description 处理登录xhr请求 获取code并跳转到redirect_uri
// @route POST /v1/info
func Code(c *gin.Context) {
	appId := c.Query("appid")
	redirectUri := c.Query("redirect_uri")
	if appId == "" || redirectUri == "" {

	}
}
