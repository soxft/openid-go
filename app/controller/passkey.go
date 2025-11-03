package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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

	// 先获取 remark 参数（可选）
	var remark string
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// 从表单获取 remark
		remark = c.PostForm("remark")
	} else {
		// 从 JSON 获取 remark
		var jsonBody map[string]interface{}
		bodyBytes, _ := c.GetRawData()
		if len(bodyBytes) > 0 {
			json.Unmarshal(bodyBytes, &jsonBody)
			if r, ok := jsonBody["remark"].(string); ok {
				remark = r
			}
		}
		// 重置请求体供后续解析使用
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	var parsed *protocol.ParsedCredentialCreationData
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// 解析表单数据
		err := c.Request.ParseForm()
		if err != nil {
			api.Fail("invalid form data")
			return
		}
		form := c.Request.Form

		// 打印所有表单字段用于调试
		log.Printf("[DEBUG] All form fields:")
		for key, values := range form {
			log.Printf("  %s = %s", key, strings.Join(values, ", "))
		}
		log.Printf("[DEBUG] Form data: id=%s, rawId=%s, type=%s",
			form.Get("id"), form.Get("rawId"), form.Get("type"))

		// 获取表单字段
		attestationObjectB64 := form.Get("response[attestationObject]")
		clientDataJSONB64 := form.Get("response[clientDataJSON]")
		if attestationObjectB64 == "" || clientDataJSONB64 == "" {
			api.Fail("missing attestation data")
			return
		}

		// 构造符合 webauthn 期望的 JSON 结构
		credentialData := map[string]interface{}{
			"id":    form.Get("id"),
			"rawId": form.Get("rawId"),
			"type":  form.Get("type"),
			"response": map[string]interface{}{
				"attestationObject": attestationObjectB64,
				"clientDataJSON":    clientDataJSONB64,
			},
		}

		// 如果有 authenticatorAttachment 字段，也添加进去
		if attachment := form.Get("authenticatorAttachment"); attachment != "" {
			credentialData["authenticatorAttachment"] = attachment
		}

		// 如果有 transports 字段
		if transports := form.Get("response[transports]"); transports != "" {
			response := credentialData["response"].(map[string]interface{})
			response["transports"] = []string{transports}
		}

		// 序列化为JSON
		jsonData, err := json.Marshal(credentialData)
		if err != nil {
			api.Fail("marshal error")
			return
		}

		log.Printf("[DEBUG] JSON data to parse: %s", string(jsonData))

		// 使用 ParseCredentialCreationResponseBody 解析
		parsed, err = protocol.ParseCredentialCreationResponseBody(strings.NewReader(string(jsonData)))
		if err != nil {
			log.Printf("[ERROR] ParseCredentialCreationResponseBody failed: %v", err)
		}
		if parsed == nil {
			log.Printf("[ERROR] parsed is nil after ParseCredentialCreationResponseBody")
		}
	} else {
		// 默认JSON处理
		parsed, err = protocol.ParseCredentialCreationResponse(c.Request)
	}

	if err != nil {
		api.Fail("invalid payload")
		return
	}

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

	var parsed *protocol.ParsedCredentialAssertionData
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// 处理表单格式的数据
		err := c.Request.ParseForm()
		if err != nil {
			api.Fail("invalid form data")
			return
		}
		form := c.Request.Form

		// 打印所有表单字段用于调试
		log.Printf("[DEBUG] Login form fields:")
		for key, values := range form {
			log.Printf("  %s = %s", key, strings.Join(values, ", "))
		}

		// 获取表单字段
		id := form.Get("id")
		rawId := form.Get("rawId")
		responseType := form.Get("type")
		clientDataJSON := form.Get("response[clientDataJSON]")
		authenticatorData := form.Get("response[authenticatorData]")
		signature := form.Get("response[signature]")
		userHandle := form.Get("response[userHandle]")

		if id == "" || clientDataJSON == "" || authenticatorData == "" || signature == "" {
			log.Printf("[ERROR] Missing required fields - id:%v, clientDataJSON:%v, authenticatorData:%v, signature:%v",
				id != "", clientDataJSON != "", authenticatorData != "", signature != "")
			api.Fail("missing credential data")
			return
		}

		// 构造符合 webauthn 期望的 JSON 结构
		credentialData := map[string]interface{}{
			"id":    id,
			"rawId": rawId,
			"type":  responseType,
			"response": map[string]interface{}{
				"clientDataJSON":    clientDataJSON,
				"authenticatorData": authenticatorData,
				"signature":         signature,
				"userHandle":        userHandle,
			},
		}

		// 序列化为 JSON
		jsonData, err := json.Marshal(credentialData)
		if err != nil {
			api.Fail("marshal error")
			return
		}

		log.Printf("[DEBUG] Login JSON data: %s", string(jsonData))

		// 解析凭证
		parsed, err = protocol.ParseCredentialRequestResponseBody(strings.NewReader(string(jsonData)))
		if err != nil {
			log.Printf("[ERROR] ParseCredentialRequestResponseBody failed: %v", err)
			api.Fail("invalid credential format")
			return
		}
	} else {
		// 解析 JSON 格式的 WebAuthn 响应
		var err error
		parsed, err = protocol.ParseCredentialRequestResponse(c.Request)
		if err != nil {
			log.Printf("[ERROR] parse credential request failed: %v", err)
			api.Fail("invalid credential")
			return
		}
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
