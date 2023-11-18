package codeutil

import (
	"context"
	"time"
)

type Coder interface {
	Create(length int) string
	Save(topic string, email string, code string, timeout time.Duration) error
	Check(topic string, email string, code string) (bool, error)
	Consume(topic string, email string)
}

type VerifyCode struct {
	ctx context.Context
}
