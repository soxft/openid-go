package model

type PassKey struct {
	ID     int    `gorm:"autoIncrement;primaryKey"`
	UserID int    `gorm:"index"`
	Key    string `gorm:"type:varchar(128);uniqueIndex"`
}
