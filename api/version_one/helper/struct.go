package helper

import "errors"

type ApiErr = error

var (
	ErrTokenNotExists = errors.New("token not exists")

	ErrOpenIdExists   = errors.New("openId exists")
	ErrUniqueIdExists = errors.New("uniqueId exists")
)

type UserIdsStruct struct {
	OpenId   string `json:"openId"`
	UniqueId string `json:"uniqueId"`
}
