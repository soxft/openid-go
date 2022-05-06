package userutil

type User struct {
	Username string
	Password string
	Salt     string
	Email    string
	RegTime  int64
	RegIp    string
	LastTime int64
	LastIp   string
}
