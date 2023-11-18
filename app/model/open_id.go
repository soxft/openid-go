package model

type OpenId struct {
	ID       int    `gorm:"autoIncrement;primaryKey"`
	UserId   int    `gorm:"index"`
	AppId    string `gorm:"type:varchar(20);index"`
	OpenId   string `gorm:"type:varchar(128);uniqueIndex"`
	CreateAt int64  `gorm:"autoCreateTime"`
}

func (OpenId) TableName() string {
	return "open_id"
}
