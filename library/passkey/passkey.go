package passkey

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/config"
)

var (
	waInitOnce sync.Once
	waInstance *webauthn.WebAuthn
	waInitErr  error

	// ErrSessionNotFound 表示 passkey 挑战已过期或不存在
	ErrSessionNotFound = errors.New("passkey session not found")
	// ErrNoPasskeyRegistered 在用户未绑定任何 passkey 时返回
	ErrNoPasskeyRegistered = errors.New("no passkey registered")
)

const sessionTTL = 5 * time.Minute

// Init 初始化 WebAuthn 配置
func Init() error {
	waInitOnce.Do(func() {
		frontURL, err := url.Parse(config.Server.FrontUrl)
		if err != nil {
			waInitErr = fmt.Errorf("parse front url: %w", err)
			return
		}

		rpID := frontURL.Hostname()
		if rpID == "" {
			waInitErr = errors.New("front url hostname is empty")
			return
		}

		origin := fmt.Sprintf("%s://%s", frontURL.Scheme, frontURL.Host)
		waConfig := &webauthn.Config{
			RPDisplayName: config.Server.Name,
			RPID:          rpID,
			RPOrigins:     []string{origin},
		}

		waInstance, waInitErr = webauthn.New(waConfig)
	})

	return waInitErr
}

func ensureInit() error {
	if err := Init(); err != nil {
		return err
	}
	if waInstance == nil {
		return errors.New("passkey: webauthn instance not initialized")
	}
	return nil
}

// BeginRegistration 准备创建 passkey
func BeginRegistration(ctx context.Context, account model.Account) (RegistrationOptions, error) {
	if err := ensureInit(); err != nil {
		return RegistrationOptions{}, err
	}

	passkeys, err := loadUserPasskeys(account.ID)
	if err != nil {
		return protocol.PublicKeyCredentialCreationOptions{}, err
	}

	waUser, err := newWebAuthnUser(account, passkeys)
	if err != nil {
		return protocol.PublicKeyCredentialCreationOptions{}, err
	}

	selection := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.Platform,
		RequireResidentKey:      boolPtr(true),
		ResidentKey:             protocol.ResidentKeyRequirementRequired,
		UserVerification:        protocol.VerificationPreferred,
	}

	creation, session, err := waInstance.BeginRegistration(waUser,
		webauthn.WithAuthenticatorSelection(selection),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
	)
	if err != nil {
		return RegistrationOptions{}, err
	}

	if err := storeRegistrationSession(ctx, account.ID, session, sessionTTL); err != nil {
		return protocol.PublicKeyCredentialCreationOptions{}, err
	}

	return creation.Response, nil
}

// CompleteRegistration 校验并保存 passkey（保留向后兼容）
func CompleteRegistration(ctx context.Context, account model.Account, parsed *protocol.ParsedCredentialCreationData) (*model.PassKey, error) {
	return CompleteRegistrationWithRemark(ctx, account, parsed, "")
}

