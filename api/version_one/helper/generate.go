package helper

import (
	"database/sql"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/library/toolutil"
	"openid/process/mysqlutil"
	"openid/process/redisutil"
	"strconv"
	"strings"
	"time"
)

// GenerateToken
// @description: v1 获取token (用于跳转redirect_uri携带)
func GenerateToken(appId, userId int) (string, error) {
	// check if exists in redis
	_redis := redisutil.R.Get()
	defer func() {
		_ = _redis.Close()
	}()

	a := toolutil.RandStr(10)
	b := toolutil.Md5(time.Now().Format("15:04:05"))[:10]
	c := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:10]

	token := a + "_" + b + toolutil.RandStr(9) + c
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
func generateOpenId(appId int, userId int) (string, error) {
	a := toolutil.Md5(strconv.Itoa(appId))[:10]
	b := toolutil.Md5(strconv.Itoa(userId))[:10]
	c := toolutil.Md5(time.Now().Format("15:04:05"))[:10]
	d := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:10]
	e := toolutil.RandStr(22)
	openId := a + "_" + b + "_" + c + d + e
	openId = strings.ToLower(openId)

	db, _ := mysqlutil.D.Prepare("SELECT `id` FROM `openId` WHERE `openId` = ? ")
	defer func() {
		_ = db.Close()
	}()

	var id int
	err := db.QueryRow(openId).Scan(&id)
	if err == sql.ErrNoRows {
		// 不存在
		db, _ := mysqlutil.D.Prepare("INSERT INTO `openId` (`userId`,`appId`, `openId`,`time`) VALUES (?, ?, ?, ?)")
		_, err := db.Exec(userId, appId, openId, time.Now().Unix())
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
	c := toolutil.Md5(time.Now().Format("15:04:05"))[:10]
	d := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:10]
	e := toolutil.RandStr(22)
	uniqueId := a + "_" + b + "_" + c + d + e
	uniqueId = strings.ToLower(uniqueId)

	db, _ := mysqlutil.D.Prepare("SELECT `id` FROM `uniqueId` WHERE `uniqueId` = ? ")
	defer func() {
		_ = db.Close()
	}()

	var id int
	err := db.QueryRow(uniqueId).Scan(&id)
	if err == sql.ErrNoRows {
		// 不存在
		db, _ := mysqlutil.D.Prepare("INSERT INTO `uniqueId` (`userId`, `DevUserId`, `uniqueId`,`time`) VALUES (?, ?, ?, ?)")
		_, err := db.Exec(userId, devUserId, uniqueId, time.Now().Unix())
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
