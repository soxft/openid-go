package model

import "time"

type PassKey struct {
	ID           int        `gorm:"autoIncrement;primaryKey"`
	UserID       int        `gorm:"index;not null"`
	CredentialID string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	PublicKey    string     `gorm:"type:text;not null"`
	Attestation  string     `gorm:"type:varchar(32)"`
	AAGUID       string     `gorm:"type:varchar(64)"`
	SignCount    uint32     `gorm:"type:int unsigned"`
	Transport    string     `gorm:"type:varchar(255)"`
	CloneWarning bool       `gorm:"type:tinyint(1);default:0"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	LastUsedAt   *time.Time `gorm:"type:datetime"`
}

func (PassKey) TableName() string {
	return "pass_keys"
}
