package helper

import (
	"context"
	"errors"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/redisutil"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

// GenerateToken
// @description: v1 获取token (用于跳转redirect_uri携带)
func GenerateToken(ctx context.Context, appId string, userId int) (string, error) {
	// check if exists in redis
	_redis := redisutil.RDB

	a := toolutil.RandStr(10)
	b := toolutil.Md5(time.Now().Format("15:04:05"))[:10]
	c := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:10]

	token := a + "." + b + c + toolutil.RandStr(9)
	token = strings.ToLower(token)

	_redisKey := getTokenRedisKey(appId, token)

	if exists, err := _redis.Exists(ctx, _redisKey).Result(); err != nil {
		log.Printf("[error] redis.Bool: %s", err.Error())
		return "", errors.New("system error")
	} else if exists == 0 {
		// 不存在 则存入redis 并返回
		if _, err := _redis.SetEx(ctx, _redisKey, userId, 3*time.Minute).Result(); err != nil {
			log.Printf("[ERROR] GetToken error: %s", err)
			return "", errors.New("server error")
		}
		return token, nil
	}
	// 存在
	return GenerateToken(ctx, appId, userId)
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
	err := dbutil.D.Model(&model.OpenId{}).Where(model.OpenId{OpenId: openId}).First(&model.OpenId{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在
		err := dbutil.D.Create(&model.OpenId{
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

	err := dbutil.D.Model(&model.UniqueId{}).Where(model.UniqueId{UniqueId: uniqueId}).First(&model.UniqueId{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在
		err := dbutil.D.Create(&model.UniqueId{
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
