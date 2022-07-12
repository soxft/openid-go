package helper

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"log"
	"openid/config"
	"openid/library/apputil"
	"openid/library/toolutil"
	"openid/process/dbutil"
	"openid/process/redisutil"
)

// GetUserIdByToken
// 通过Token和appid 获取用户ID
func GetUserIdByToken(appId string, token string) (int, error) {
	_redis := redisutil.R.Get()
	defer func() {
		_ = _redis.Close()
	}()

	userId, err := redis.Int(_redis.Do("GET", getTokenRedisKey(appId, token)))
	if err != nil {
		if err == redis.ErrNil {
			return 0, ErrTokenNotExists
		}
		log.Printf("[ERROR] GetUserIdByToken error: %s", err)
		return 0, errors.New("server error")
	}
	return userId, nil
}

// GetUserIds
// @description 获取用户ID
func GetUserIds(appId string, userId int) (UserIdsStruct, error) {
	openId, err := getUserOpenId(appId, userId)
	if err != nil {
		return UserIdsStruct{}, err
	}
	appInfo, err := apputil.GetAppInfo(appId)
	if err != nil {
		return UserIdsStruct{}, err
	}
	uniqueId, err := getUserUniqueId(userId, appInfo.AppUserId)
	if err != nil {
		return UserIdsStruct{}, err
	}

	return UserIdsStruct{
		OpenId:   openId,
		UniqueId: uniqueId,
	}, nil
}

func DeleteToken(appId string, token string) error {
	_redis := redisutil.R.Get()
	defer func() {
		_ = _redis.Close()
	}()

	_, err := _redis.Do("DEL", getTokenRedisKey(appId, token))
	if err != nil {
		log.Printf("[ERROR] DeleteToken error: %s", err)
		return errors.New("server error")
	}
	return nil
}

// GetUserOpenId
// 获取 用户openID
func getUserOpenId(appId string, userId int) (string, error) {
	var openId string
	err := dbutil.D.Model(&dbutil.OpenId{}).Where(dbutil.OpenId{AppId: appId, UserId: userId}).Select("open_id").First(&openId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return generateOpenId(appId, userId)
	} else if err != nil {
		log.Printf("[ERROR] GetUserOpenId error: %s", err)
		return "", errors.New("server error")
	}
	return openId, nil
}

// getUserUniqueId
// 获取用户UniqueId
func getUserUniqueId(userId, DevUserId int) (string, error) {
	var uniqueId string
	err := dbutil.D.Model(&dbutil.UniqueId{}).Where(dbutil.UniqueId{UserId: userId, DevUserId: DevUserId}).Select("unique_id").First(&uniqueId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return generateUniqueId(userId, DevUserId)
	} else if err != nil {
		log.Printf("[ERROR] GetUserUniqueId error: %s", err)
		return "", errors.New("server error")
	}
	return uniqueId, nil
}

func getTokenRedisKey(appId string, token string) string {
	return config.RedisPrefix + ":app:" + toolutil.Md5(appId) + ":" + toolutil.Md5(token)
}
