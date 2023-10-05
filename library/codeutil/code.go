package codeutil

import (
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/redisutil"
	"math/rand"
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
// timeout: expire time (second)
func (c VerifyCode) Save(topic string, email string, code string, timeout time.Duration) error {
	_redis := redisutil.R

	redisKey := config.RedisPrefix + ":code:" + topic + ":" + toolutil.Md5(email)

	return _redis.SetEx(c.ctx, redisKey, toolutil.Md5(code), timeout).Err()
}

// Check
// @description: 判断验证码是否正确
func (c VerifyCode) Check(topic string, email string, code string) (bool, error) {
	_redis := redisutil.R

	redisKey := config.RedisPrefix + ":code:" + topic + ":" + toolutil.Md5(email)
	if realCode, err := _redis.Get(c.ctx, redisKey).Result(); err != nil {
		return false, err
	} else if realCode == toolutil.Md5(code) {
		// delete key
		return true, nil
	}

	return false, nil
}

// Consume
// @description: 消费(删除)验证码
func (c VerifyCode) Consume(topic string, email string) {
	_redis := redisutil.R

	redisKey := config.RedisPrefix + ":code:" + topic + ":" + toolutil.Md5(email)

	_redis.Del(c.ctx, redisKey)
}
