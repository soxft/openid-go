package dto

// AppCreateRequest 创建应用请求
type AppCreateRequest struct {
	AppName string `json:"app_name" binding:"required"`
}

// AppEditRequest 编辑应用请求
type AppEditRequest struct {
	AppName    string `json:"app_name" binding:"required"`
	AppGateway string `json:"app_gateway" binding:"required,max=200"`
}

// AppListRequest 获取应用列表请求
type AppListRequest struct {
	Page    int `json:"page,omitempty" form:"page" binding:"omitempty,min=1"`
	PerPage int `json:"per_page,omitempty" form:"per_page" binding:"omitempty,min=1,max=100"`
}

// AppListResponse 应用列表响应
type AppListResponse struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

// AppInfoResponse 应用信息响应
type AppInfoResponse struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppGateway string `json:"app_gateway"`
	AppSecret  string `json:"app_secret,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	UpdateTime int64  `json:"update_time,omitempty"`
}

// AppSecretResponse 重置密钥响应
type AppSecretResponse struct {
	Secret string `json:"secret"`
}