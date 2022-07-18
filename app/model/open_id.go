package model

type OpenId struct {
	ID       int
	UserId   int
	AppId    string
	OpenId   string
	CreateAt int64 `gorm:"autoCreateTime"`
}

func (OpenId) TableName() string {
	return "open_id"
}
