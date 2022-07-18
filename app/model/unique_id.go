package model

type UniqueId struct {
	ID        int
	UserId    int
	DevUserId int
	UniqueId  string
	CreateAt  int64 `gorm:"autoCreateTime"`
}

func (UniqueId) TableName() string {
	return "unique_id"
}
