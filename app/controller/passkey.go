package controller

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/passkey"
	"log"
	"strconv"
)

// PasskeyGetOption 获取创建 passkey 选项
//
//	GET /passkey/create/option
func PasskeyGetOption(c *gin.Context) {
	api := apiutil.New(c)

	userID := c.GetInt("userId")
	user := passkey.User{
		Id:          base64.URLEncoding.EncodeToString([]byte(strconv.Itoa(userID))),
		Name:        c.GetString("username"),
		DisplayName: c.GetString("username"),
	}

	response, err := passkey.PreparePasskey(user)
	if err != nil {
		log.Printf("[ERROR] passkey prepare failed: %v", err)

		api.FailWithData("get option failed", gin.H{
			"error": "failed",
		})
		return
	}

	api.SuccessWithData("success", response)
}
