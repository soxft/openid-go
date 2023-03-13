package model

type Account struct {
	ID       int
	Username string
	Password string
	Email    string
	RegTime  int64
	RegIp    string
	LastTime int64
	LastIp   string
}
