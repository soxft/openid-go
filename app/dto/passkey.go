package dto

// PasskeyRegistrationFinishRequest Passkey注册完成请求
type PasskeyRegistrationFinishRequest struct {
	ID     string `json:"id" binding:"required"`
	RawID  string `json:"rawId" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Remark string `json:"remark,omitempty"`
	Response struct {
		AttestationObject string   `json:"attestationObject" binding:"required"`
		ClientDataJSON    string   `json:"clientDataJSON" binding:"required"`
		Transports        []string `json:"transports,omitempty"`
	} `json:"response" binding:"required"`
	AuthenticatorAttachment string `json:"authenticatorAttachment,omitempty"`
}

// PasskeyLoginFinishRequest Passkey登录完成请求
type PasskeyLoginFinishRequest struct {
	ID    string `json:"id" binding:"required"`
	RawID string `json:"rawId" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Response struct {
		ClientDataJSON    string `json:"clientDataJSON" binding:"required"`
		AuthenticatorData string `json:"authenticatorData" binding:"required"`
		Signature         string `json:"signature" binding:"required"`
		UserHandle        string `json:"userHandle,omitempty"`
	} `json:"response" binding:"required"`
}

// PasskeyLoginResponse Passkey登录响应
type PasskeyLoginResponse struct {
	Token     string `json:"token"`
	PasskeyID string `json:"passkeyId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

// PasskeyRegistrationResponse Passkey注册响应
type PasskeyRegistrationResponse struct {
	PasskeyID string `json:"passkeyId"`
}

// PasskeySummaryResponse Passkey摘要响应
type PasskeySummaryResponse struct {
	ID           int      `json:"id"`
	Remark       string   `json:"remark,omitempty"`
	CreatedAt    int64    `json:"created_at"`
	LastUsedAt   int64    `json:"last_used_at"`
	CloneWarning bool     `json:"clone_warning"`
	SignCount    uint32   `json:"sign_count"`
	Transports   []string `json:"transports,omitempty"`
}