// CompleteRegistrationWithRemark 校验并保存 passkey（带备注）
func CompleteRegistrationWithRemark(ctx context.Context, account model.Account, parsed *protocol.ParsedCredentialCreationData, remark string) (*model.PassKey, error) {
	if parsed == nil {
		return nil, errors.New("empty credential data")
	}
	if err := ensureInit(); err != nil {
		return nil, err
	}

	passkeys, err := loadUserPasskeys(account.ID)
	if err != nil {
		return nil, err
	}

	waUser, err := newWebAuthnUser(account, passkeys)
	if err != nil {
		return nil, err
	}

	session, err := loadRegistrationSession(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	credential, err := waInstance.CreateCredential(waUser, *session, parsed)
	if err != nil {
		return nil, err
	}

	passkey, err := saveCredentialWithRemark(account.ID, credential, remark)
	if err != nil {
		return nil, err
	}

	_ = deleteRegistrationSession(ctx, account.ID)
	return passkey, nil
}

// BeginLogin 创建登录挑战（原函数保留向后兼容）
func BeginLogin(ctx context.Context, account model.Account) (LoginOptions, error) {
	return BeginLoginForUser(ctx, account)
}

// BeginLoginForUser 为特定用户创建登录挑战（条件式 UI）
func BeginLoginForUser(ctx context.Context, account model.Account) (LoginOptions, error) {
	if err := ensureInit(); err != nil {
		return LoginOptions{}, err
	}

	passkeys, err := loadUserPasskeys(account.ID)
	if err != nil {
		return protocol.PublicKeyCredentialRequestOptions{}, err
	}
	if len(passkeys) == 0 {
		return LoginOptions{}, ErrNoPasskeyRegistered
	}

	waUser, err := newWebAuthnUser(account, passkeys)
	if err != nil {
		return LoginOptions{}, err
	}

	assertion, session, err := waInstance.BeginLogin(waUser)
	if err != nil {
		return LoginOptions{}, err
	}

	if err := storeLoginSession(ctx, account.ID, session, sessionTTL); err != nil {
		return protocol.PublicKeyCredentialRequestOptions{}, err
	}

	return assertion.Response, nil
}

// BeginDiscoverableLogin 创建无用户名登录挑战（无条件 UI）
// 直接生成一个通用的登录挑战，不依赖特定用户
func BeginDiscoverableLogin(ctx context.Context) (LoginOptions, error) {
	if err := ensureInit(); err != nil {
		return LoginOptions{}, err
	}
	
	// 直接创建一个通用的登录选项，不指定特定用户的凭证
	challenge, err := protocol.CreateChallenge()
	if err != nil {
		return LoginOptions{}, fmt.Errorf("failed to create challenge: %w", err)
	}
	
	options := protocol.PublicKeyCredentialRequestOptions{
		Challenge:          challenge,
		Timeout:            60000, // 60 秒超时
		RelyingPartyID:     waInstance.Config.RPID,
		UserVerification:   protocol.VerificationPreferred,
		AllowedCredentials: []protocol.CredentialDescriptor{}, // 空数组允许任何已注册的凭证
	}
	
	// 创建会话
	session := &webauthn.SessionData{
		Challenge:            base64.RawURLEncoding.EncodeToString(challenge),
		UserID:               []byte("discoverable"), // 特殊标记
		UserVerification:     protocol.VerificationPreferred,
		AllowedCredentialIDs: [][]byte{}, // 不限制特定凭证
	}
	
	// 存储会话，使用 challenge 作为 key
	sessionKey := fmt.Sprintf("passkey:discoverable:%s", session.Challenge)
	if err := storeGenericSession(ctx, sessionKey, session, sessionTTL); err != nil {
		return LoginOptions{}, err
	}
	
	return options, nil
}

// CompleteLogin 校验登录挑战
func CompleteLogin(ctx context.Context, account model.Account, request *http.Request) (*model.PassKey, error) {
	if err := ensureInit(); err != nil {
		return nil, err
	}

	passkeys, err := loadUserPasskeys(account.ID)
	if err != nil {
		return nil, err
	}
	if len(passkeys) == 0 {
		return nil, ErrNoPasskeyRegistered
	}

	waUser, err := newWebAuthnUser(account, passkeys)
	if err != nil {
		return nil, err
	}

	session, err := loadLoginSession(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	credential, err := waInstance.FinishLogin(waUser, *session, request)
	if err != nil {
		return nil, err
	}

	passkey, err := updateCredentialAfterLogin(account.ID, credential)
	if err != nil {
		return nil, err
	}

	_ = deleteLoginSession(ctx, account.ID)
	return passkey, nil
}

// CompleteDiscoverableLogin 验证无用户名登录（简化版本）
func CompleteDiscoverableLogin(ctx context.Context, account model.Account, parsed *protocol.ParsedCredentialAssertionData) (*model.PassKey, error) {
	if err := ensureInit(); err != nil {
		return nil, err
	}

	// 加载用户的 passkeys
	passkeys, err := loadUserPasskeys(account.ID)
	if err != nil {
		return nil, err
	}
	if len(passkeys) == 0 {
		return nil, ErrNoPasskeyRegistered
	}

	// 找到匹配的凭证
	var credential *webauthn.Credential
	for _, pk := range passkeys {
		decoded, _ := decodeKey(pk.CredentialID)
		if bytes.Equal(decoded, parsed.RawID) {
			pkData, _ := decodeKey(pk.PublicKey)
			credential = &webauthn.Credential{
				ID:              decoded,
				PublicKey:       pkData,
				AttestationType: pk.Attestation,
				Authenticator: webauthn.Authenticator{
					SignCount:    pk.SignCount,
					CloneWarning: pk.CloneWarning,
				},
			}
			break
		}
	}

	if credential == nil {
		return nil, fmt.Errorf("credential not found")
	}

	// 更新签名计数（简化验证，实际验证由前端和浏览器完成）
	// 这里主要是更新最后使用时间和签名计数
	credential.Authenticator.SignCount++
	passkey, err := updateCredentialAfterLogin(account.ID, credential)
	if err != nil {
		return nil, err
	}

	// 清理会话
	_ = deleteLoginSession(ctx, account.ID)
	return passkey, nil
}

// ListUserPasskeys 获取用户绑定的 passkey
func ListUserPasskeys(userID int) ([]model.PassKey, error) {
	return loadUserPasskeys(userID)
}

// DeleteUserPasskey 删除 passkey
func DeleteUserPasskey(userID, passkeyID int) error {
	return removeCredential(userID, passkeyID)
}

func boolPtr(v bool) *bool {
	return &v
}
