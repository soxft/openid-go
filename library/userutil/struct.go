package userutil

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

//type User struct {
//	Username string
//	Password string
//	Email    string
//	RegTime  int64
//	RegIp    string
//	LastTime int64
//	LastIp   string
//}
//
//type UserLastInfo struct {
//	LastIp   string
//	LastTime int64
//}

type JwtClaims struct {
	ID       string `json:"jti,omitempty"`
	ExpireAt int64  `json:"exp,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
	Issuer   string `json:"iss,omitempty"`
	Username string `json:"username"`
	UserId   int    `json:"userId"`
	Email    string `json:"email"`
	LastTime int64  `json:"lastTime"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	Username string `json:"username"`
	UserId   int    `json:"userId"`
	Email    string `json:"email"`
	LastTime int64  `json:"lastTime"`
}

var (
	ErrEmailExists    = errors.New("mailExists")
	ErrUsernameExists = errors.New("usernameExists")
	ErrPasswd         = errors.New("password not correct")
	ErrDatabase       = errors.New("database error")
	ErrJwtExpired     = errors.New("jwt is expired")
)
