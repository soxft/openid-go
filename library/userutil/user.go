package userutil

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"log"
	"openid/config"
	"openid/library/mailutil"
	"openid/library/toolutil"
	"openid/process/dbutil"
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
		return ErrUsernameExists
	}

	if exists, err := CheckEmailExists(email); err != nil {
		return err
	} else if exists {
		return ErrEmailExists
	}

	return nil
}

// CheckUserNameExists
// @description Check username if exists in database
func CheckUserNameExists(username string) (bool, error) {
	var ID int64
	err := dbutil.D.Model(&dbutil.Account{}).Select("id").Where("username = ?", username).First(ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		log.Printf("[ERROR] CheckUserNameExists: %s", err.Error())
		return false, errors.New("system error")
	}
	return true, nil
}

// CheckEmailExists
// @description Check email if exists in database
func CheckEmailExists(email string) (bool, error) {
	var ID string
	err := dbutil.D.Model(&dbutil.Account{}).Select("id").Where("email = ?", email).First(ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		log.Printf("[ERROR] CheckUserNameExists: %s", err.Error())
		return false, errors.New("system error")
	}
	return true, nil
}

// CheckPassword
// @description 验证用户登录
func CheckPassword(username, password string) (int, error) {
	var err error

	var account dbutil.Account
	if toolutil.IsEmail(username) {
		err = dbutil.D.Select("id, salt, password").Where("email = ?", username).Take(&account).Error
	} else {
		err = dbutil.D.Select("id, salt, password").Where("username = ?", username).Take(&account).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("用户名或密码错误")
	} else if err != nil {
		return 0, errors.New("system error")
	}
	if toolutil.Sha1(password+account.Salt) != account.Password {
		return 0, errors.New("用户名或密码错误")
	}
	return account.ID, nil
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
	// rewrite by gorm
	var account dbutil.Account
	err := dbutil.D.Model(&dbutil.Account{}).Select("id, salt, password").Where("id = ?", userId).Take(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, errors.New("system error")
	}
	if toolutil.Sha1(password+account.Salt) != account.Password {
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

func EmailChangeNotify(email string, timestamp time.Time) {
	_msg, _ := json.Marshal(mailutil.Mail{
		ToAddress: email,
		Subject:   "您的邮箱已修改",
		Content:   "您的邮箱已于" + timestamp.Format("2006-01-02 15:04:05") + "修改, 如果不是您本人操作, 请及时联系管理员",
		Typ:       "emailChangeNotify",
	})
	_ = queueutil.Q.Publish("mail", string(_msg), 5)
}
