package passkey

import (
	"context"
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
			RPOrigin:      origin,
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

// CompleteRegistration 校验并保存 passkey
func CompleteRegistration(ctx context.Context, account model.Account, parsed *protocol.ParsedCredentialCreationData) (*model.PassKey, error) {
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

	passkey, err := saveCredential(account.ID, credential)
	if err != nil {
		return nil, err
	}

	_ = deleteRegistrationSession(ctx, account.ID)
	return passkey, nil
}

// BeginLogin 创建登录挑战
func BeginLogin(ctx context.Context, account model.Account) (LoginOptions, error) {
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
