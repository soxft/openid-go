package codeutil

type Coder interface {
	Create(length int) string
	Save(topic string, timeout int64, email string, code string) error
	Check(topic string, email string, code string) (bool, error)
	Consume(topic string, email string)
}

type VerifyCode struct {
}
