package model

type OpenId struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	UserId   int    `gorm:"index"`
	AppId    string `gorm:"index"`
	OpenId   string `gorm:"uniqueIndex"`
	CreateAt int64  `gorm:"autoCreateTime"`
}

func (OpenId) TableName() string {
	return "open_id"
}
