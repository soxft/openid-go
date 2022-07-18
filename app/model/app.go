package model

type App struct {
	ID         int
	UserId     int
	AppId      string
	AppName    string
	AppSecret  string
	AppGateway string
	CreateAt   int64 `gorm:"autoCreateTime"`
}
