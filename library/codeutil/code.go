package codeutil

import (
	"github.com/gomodule/redigo/redis"
	"openid/config"
	"openid/library/tool"
	"openid/redisutil"
)

type Coder interface {
	Create(length int) string
	Save(topic string, timeout int64, email string, code string) error
	Check(topic string, email string, code string) (bool, error)
	Consume(topic string, email string)
}

type VerifyCode struct {
}

func (c VerifyCode) Create(length int) string {
	return tool.RandStr(length)
}

func (c VerifyCode) Save(topic string, email string, code string, timeout int64) error {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	redisKey := config.C.Redis.Prefix + ":code:" + topic + ":" + tool.Md5(email)
	if _, err := _redis.Do("SETEX", redisKey, timeout, tool.Md5(code)); err != nil {
		return err
	}
	return nil
}

func (c VerifyCode) Check(topic string, email string, code string) (bool, error) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	redisKey := config.C.Redis.Prefix + ":code:" + topic + ":" + tool.Md5(email)
	if realCode, err := redis.String(_redis.Do("GET", redisKey)); err != nil {
		return false, err
	} else {
		if realCode == tool.Md5(code) {
			// delete key
			return true, nil
		}
	}
	return false, nil
}

// Consume
// @description: del code
func (c VerifyCode) Consume(topic string, email string) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	redisKey := config.C.Redis.Prefix + ":code:" + topic + ":" + tool.Md5(email)
	_, _ = _redis.Do("DEL", redisKey)
}
