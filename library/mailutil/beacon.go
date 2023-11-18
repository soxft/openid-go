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
	_redis := redisutil.RDB

	ipKey, mailKey := getRKeys(c, mail)

	_redis.SetEx(c, ipKey, "1", timeout)
	_redis.SetEx(c, mailKey, "1", timeout)
	return nil
}

// DeleteBeacon
// @description: 手动删除邮件创建新信标
func DeleteBeacon(c *gin.Context, mail string) {
	_redis := redisutil.RDB

	ipKey, mailKey := getRKeys(c, mail)

	_redis.Del(c, ipKey)
	_redis.Del(c, mailKey)
}

// CheckBeacon
// @description: 检查邮件发送信标 避免频繁发信
func CheckBeacon(c *gin.Context, mail string) (bool, error) {
	_redis := redisutil.RDB

	ipKey, mailKey := getRKeys(c, mail)

	if ipExists, err := _redis.Exists(c, ipKey).Result(); err != nil || ipExists != 1 {
		return false, err
	}

	if mailExists, err := _redis.Exists(c, mailKey).Result(); err != nil || mailExists != 1 {
		return false, err
	}

	return true, nil
}

func generateUnique(c *gin.Context) string {
	// get user ip
	userIp := c.ClientIP()
	userAgent := c.Request.UserAgent()
	return toolutil.Md5(userIp + userAgent)
}

func getRKeys(c *gin.Context, mail string) (string, string) {
	unique := generateUnique(c)
	redisPrefix := config.Redis.Prefix

	ipKey := redisPrefix + ":beacon:ip:" + unique
	mailKey := redisPrefix + ":beacon:mail:" + toolutil.Md5(mail)

	return ipKey, mailKey
}
