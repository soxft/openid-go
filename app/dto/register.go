package dto

// RegisterCodeRequest 发送注册验证码请求
type RegisterCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RegisterSubmitRequest 提交注册请求
type RegisterSubmitRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Code       string `json:"code" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}