package dto

// ForgetPasswordCodeRequest 忘记密码发送验证码请求
type ForgetPasswordCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgetPasswordUpdateRequest 忘记密码更新请求
type ForgetPasswordUpdateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}