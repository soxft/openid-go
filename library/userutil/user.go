package userutil

import (
	"database/sql"
	"openid/library/tool"
	"openid/process/mysqlutil"
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

// CheckPassword
// @description 验证用户登录
func CheckPassword(username, password string) (bool, int, string) {
	var row *sql.Stmt
	var err error

	if tool.IsEmail(username) {
		row, err = mysqlutil.D.Prepare("SELECT `id`,`salt`,`password` FROM `account` WHERE `email` = ? ")
	} else {
		row, err = mysqlutil.D.Prepare("SELECT `id`,`salt`,`password` FROM `account` WHERE `username` = ? ")
	}
	if err != nil {
		return false, 0, "system error"
	}
	res := row.QueryRow(username)
	if res.Err() == sql.ErrNoRows {
		return false, 0, "用户名或密码错误"
	} else if res.Err() != nil {
		return false, 0, "system error"
	} else {
		var id int
		var salt string
		var passwordDb string
		_ = res.Scan(&id, &salt, &passwordDb)
		if tool.Sha1(password+salt) != passwordDb {
			return false, 0, "用户名或密码错误"
		}
		return true, id, ""
	}
}
