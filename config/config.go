package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	C           *Config
	Server      ServerConfig
	Redis       RedisConfig
	Mysql       MysqlConfig
	Smtp        SmtpConfig
	Aliyun      AliyunConfig
	Jwt         JwtConfig
	Developer   DeveloperConfig
	RedisPrefix string
)

func init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Panicf("error when reading yaml: %v", err)
	}
	C = &Config{}
	if err := yaml.Unmarshal(data, C); err != nil {
		log.Panicf("error when unmarshal yaml: %v", err)
	}

	Server = C.ServerConfig
	Redis = C.RedisConfig
	Mysql = C.MysqlConfig
	Smtp = C.SmtpConfig
	Aliyun = C.AliyunConfig
	Jwt = C.JwtConfig
	Developer = C.DeveloperConfig
	RedisPrefix = C.RedisConfig.Prefix
}
