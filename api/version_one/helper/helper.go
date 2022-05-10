package helper

import (
	"database/sql"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
	"openid/library/apputil"
	"openid/library/toolutil"
	"openid/process/mysqlutil"
	"openid/process/redisutil"
	"strconv"
)

// GetUserIdByToken
// 通过Token和appid 获取用户ID
func GetUserIdByToken(appId int, token string) (int, error) {
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
func GetUserIds(appId, userId int) (UserIdsStruct, error) {
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

func DeleteToken(appId int, token string) error {
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
func getUserOpenId(appId int, userId int) (string, error) {
	db, err := mysqlutil.D.Prepare("SELECT `openId` FROM `openId` WHERE `userId` = ? AND `appId` = ?")
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		log.Printf("[ERROR] GetUserOpenId error: %s", err)
		return "", errors.New("server error")
	}
	var openId string
	if err := db.QueryRow(userId, appId).Scan(&openId); err != nil {
		if err == sql.ErrNoRows {
			// 无结果, 则创建一个
			return generateOpenId(appId, userId)
		}
		log.Printf("[ERROR] GetUserOpenId error: %s", err)
		return "", errors.New("server error")
	}
	return openId, nil
}

// getUserUniqueId
// 获取用户UniqueId
func getUserUniqueId(userId, DevUserId int) (string, error) {
	db, err := mysqlutil.D.Prepare("SELECT `uniqueId` FROM `uniqueId` WHERE `userId` = ? AND `devUserId` = ?")
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		log.Printf("[ERROR] GetUserOpenId error: %s", err)
		return "", errors.New("server error")
	}
	var openId string
	if err := db.QueryRow(userId, DevUserId).Scan(&openId); err != nil {
		if err == sql.ErrNoRows {
			// 无结果, 则创建一个
			return generateUniqueId(userId, DevUserId)
		}
		log.Printf("[ERROR] GetUserOpenId error: %s", err)
		return "", errors.New("server error")
	}
	return openId, nil
}

func getTokenRedisKey(appId int, token string) string {
	return config.RedisPrefix + ":app:" + toolutil.Md5(strconv.Itoa(appId)) + ":" + toolutil.Md5(token)
}
