package model

type UniqueId struct {
	ID        int    `gorm:"autoIncrement;primaryKey"`
	UserId    int    `gorm:"index"`
	DevUserId int    `gorm:"index"`
	UniqueId  string `gorm:"type:varchar(128);uniqueIndex"`
	CreateAt  int64  `gorm:"autoCreateTime"`
}

func (UniqueId) TableName() string {
	return "unique_id"
}
