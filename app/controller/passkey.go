package controller

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"gorm.io/gorm"

	"github.com/soxft/openid-go/app/dto"
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

	// 解析 JSON 格式的 WebAuthn 响应
	parsed, err := protocol.ParseCredentialCreationResponse(c.Request)
	if err != nil {
		log.Printf("[ERROR] parse credential creation failed: %v", err)
		api.Fail("invalid payload")
		return
	}

	// 获取备注（可选）
	var reqBody dto.PasskeyRegistrationFinishRequest
	_ = c.ShouldBindJSON(&reqBody)
	remark := reqBody.Remark

	// 完成注册，并传递备注
	credential, err := passkey.CompleteRegistrationWithRemark(c.Request.Context(), *account, parsed, remark)
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

// PasskeyLoginOptions 获取 Passkey 登录选项（无需用户名）
//
//	GET /passkey/login/options
func PasskeyLoginOptions(c *gin.Context) {
	api := apiutil.New(c)

	// 模式2：无用户名登录（无条件 UI）
	// 生成通用的登录挑战，不指定 allowCredentials
	options, err := passkey.BeginDiscoverableLogin(c.Request.Context())
	if err != nil {
		log.Printf("[ERROR] passkey begin discoverable login failed: %v", err)
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

	// 解析 JSON 格式的 WebAuthn 响应
	parsed, err := protocol.ParseCredentialRequestResponse(c.Request)
	if err != nil {
		log.Printf("[ERROR] parse credential request failed: %v", err)
		api.Fail("invalid credential")
		return
	}

	if parsed == nil {
		log.Printf("[ERROR] parsed credential is nil")
		api.Fail("invalid credential")
		return
	}

	// 方式1：通过 userHandle（如果客户端提供了）
	var account model.Account
	var found bool

	if len(parsed.Response.UserHandle) > 0 {
		// userHandle 是用户 ID 的 base64 编码
		userIDStr := string(parsed.Response.UserHandle)
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			if err := dbutil.D.Where("id = ?", userID).Take(&account).Error; err == nil {
				found = true
				log.Printf("[DEBUG] Found user by userHandle: %d", userID)
			}
		}
	}

	// 方式2：通过 credential ID 查找
	if !found {
		credentialID := passkey.EncodeKey(parsed.RawID)
		var passkeyRecord model.PassKey
		if err := dbutil.D.Where("credential_id = ?", credentialID).Take(&passkeyRecord).Error; err == nil {
			if err := dbutil.D.Where("id = ?", passkeyRecord.UserID).Take(&account).Error; err == nil {
				found = true
				log.Printf("[DEBUG] Found user by credential ID: %d", passkeyRecord.UserID)
			}
		}
	}

	if !found {
		log.Printf("[ERROR] User not found for credential ID: %s", base64.StdEncoding.EncodeToString(parsed.RawID))
		api.Fail("invalid credential")
		return
	}

	// 验证登录
	passkeyCredential, err := passkey.CompleteDiscoverableLogin(c.Request.Context(), account, parsed)
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

	// 登录成功，生成 JWT Token
	token, err := generateLoginToken(c, account.ID)
	if err != nil {
		api.Fail("system error")
		return
	}

	api.SuccessWithData("success", gin.H{
		"token":     token,
		"passkeyId": passkeyCredential.ID,
		"username":  account.Username,
		"email":     account.Email,
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
			Remark:       item.Remark, // 包含备注信息
			CreatedAt:    item.CreatedAt,
			LastUsedAt:   item.LastUsedAt,
			CloneWarning: item.CloneWarning,
			SignCount:    item.SignCount,
			Transports:   passkey.SplitTransports(item.Transport),
		})
	}

	api.SuccessWithData("success", summaries)
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
