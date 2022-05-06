package version_one

import "github.com/gin-gonic/gin"

// Login
// @description v1 登录
func Login(c *gin.Context) {
	c.Redirect(302, "/login?appid="+c.Query("appid")+"&redirect_uri="+c.Query("redirect_uri"))
}

// LoginHandler
// @description 处理登录xhr请求
func LoginHandler(c *gin.Context) {

}
