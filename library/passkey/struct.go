package passkey

import (
	"time"

	"github.com/go-webauthn/webauthn/protocol"
)

// RegistrationOptions 兼容旧代码的类型别名
type RegistrationOptions = protocol.PublicKeyCredentialCreationOptions

// LoginOptions 兼容旧代码的类型别名
type LoginOptions = protocol.PublicKeyCredentialRequestOptions

// Summary 用于对外输出的 Passkey 信息
type Summary struct {
	ID           int        `json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty"`
	CloneWarning bool       `json:"cloneWarning"`
	SignCount    uint32     `json:"signCount"`
	Transports   []string   `json:"transports"`
}
