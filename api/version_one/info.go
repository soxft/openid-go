package version_one

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/api/version_one/helper"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/apputil"
)

type InfoRequest struct {
	Token     string `json:"token" binding:"required"`
	AppId     string `json:"appid" binding:"required"`
	AppSecret string `json:"app_secret" binding:"required"`
}

type InfoResponse struct {
	OpenId   string `json:"openId"`
	UniqueId string `json:"uniqueId"`
}

// Info
// @description 获取openid 和 uniqueId
// @route POST /v1/info
func Info(c *gin.Context) {
	api := apiutil.New(c)

	var req InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.Fail("Invalid params")
		return
	}

	// 判断appId与appSecret是否正确
	if err := apputil.CheckAppSecret(req.AppId, req.AppSecret); err != nil {
		api.Fail(err.Error())
		return
	}

	// 检测token是否正确 并获取userId
	userId, err := helper.GetUserIdByToken(c, req.AppId, req.AppSecret)
	if err != nil {
		if errors.Is(err, helper.ErrTokenNotExists) {
			api.Fail("Token not exists")
			return
		}
		api.Fail(err.Error())
		return
	}
	userIds, err := helper.GetUserIds(req.AppId, userId)
	if err != nil {
		api.Fail(err.Error())
		return
	}
	// delete token
	_ = helper.DeleteToken(c, req.AppId, req.Token)
	api.SuccessWithData("success", InfoResponse{
		OpenId:   userIds.OpenId,
		UniqueId: userIds.UniqueId,
	})
}
