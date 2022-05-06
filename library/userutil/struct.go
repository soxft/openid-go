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

type UserRedis struct {
	UserId   int64
	Username string
	Email    string
	LastTime int64
	LastIp   string
}

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JwtPayload struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Iss      string `json:"iss"`
	Iat      int64  `json:"iat"`
	Jti      string `json:"jti"`
}
