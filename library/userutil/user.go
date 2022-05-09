package userutil

import (
	"database/sql"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
	"openid/library/tool"
	"openid/process/mysqlutil"
	"openid/process/redisutil"
	"strconv"
	"time"
)

// GenerateSalt
// @description Generate a random salt
func GenerateSalt() string {
	str := tool.RandStr(16)
	timestamp := time.Now().Unix()
	return tool.Md5(str + strconv.FormatInt(timestamp, 10))
}

func RegisterCheck(username, email string) (bool, string) {
	if exists, err := CheckUserNameExists(username); err != nil {
		return false, "server error"
	} else {
		if exists {
			return false, "用户名已存在"
		}
	}
	if exists, err := CheckEmailExists(email); err != nil {
		return false, "server error"
	} else if exists {
		return false, "邮箱已存在"
	}

	return true, ""
}

// CheckUserNameExists
// @description Check username if exists in database
func CheckUserNameExists(username string) (bool, error) {
	row, err := mysqlutil.D.Prepare("SELECT `id` FROM `account` WHERE `username` = ?")
	if err != nil {
		return false, err
	}
	res, err := row.Query(username)
	if err != nil {
		return false, err
	}
	defer func(res *sql.Rows) {
		_ = res.Close()
	}(res)

	if res.Next() {
		return true, nil
	}
	return false, nil
}

// CheckEmailExists
// @description Check email if exists in database
func CheckEmailExists(email string) (bool, error) {
	row, err := mysqlutil.D.Prepare("SELECT `id` FROM `account` WHERE `email` = ?")
	if err != nil {
		return false, err
	}
	res, err := row.Query(email)
	if err != nil {
		return false, err
	}
	defer func(res *sql.Rows) {
		_ = res.Close()
	}(res)

	if res.Next() {
		return true, nil
	}
	return false, nil
}

// CheckPassword
// @description 验证用户登录
func CheckPassword(username, password string) (int, error) {
	var row *sql.Stmt
	var err error

	if tool.IsEmail(username) {
		row, err = mysqlutil.D.Prepare("SELECT `id`,`salt`,`password` FROM `account` WHERE `email` = ? ")
	} else {
		row, err = mysqlutil.D.Prepare("SELECT `id`,`salt`,`password` FROM `account` WHERE `username` = ? ")
	}
	if err != nil {
		return 0, errors.New("system error")
	}
	res := row.QueryRow(username)
	if res.Err() == sql.ErrNoRows {
		return 0, errors.New("用户名或密码错误")
	} else if res.Err() != nil {
		return 0, errors.New("system error")
	}
	var id int
	var salt string
	var passwordDb string
	_ = res.Scan(&id, &salt, &passwordDb)
	if tool.Sha1(password+salt) != passwordDb {
		return 0, errors.New("用户名或密码错误")
	}
	return id, nil
}

// GetUserLast
// @description Get user last login time and ip
func GetUserLast(userId int) UserLastInfo {
	// get from redis
	_redis := redisutil.R.Get()
	_redisKey := config.RedisPrefix + ":user:last:" + strconv.Itoa(userId)

	var userLastInfo UserLastInfo
	row, err := redis.Values(_redis.Do("HGETALL", _redisKey))
	if err != nil {
		return userLastInfo
	}
	_ = redis.ScanStruct(row, &userLastInfo)
	_ = _redis.Close()
	return userLastInfo
}

// CheckPasswordByUserId
// @description 通过userid验证用户password
func CheckPasswordByUserId(userId int, password string) (bool, error) {
	var row *sql.Stmt
	var err error
	row, err = mysqlutil.D.Prepare("SELECT `salt`,`password` FROM `account` WHERE `id` = ? ")
	if err != nil {
		log.Printf("[ERROR] CheckPasswordByUserId: %s", err.Error())
		return false, errors.New("system error")
	}
	res := row.QueryRow(userId)
	if res.Err() == sql.ErrNoRows {
		return false, nil
	} else if res.Err() != nil {
		log.Printf("[ERROR] CheckPasswordByUserId: %s", res.Err())
		return false, errors.New("system error")
	}
	var salt string
	var passwordDb string
	_ = res.Scan(&salt, &passwordDb)
	if tool.Sha1(password+salt) != passwordDb {
		return false, nil
	}
	return true, nil
}
