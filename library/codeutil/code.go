package codeutil

import (
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"openid/config"
	"openid/library/toolutil"
	"openid/process/redisutil"
	"strconv"
	"time"
)

func New() *VerifyCode {
	return &VerifyCode{}
}

// Create
// @description: create verify code
func (c VerifyCode) Create(length int) string {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))

	var code string
	for i := 0; i < length; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

// Save
// @description: save verify code 存储验证码
func (c VerifyCode) Save(topic string, email string, code string, timeout int64) error {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	redisKey := config.RedisPrefix + ":code:" + topic + ":" + toolutil.Md5(email)
	if _, err := _redis.Do("SETEX", redisKey, timeout, toolutil.Md5(code)); err != nil {
		return err
	}
	return nil
}

// Check
// @description: 判断验证码是否正确
func (c VerifyCode) Check(topic string, email string, code string) (bool, error) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	redisKey := config.RedisPrefix + ":code:" + topic + ":" + toolutil.Md5(email)
	if realCode, err := redis.String(_redis.Do("GET", redisKey)); err != nil {
		return false, err
	} else {
		if realCode == toolutil.Md5(code) {
			// delete key
			return true, nil
		}
	}
	return false, nil
}

// Consume
// @description: 消费(删除)验证码
func (c VerifyCode) Consume(topic string, email string) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	redisKey := config.RedisPrefix + ":code:" + topic + ":" + toolutil.Md5(email)
	_, _ = _redis.Do("DEL", redisKey)
}
