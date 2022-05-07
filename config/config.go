package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	C           *Config
	RedisPrefix string
)

func init() {
	log.SetOutput(os.Stdout)

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Panicf("error when reading yaml: %v", err)
	}
	C = new(Config)
	if err := yaml.Unmarshal(data, C); err != nil {
		log.Panicf("error when unmarshal yaml: %v", err)
	}
	RedisPrefix = C.Redis.Prefix
	log.Print("config init")
}
