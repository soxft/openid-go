package helper

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/soxft/openid/library/toolutil"
	"github.com/soxft/openid/process/dbutil"
	"github.com/soxft/openid/process/redisutil"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

// GenerateToken
// @description: v1 获取token (用于跳转redirect_uri携带)
func GenerateToken(appId string, userId int) (string, error) {
	// check if exists in redis
	_redis := redisutil.R.Get()
	defer func() {
		_ = _redis.Close()
	}()

	a := toolutil.RandStr(10)
	b := toolutil.Md5(time.Now().Format("15:04:05"))[:10]
	c := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:10]

	token := a + "." + b + c + toolutil.RandStr(9)
	token = strings.ToLower(token)

	_redisKey := getTokenRedisKey(appId, token)

	if exists, err := redis.Bool(_redis.Do("EXISTS", _redisKey)); err != nil {
		log.Printf("[error] redis.Bool: %s", err.Error())
		return "", errors.New("system error")
	} else if !exists {
		// 不存在 则存入redis 并返回
		if _, err := _redis.Do("SETEX", _redisKey, 60*3, userId); err != nil {
			log.Printf("[ERROR] GetToken error: %s", err)
			return "", errors.New("server error")
		}
		return token, nil
	}
	// 存在
	return GenerateToken(appId, userId)
}

// generateOpenId
// 创建一个唯一的openId
func generateOpenId(appId string, userId int) (string, error) {
	a := toolutil.Md5(appId)[:10]
	b := toolutil.Md5(strconv.Itoa(userId))[:10]
	d := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	c := toolutil.Md5(toolutil.RandStr(16))

	openId := a + "." + b + "." + c + d
	openId = strings.ToLower(openId)
	err := dbutil.D.Model(&dbutil.OpenId{}).Where(dbutil.OpenId{OpenId: openId}).First(&dbutil.OpenId{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在
		err := dbutil.D.Create(&dbutil.OpenId{
			UserId: userId,
			AppId:  appId,
			OpenId: openId,
		}).Error
		if err != nil {
			log.Printf("[ERROR] generateOpenId error: %s", err)
			return "", errors.New("server error")
		}
		return openId, nil
	} else if err != nil {
		log.Printf("[ERROR] generateOpenId error: %s", err)
		return "", errors.New("server error")
	} else {
		// 存在
		return generateOpenId(appId, userId)
	}
}

// checkOpenIdExists
// 创建一个唯一的uniqueId
func generateUniqueId(userId, devUserId int) (string, error) {
	a := toolutil.Md5(strconv.Itoa(devUserId))[:10]
	b := toolutil.Md5(strconv.Itoa(userId))[:10]
	d := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	c := toolutil.Md5(toolutil.RandStr(16))

	uniqueId := a + "." + b + "." + c + d
	uniqueId = strings.ToLower(uniqueId)

	err := dbutil.D.Model(&dbutil.UniqueId{}).Where(dbutil.UniqueId{UniqueId: uniqueId}).First(&dbutil.UniqueId{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在
		err := dbutil.D.Create(&dbutil.UniqueId{
			UserId:    userId,
			DevUserId: devUserId,
			UniqueId:  uniqueId,
		}).Error

		if err != nil {
			log.Printf("[ERROR] generateUniqueId error: %s", err)
			return "", errors.New("server error")
		}
		return uniqueId, nil
	} else if err != nil {
		log.Printf("[ERROR] generateUniqueId error: %s", err)
		return "", errors.New("server error")
	} else {
		// 存在
		return generateUniqueId(userId, devUserId)
	}
}
