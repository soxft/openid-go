package model

type App struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	UserId     int    `gorm:"index"`
	AppId      string `gorm:"uniqueIndex"`
	AppName    string
	AppSecret  string `gorm:"uniqueIndex"`
	AppGateway string
	CreateAt   int64 `gorm:"autoCreateTime"`
}
