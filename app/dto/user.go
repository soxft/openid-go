package dto

// UserPasswordUpdateRequest 用户更新密码请求
type UserPasswordUpdateRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=64"`
}

// UserEmailUpdateCodeRequest 发送邮箱更新验证码请求
type UserEmailUpdateCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UserEmailUpdateRequest 更新邮箱请求
type UserEmailUpdateRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RegTime  int64  `json:"reg_time"`
	LastTime int64  `json:"last_time"`
}

// UserStatusResponse 用户状态响应
type UserStatusResponse struct {
	Login    bool   `json:"login"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}