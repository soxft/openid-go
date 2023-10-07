package model

type Account struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"uniqueIndex"`
	Password string
	Email    string `gorm:"uniqueIndex"`
	RegTime  int64
	RegIp    string
	LastTime int64
	LastIp   string
}
