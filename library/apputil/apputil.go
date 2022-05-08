package apputil

import (
	"database/sql"
	"errors"
	"html"
	"log"
	"math/rand"
	"openid/config"
	"openid/library/tool"
	"openid/process/mysqlutil"
	"strconv"
	"strings"
	"time"
)

func CheckName(name string) bool {
	if html.EscapeString(name) != name {
		return false
	}
	if len(name) < 2 || len(name) > 20 {
		return false
	}
	return true
}

func CreateApp(userId int, appName string) (bool, error) {
	if !CheckName(appName) {
		return false, errors.New("app name is not valid")
	}
	// 判断用户app数量是否超过限制
	counts, err := GetUserAppCount(userId)
	if err != nil {
		return false, err
	}
	if counts >= config.C.Developer.AppLimit {
		return false, errors.New("app count is over limit")
	}

	// 创建app
	appId, err := GenerateAppId()
	if err != nil {
		return false, err
	}
	appSecret := GenerateAppSecret()
	db, err := mysqlutil.D.Prepare("INSERT INTO `app` (`userId`,`appId`,`appName`,`appSecret`,`appGateway`,`time`) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[apputil] create app failed: %s", err.Error())
		return false, errors.New("server error")
	}
	_, err = db.Exec(userId, appId, appName, appSecret, "", time.Now().Unix())
	if err != nil {
		log.Printf("[apputil] create app failed: %s", err.Error())
		return false, errors.New("server error")
	}
	return true, nil
}

// GenerateAppId
// 创建唯一的appid
func GenerateAppId() (string, error) {
	timeUnix := time.Now().Unix()
	Tp := strconv.FormatInt(timeUnix, 10)
	// 随机数种子
	rand.Seed(time.Now().UnixNano())
	appId := time.Now().Format("20060102") + Tp[len(Tp)-4:] + strconv.Itoa(tool.RandInt(1000, 9999))
	if exists, err := checkAppIdExists(appId); err != nil {
		return "", err
	} else {
		if exists {
			return GenerateAppId()
		}
		return appId, nil
	}
}

// GenerateAppSecret
// 创建唯一的appSecret
func GenerateAppSecret() string {
	a := tool.Md5(time.Now().Format("20060102"))[:16]
	b := tool.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:16]
	c := tool.RandStr(16)
	d := tool.RandStr(16)
	return strings.Join([]string{a, b, c, d}, ".")
}

// CheckAppIdExists
// @description: check if appid exists
func checkAppIdExists(appid string) (bool, error) {
	db, err := mysqlutil.D.Prepare("select `id` from `app` where `appid` = ?")
	if err != nil {
		log.Printf("[ERROR] CheckAppIdExists: %s", err.Error())
		return false, errors.New("system error")
	}
	row := db.QueryRow(appid)
	var id int
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			log.Printf("[ERROR] CheckAppIdExists: %s", err.Error())
			return false, errors.New("system error")
		}
	}
	_ = db.Close()
	return false, nil
}

// GetUserAppList
// @description: 获取用户app列表
func GetUserAppList(userId, limit, offset int) ([]AppBaseStruct, error) {
	// 开始获取
	db, err := mysqlutil.D.Prepare("SELECT `id`,`appId`,`appName` FROM `app` WHERE `userId` = ? LIMIT ? OFFSET ?")
	if err != nil {
		log.Printf("[ERROR] AppGetList error: %s", err)
		return nil, errors.New("AppGetList error")
	}
	// process data
	row, err := db.Query(userId, limit, offset)
	if err != nil {
		log.Printf("[ERROR] AppGetList error: %s", err)
		return nil, errors.New("AppGetList error")
	}
	var appList []AppBaseStruct
	for row.Next() {
		var app AppBaseStruct
		err := row.Scan(&app.Id, &app.AppId, &app.AppName)
		if err != nil {
			log.Printf("[ERROR] AppGetList error: %s", err)
			return nil, errors.New("server error")
		}
		appList = append(appList, app)
	}

	_ = row.Close()
	_ = db.Close()
	return appList, nil
}

// GetUserAppCount
// 获取用户的app数量
func GetUserAppCount(userId int) (int, error) {
	db, err := mysqlutil.D.Prepare("SELECT COUNT(*) FROM `app` WHERE `userId` = ?")
	if err != nil {
		log.Printf("[ERROR] AppGetCount error: %s", err)
		return 0, errors.New("AppGetCount error")
	}
	var count int
	err = db.QueryRow(userId).Scan(&count)
	if err != nil {
		log.Printf("[ERROR] AppGetCount error: %s", err)
		return 0, errors.New("AppGetCount error")
	}
	return count, nil
}
