package model

type UniqueId struct {
	ID        int    `gorm:"type:bigint(20);primaryKey;autoIncrement"`
	UserId    int    `gorm:"type:bigint(20);index"`
	DevUserId int    `gorm:"type:bigint(20);index"`
	UniqueId  string `gorm:"type:varchar(128);uniqueIndex"`
	CreateAt  int64  `gorm:"autoCreateTime"`
}

func (UniqueId) TableName() string {
	return "unique_id"
}
