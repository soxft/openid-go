package helper

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/apputil"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/redisutil"
	"gorm.io/gorm"
	"log"
)

// GetUserIdByToken
// 通过Token和appid 获取用户ID
func GetUserIdByToken(ctx context.Context, appId string, token string) (int, error) {
	_redis := redisutil.RDB

	userId, err := _redis.Get(ctx, getTokenRedisKey(appId, token)).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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

func DeleteToken(ctx context.Context, appId string, token string) error {
	_redis := redisutil.RDB

	if err := _redis.Del(ctx, getTokenRedisKey(appId, token)).Err(); err != nil {
		log.Printf("[ERROR] DeleteToken error: %s", err)
		return errors.New("server error")
	}
	return nil
}

// GetUserOpenId
// 获取 用户openID
func getUserOpenId(appId string, userId int) (string, error) {
	var openId string
	err := dbutil.D.Model(&model.OpenId{}).Where(model.OpenId{AppId: appId, UserId: userId}).Select("open_id").First(&openId).Error

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
	err := dbutil.D.Model(&model.UniqueId{}).Where(model.UniqueId{UserId: userId, DevUserId: DevUserId}).Select("unique_id").First(&uniqueId).Error

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
