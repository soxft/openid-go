package apputil

import "errors"

type AppBaseStruct struct {
	Id      int    `json:"id"`
	AppId   int    `json:"app_id"`
	AppName string `json:"app_name"`
}

type AppFullInfoStruct struct {
	Id         int    `json:"id"`
	AppId      int    `json:"app_id"`
	AppName    string `json:"app_name"`
	AppSecret  string `json:"app_secret"`
	AppGateway string `json:"app_gateway"`
	Time       int    `json:"create_time"`
}

var (
	ErrAppNotExist = errors.New("app not exist")
)
