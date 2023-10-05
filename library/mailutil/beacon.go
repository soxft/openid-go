package mailutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/redisutil"
	"time"
)

// CreateBeacon
// @description: 创建邮件发送信标
func CreateBeacon(c *gin.Context, mail string, timeout time.Duration) error {
	_redis := redisutil.R

	unique := generateUnique(c)
	redisPrefix := config.Redis.Prefix

	ipKey := redisPrefix + ":beacon:ip:" + unique
	mailKey := redisPrefix + ":beacon:mail:" + toolutil.Md5(mail)

	_redis.SetEx(c, ipKey, "1", timeout)
	_redis.SetEx(c, mailKey, "1", timeout)
	return nil
}

// CheckBeacon
// @description: 检查邮件发送信标 避免频繁发信
func CheckBeacon(c *gin.Context, mail string) (bool, error) {
	_redis := redisutil.R

	unique := generateUnique(c)
	redisPrefix := config.Redis.Prefix

	ipExists, err := _redis.Exists(c, redisPrefix+":beacon:ip:"+unique).Result()
	if err != nil {
		return false, err
	}
	mailExists, err := _redis.Exists(c, redisPrefix+":beacon:mail:"+toolutil.Md5(mail)).Result()
	if err != nil {
		return false, err
	}
	if ipExists == 1 && mailExists == 1 {
		return true, nil
	}

	return false, nil
}

func generateUnique(c *gin.Context) string {
	// get user ip
	userIp := c.ClientIP()
	userAgent := c.Request.UserAgent()
	return toolutil.Md5(userIp + userAgent)
}
