package passkey

import (
	"encoding/base64"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/toolutil"
	"net/url"
)

// PreparePasskey
// @description 准备创建 passkey, 返回服务端支持的参数
func PreparePasskey(user User) (KeyCreatePrepare, error) {
	fullRand := "XOpenID_" + toolutil.RandStr(6) + "_" + toolutil.RandStrInt(8)

	// RP
	u, er := url.Parse(config.Server.FrontUrl)
	if er != nil {
		return KeyCreatePrepare{}, er
	}

	rp := Rp{
		Name: config.Server.Name,
		Id:   u.Host,
	}

	authenticatorSelection := AuthenticatorSelection{
		AuthenticatorAttachment: "platform",
		RequireResidentKey:      true,
		ResidentKey:             "required",
	}

	return KeyCreatePrepare{
		Challenge: base64.URLEncoding.EncodeToString([]byte(fullRand)),
		Rp:        rp,
		User:      user,
		PubKeyCredParams: []PubKeyCredParams{{
			Alg:  -7,
			Type: "public-key",
		}},
		Timeout:                60000,
		Attestation:            "none",
		AuthenticatorSelection: authenticatorSelection,
		Extensions: Extensions{
			CredProps: true,
		},
	}, nil
}

func (k KeyCreatePrepare) GetChallenge() string {
	return k.Challenge
}
