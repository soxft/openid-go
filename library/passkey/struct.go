package passkey

import (
	"github.com/go-webauthn/webauthn/protocol"
)

// RegistrationOptions 兼容旧代码的类型别名
type RegistrationOptions = protocol.PublicKeyCredentialCreationOptions

// LoginOptions 兼容旧代码的类型别名
type LoginOptions = protocol.PublicKeyCredentialRequestOptions

// Summary 用于对外输出的 Passkey 信息
type Summary struct {
	ID           int      `json:"id"`
	CreatedAt    int64    `json:"createdAt"`
	LastUsedAt   int64    `json:"lastUsedAt"`
	CloneWarning bool     `json:"cloneWarning"`
	SignCount    uint32   `json:"signCount"`
	Transports   []string `json:"transports"`
}
