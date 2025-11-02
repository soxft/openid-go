package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"gorm.io/gorm"

	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/library/apiutil"
	"github.com/soxft/openid-go/library/passkey"
	"github.com/soxft/openid-go/library/userutil"
	"github.com/soxft/openid-go/process/dbutil"
)

// PasskeyRegistrationOptions 获取 Passkey 注册选项
//
//	GET /passkey/register/options
func PasskeyRegistrationOptions(c *gin.Context) {
	api := apiutil.New(c)

	account, err := getAccount(c)
	if err != nil {
		api.Fail("user not found")
		return
	}

	options, err := passkey.BeginRegistration(c.Request.Context(), *account)
	if err != nil {
		log.Printf("[ERROR] passkey begin registration failed: %v", err)
		api.Fail("生成 Passkey 参数失败")
		return
	}

	api.SuccessWithData("success", options)
}

// PasskeyRegistrationFinish 完成 Passkey 注册
//
//	POST /passkey/register
func PasskeyRegistrationFinish(c *gin.Context) {
	api := apiutil.New(c)
	account, err := getAccount(c)
	if err != nil {
		api.Fail("user not found")
		return
	}

	var parsed *protocol.ParsedCredentialCreationData
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// 解析表单数据
		err := c.Request.ParseForm()
		if err != nil {
			api.Fail("invalid form data")
			return
		}
		form := c.Request.Form

		// 解码base64数据
		attestationObjectB64 := form.Get("response[attestationObject]")
		clientDataJSONB64 := form.Get("response[clientDataJSON]")
		if attestationObjectB64 == "" || clientDataJSONB64 == "" {
			api.Fail("missing attestation data")
			return
		}

		attestationObject, err := base64.StdEncoding.DecodeString(attestationObjectB64)
		if err != nil {
			api.Fail("invalid attestationObject")
			return
		}
		clientDataJSON, err := base64.StdEncoding.DecodeString(clientDataJSONB64)
		if err != nil {
			api.Fail("invalid clientDataJSON")
			return
		}

		// 构造CredentialCreationResponse
		ccr := protocol.CredentialCreationResponse{
			PublicKeyCredential: protocol.PublicKeyCredential{
				Credential: protocol.Credential{
					ID:   form.Get("id"),
					Type: form.Get("type"),
				},
				RawID: protocol.URLEncodedBase64(form.Get("rawId")),
			},
			AttestationResponse: protocol.AuthenticatorAttestationResponse{
				AttestationObject: attestationObject,
				ClientDataJSON:    clientDataJSON,
			},
		}

		// 序列化为JSON并解析
		jsonData, err := json.Marshal(ccr)
		if err != nil {
			api.Fail("marshal error")
			return
		}
		parsed, err = protocol.ParseCredentialCreationResponseBody(strings.NewReader(string(jsonData)))
	} else {
		// 默认JSON处理
		parsed, err = protocol.ParseCredentialCreationResponse(c.Request)
	}

	if err != nil {
		api.Fail("invalid payload")
		return
	}

	credential, err := passkey.CompleteRegistration(c.Request.Context(), *account, parsed)
	if err != nil {
		if errors.Is(err, passkey.ErrSessionNotFound) {
			api.Fail("挑战已过期，请重试")
			return
		}
		log.Printf("[ERROR] passkey finish registration failed: %v", err)
		api.Fail("注册失败")
		return
	}

	api.SuccessWithData("success", gin.H{
		"passkeyId": credential.ID,
	})
}

// PasskeyLoginOptions 获取 Passkey 登录选项
//
//	GET /passkey/login/options
func PasskeyLoginOptions(c *gin.Context) {
	api := apiutil.New(c)
	account, err := getAccount(c)
	if err != nil {
		api.Fail("user not found")
		return
	}

	options, err := passkey.BeginLogin(c.Request.Context(), *account)
	if err != nil {
		if errors.Is(err, passkey.ErrNoPasskeyRegistered) {
			api.Fail("未绑定 Passkey")
			return
		}
		log.Printf("[ERROR] passkey begin login failed: %v", err)
		api.Fail("生成登录参数失败")
		return
	}

	api.SuccessWithData("success", options)
}

// PasskeyLoginFinish 完成 Passkey 登录
//
//	POST /passkey/login
func PasskeyLoginFinish(c *gin.Context) {
	api := apiutil.New(c)
	account, err := getAccount(c)
	if err != nil {
		api.Fail("user not found")
		return
	}

	passkeyCredential, err := passkey.CompleteLogin(c.Request.Context(), *account, c.Request)
	if err != nil {
		if errors.Is(err, passkey.ErrSessionNotFound) {
			api.Fail("挑战已过期，请重试")
			return
		}
		if errors.Is(err, passkey.ErrNoPasskeyRegistered) {
			api.Fail("未绑定 Passkey")
			return
		}
		log.Printf("[ERROR] passkey finish login failed: %v", err)
		api.Fail("登录失败")
		return
	}

	token, err := generateLoginToken(c, account.ID)
	if err != nil {
		api.Fail("system error")
		return
	}

	api.SuccessWithData("success", gin.H{
		"token":     token,
		"passkeyId": passkeyCredential.ID,
	})
}

// PasskeyList 列出用户绑定的所有 Passkey
//
//	GET /passkey
func PasskeyList(c *gin.Context) {
	api := apiutil.New(c)
	account, err := getAccount(c)
	if err != nil {
		api.Fail("user not found")
		return
	}

	passkeys, err := passkey.ListUserPasskeys(account.ID)
	if err != nil {
		log.Printf("[ERROR] passkey list failed: %v", err)
		api.Fail("获取失败")
		return
	}

	summaries := make([]passkey.Summary, 0, len(passkeys))
	for _, item := range passkeys {
		summaries = append(summaries, passkey.Summary{
			ID:           item.ID,
			CreatedAt:    item.CreatedAt,
			LastUsedAt:   item.LastUsedAt,
			CloneWarning: item.CloneWarning,
			SignCount:    item.SignCount,
			Transports:   passkey.SplitTransports(item.Transport),
		})
	}

	api.SuccessWithData("success", gin.H{
		"items": summaries,
	})
}

// PasskeyDelete 删除指定 Passkey
//
//	DELETE /passkey/:id
func PasskeyDelete(c *gin.Context) {
	api := apiutil.New(c)
	account, err := getAccount(c)
	if err != nil {
		api.Fail("user not found")
		return
	}

	passkeyID, err := strconv.Atoi(c.Param("id"))
	if err != nil || passkeyID <= 0 {
		api.Fail("invalid id")
		return
	}

	if err := passkey.DeleteUserPasskey(account.ID, passkeyID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.Fail("Passkey 不存在")
			return
		}
		log.Printf("[ERROR] passkey delete failed: %v", err)
		api.Fail("删除失败")
		return
	}

	api.Success("success")
}

func getAccount(c *gin.Context) (*model.Account, error) {
	userID := c.GetInt("userId")
	if userID == 0 {
		return nil, errors.New("empty user id")
	}

	var account model.Account
	if err := dbutil.D.Where("id = ?", userID).Take(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func generateLoginToken(c *gin.Context, userID int) (string, error) {
	return userutil.GenerateJwt(userID, c.ClientIP())
}
