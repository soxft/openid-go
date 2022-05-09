package userutil

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
	"openid/library/mailutil"
	"openid/library/toolutil"
	"openid/process/mysqlutil"
	"openid/process/queueutil"
	"openid/process/redisutil"
	"strconv"
	"time"
)

// GenerateSalt
// @description Generate a random salt
func GenerateSalt() string {
	str := toolutil.RandStr(16)
	timestamp := time.Now().Unix()
	return toolutil.Md5(str + strconv.FormatInt(timestamp, 10))
}

// RegisterCheck
// @description Check users email or user if already exists
func RegisterCheck(username, email string) error {
	if exists, err := CheckUserNameExists(username); err != nil {
		return err
	} else if exists {
		return ErrEmailExists
	}

	if exists, err := CheckEmailExists(email); err != nil {
		return err
	} else if exists {
		return ErrUsernameExists
	}

	return nil
}

// CheckUserNameExists
// @description Check username if exists in database
func CheckUserNameExists(username string) (bool, error) {
	row, err := mysqlutil.D.Prepare("SELECT `id` FROM `account` WHERE `username` = ?")
	if err != nil {
		log.Printf("CheckUserNameExists: %s", err.Error())
		return false, errors.New("server error")
	}
	res := row.QueryRow(username)
	var id int
	err = res.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("[ERROR] CheckEmailExists err: %v", err)
		return false, errors.New("server error")
	}
	return true, nil
}

// CheckEmailExists
// @description Check email if exists in database
func CheckEmailExists(email string) (bool, error) {
	row, err := mysqlutil.D.Prepare("SELECT `id` FROM `account` WHERE `email` = ?")
	if err != nil {
		log.Printf("[ERROR] CheckEmailExists err: %v", err)
		return false, errors.New("server error")
	}
	res := row.QueryRow(email)
	var id int
	err = res.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("[ERROR] CheckEmailExists err: %v", err)
		return false, errors.New("server error")
	}
	return true, nil
}

// CheckPassword
// @description 验证用户登录
func CheckPassword(username, password string) (int, error) {
	var row *sql.Stmt
	var err error

	if toolutil.IsEmail(username) {
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
	if toolutil.Sha1(password+salt) != passwordDb {
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
	if toolutil.Sha1(password+salt) != passwordDb {
		return false, nil
	}
	return true, nil
}

func PasswordChangeNotify(email string, timestamp time.Time) {
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   "您的密码已修改",
		Content:   "您的密码已于" + timestamp.Format("2006-01-02 15:04:05") + "修改, 如果不是您本人操作, 请及时联系管理员",
		Typ:       "passwordChangeNotify",
	})
	_ = queueutil.Q.Publish("mail", string(_msg), 5)
}
