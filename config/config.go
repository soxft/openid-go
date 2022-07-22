package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	C           *Config
	Server      ServerConfig
	Redis       RedisConfig
	Mysql       MysqlConfig
	Aliyun      AliyunConfig
	Jwt         JwtConfig
	Developer   DeveloperConfig
	Github      GithubConfig
	RedisPrefix string
)

func init() {
	data, err := ioutil.ReadFile("config.yaml")
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
	Aliyun = C.AliyunConfig
	Jwt = C.JwtConfig
	Developer = C.DeveloperConfig
	Github = C.GithubConfig
	RedisPrefix = C.RedisConfig.Prefix
}
