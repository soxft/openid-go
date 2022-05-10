package version_one

import "github.com/gin-gonic/gin"

// Login
// @description v1 登录
// @route GET /v1/login
func Login(c *gin.Context) {
	c.Redirect(302, "/login?"+c.Request.URL.Query().Encode())
}
