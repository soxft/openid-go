package thirdpartutil

import "github.com/gin-gonic/gin"

type thirdPart interface {
	GetUserInfo(token string) (interface{}, error)
	Redirect(ctx gin.Context)
}

type ThirdPartCtx struct {
	RedirectUrl string
	SecretId    string
	SecretKey   string
}

func NewThirdPartCtx(redirectUrl, secretId, secretKey string) *ThirdPartCtx {
	return &ThirdPartCtx{
		RedirectUrl: redirectUrl,
		SecretId:    secretId,
		SecretKey:   secretKey,
	}
}
