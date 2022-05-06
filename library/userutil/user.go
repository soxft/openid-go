package userutil

import (
	"database/sql"
	"openid/library/tool"
	"openid/mysqlutil"
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
	if exists, err := CheckUserName(username); err != nil {
		return false, "server error"
	} else {
		if exists {
			return false, "用户名已存在"
		}
	}
	if exists, err := CheckEmail(email); err != nil {
		return false, "server error"
	} else {
		if exists {
			return false, "邮箱已存在"
		}
	}

	return true, ""
}

// CheckUserName
// @description Check username if exists in database
func CheckUserName(username string) (bool, error) {
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

// CheckEmail
// @description Check email if exists in database
func CheckEmail(email string) (bool, error) {
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
