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

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JwtPayload struct {
	UserId int    `json:"userId"`
	Iss    string `json:"iss"`
	Iat    int64  `json:"iat"`
	Jti    string `json:"jti"`
}

// UserInfo
// user_permit 中间件 中的返回参数 同时也是redis结构
type UserInfo struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	LastTime int64  `json:"lastTime"`
	LastIp   string `json:"lastIp"`
}
