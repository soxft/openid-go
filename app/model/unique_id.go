package model

type UniqueId struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	UserId    int    `gorm:"index"`
	DevUserId int    `gorm:"index"`
	UniqueId  string `gorm:"uniqueIndex"`
	CreateAt  int64  `gorm:"autoCreateTime"`
}

func (UniqueId) TableName() string {
	return "unique_id"
}
