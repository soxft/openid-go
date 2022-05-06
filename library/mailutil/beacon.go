package mailutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"openid/config"
	"openid/library/tool"
	"openid/redisutil"
)

// CreateBeacon
// @description: 创建邮件发送信标
func CreateBeacon(c *gin.Context, mail string, timeout int) error {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	unique := generateUnique(c)
	redisPrefix := config.C.Redis.Prefix

	ipKey := redisPrefix + ":beacon:ip:" + unique
	mailKey := redisPrefix + ":beacon:mail:" + tool.Md5(mail)
	_, _ = _redis.Do("SETEX", ipKey, timeout, "1")
	_, _ = _redis.Do("SETEX", mailKey, timeout, "1")
	return nil
}

// CheckBeacon
// @description: 检查邮件发送信标 避免频繁发信
func CheckBeacon(c *gin.Context, mail string) (bool, error) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	unique := generateUnique(c)
	redisPrefix := config.C.Redis.Prefix

	ipExists, err := _redis.Do("EXISTS", redisPrefix+":beacon:ip:"+unique)
	if err != nil {
		return false, err
	}
	mailExists, err := _redis.Do("EXISTS", redisPrefix+":beacon:mail:"+tool.Md5(mail))
	if err != nil {
		return false, err
	}
	if ipExists.(int64) == 1 && mailExists.(int64) == 1 {
		return true, nil
	}
	return false, nil
}

func generateUnique(c *gin.Context) string {
	// get user ip
	userIp := c.ClientIP()
	userAgent := c.Request.UserAgent()
	return tool.Md5(userIp + userAgent)
}
