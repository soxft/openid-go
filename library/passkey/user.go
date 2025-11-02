package passkey

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/soxft/openid-go/app/model"
)

type webAuthnUser struct {
	account     model.Account
	credentials []webauthn.Credential
}

func newWebAuthnUser(account model.Account, passkeys []model.PassKey) (*webAuthnUser, error) {
	credentials, err := convertPasskeys(passkeys)
	if err != nil {
		return nil, err
	}
	return &webAuthnUser{
		account:     account,
		credentials: credentials,
	}, nil
}

func (u *webAuthnUser) WebAuthnID() []byte {
	return []byte(strconv.Itoa(u.account.ID))
}

func (u *webAuthnUser) WebAuthnName() string {
	if u.account.Username != "" {
		return u.account.Username
	}
	return strconv.Itoa(u.account.ID)
}

func (u *webAuthnUser) WebAuthnDisplayName() string {
	if u.account.Email != "" {
		return u.account.Email
	}
	return u.WebAuthnName()
}

func (u *webAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *webAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func convertPasskeys(passkeys []model.PassKey) ([]webauthn.Credential, error) {
	credentials := make([]webauthn.Credential, 0, len(passkeys))
	for _, key := range passkeys {
		credentialID, err := decodeKey(key.CredentialID)
		if err != nil {
			return nil, err
		}
		publicKey, err := decodeKey(key.PublicKey)
		if err != nil {
			return nil, err
		}
		var aaguid []byte
		if key.AAGUID != "" {
			aaguid, err = decodeKey(key.AAGUID)
			if err != nil {
				return nil, err
			}
		}

		transports := make([]protocol.AuthenticatorTransport, 0)
		if key.Transport != "" {
			for _, transport := range strings.Split(key.Transport, ",") {
				transport = strings.TrimSpace(transport)
				if transport == "" {
					continue
				}
				transports = append(transports, protocol.AuthenticatorTransport(transport))
			}
		}

		credential := webauthn.Credential{
			ID:              credentialID,
			PublicKey:       publicKey,
			AttestationType: key.Attestation,
			Transport:       transports,
			Authenticator: webauthn.Authenticator{
				AAGUID:       aaguid,
				SignCount:    key.SignCount,
				CloneWarning: key.CloneWarning,
			},
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}

func encodeKey(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

func decodeKey(encoded string) ([]byte, error) {
	if encoded == "" {
		return nil, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func transportsToString(values []protocol.AuthenticatorTransport) string {
	if len(values) == 0 {
		return ""
	}
	items := make([]string, 0, len(values))
	for _, v := range values {
		items = append(items, string(v))
	}
	return strings.Join(items, ",")
}

func SplitTransports(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}
