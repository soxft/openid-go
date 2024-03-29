package userutil

import (
	"encoding/json"
	"errors"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/library/mailutil"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/queueutil"
	"gorm.io/gorm"
	"log"
	"time"
)

// GenerateSalt
// @description Generate a random salt
//func GenerateSalt() string {
//	str := toolutil.RandStr(16)
//	timestamp := time.Now().Unix()
//	return toolutil.Md5(str + strconv.FormatInt(timestamp, 10))
//}

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
	err := dbutil.D.Model(&model.Account{}).Select("id").Where(model.Account{Username: username}).First(&ID).Error
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
	err := dbutil.D.Model(&model.Account{}).Select("id").Where(&model.Account{Email: email}).Take(&ID).Error
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
// if return < 0  error, pwd error or server error
// if return > 0  success, return user id
func CheckPassword(username, password string) (int, error) {
	var err error
	var account model.Account

	if toolutil.IsEmail(username) {
		err = dbutil.D.Select("id, password").Where(model.Account{Email: username}).Take(&account).Error
	} else {
		err = dbutil.D.Select("id, password").Where(model.Account{Username: username}).Take(&account).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, ErrPasswd
	} else if err != nil {
		log.Printf("[ERROR] CheckPassword: %v", err)
		return 0, ErrDatabase
	}
	if CheckPwd(password, account.Password) != nil {
		return 0, ErrPasswd
	}
	return account.ID, nil
}

// CheckPasswordByUserId
// @description 通过userid验证用户password
//func CheckPasswordByUserId(userId int, password string) (bool, error) {
//	// rewrite by gorm
//	var account model.Account
//
//	if err := dbutil.D.Model(model.Account{}).Select("id, salt, password").
//		Where(model.Account{ID: userId}).Take(&account).Error; errors.Is(err, gorm.ErrRecordNotFound) {
//		return false, nil
//	} else if err != nil {
//		return false, errors.New("system error")
//	}
//	if CheckPwd(password, account.Password) != nil {
//		return false, nil
//	}
//	return true, nil
//}

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
