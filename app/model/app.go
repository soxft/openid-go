package model

type App struct {
	ID         int    `gorm:"type:bigint(20);primaryKey;autoIncrement"`
	UserId     int    `gorm:"type:bigint(20);index"`
	AppId      string `gorm:"type:varchar(20);uniqueIndex"`
	AppName    string `gorm:"type:varchar(128)"`
	AppSecret  string `gorm:"type:varchar(100);uniqueIndex"`
	AppGateway string `gorm:"type:varchar(200)"`
	CreateAt   int64  `gorm:"autoCreateTime"`
}
