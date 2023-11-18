package model

type Account struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(20);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(128);not null"`
	Email    string `gorm:"type:varchar(128);uniqueIndex"`
	RegTime  int64  `gorm:"bigint(20)"`
	RegIp    string `gorm:"type:varchar(128)"`
	LastTime int64  `gorm:"bigint(20)"`
	LastIp   string `gorm:"type:varchar(128)"`
}
