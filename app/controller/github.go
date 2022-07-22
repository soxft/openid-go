package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid/config"
	"github.com/soxft/openid/library/apiutil"
	"github.com/soxft/openid/thirdpart"
	"github.com/soxft/openid/thirdpart/github"
)

func GithubRedirect(c *gin.Context) {
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=read:user&redirect_uri=%s", config.Github.ClientID, config.Server.Url+"/login/github/handler")
	c.Redirect(302, url)
}

func GithubHandler(c *gin.Context) {
	// get access token
	api := apiutil.New(c)
	userId, err := github.New().Handler(c)
	if err != nil {
		api.Fail("登录失败，稍后再试")
		return
	}
	jwt, err := thirdpart.Handler(userId, github.Platform, c.ClientIP())
	if err != nil {
		api.Fail(err.Error())
		return
	}

	api.SuccessWithData("登录成功", gin.H{
		"token": jwt,
	})
}
