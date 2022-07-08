package apputil

import (
	"database/sql"
	"errors"
	"html"
	"log"
	"math/rand"
	"openid/config"
	"openid/library/toolutil"
	"openid/process/mysqlutil"
	"strconv"
	"strings"
	"time"
)

// CheckName
// 检测应用名称合法性
func CheckName(name string) bool {
	if html.EscapeString(name) != name {
		return false
	}
	if len(name) < 2 || len(name) > 20 {
		return false
	}
	return true
}

// CheckGateway
// 检测应用网关合法性
func CheckGateway(gateway string) bool {
	if len(gateway) < 4 || len(gateway) > 200 {
		return false
	}
	// 如果不包含 小数点
	if !strings.Contains(gateway, ".") {
		return false
	}
	// 是否包含特殊字符
	if strings.ContainsAny(gateway, "~!@#$%^&*()_+=|\\{}[];'\",/<>?") {
		return false
	}
	// 不能 以 . 开头 或 结尾
	if strings.HasPrefix(gateway, ".") || strings.HasSuffix(gateway, ".") {
		return false
	}
	// 以小数点分割, 每段不能超过 63 个字符
	parts := strings.Split(gateway, ".")
	for _, part := range parts {
		if len(part) > 63 {
			return false
		}
	}
	return true
}

func CreateApp(userId int, appName string) (bool, error) {
	if userId == 0 {
		return false, errors.New("userId is invalid")
	}
	if !CheckName(appName) {
		return false, errors.New("app name is invalid")
	}
	// 判断用户app数量是否超过限制
	counts, err := GetUserAppCount(userId)
	if err != nil {
		return false, err
	}
	if counts >= config.Developer.AppLimit {
		return false, errors.New("the number of app exceeds the limit")
	}

	// 创建app
	appId, err := generateAppId()
	if err != nil {
		return false, err
	}
	appSecret := generateAppSecret()
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

// DeleteUserApp
// 删除用户App
func DeleteUserApp(appId int) (bool, error) {
	// 开启 事物
	db, err := mysqlutil.D.Begin()
	if err != nil {
		log.Printf("[ERROR] DeleteUserApp error 1: %s", err)
		return false, errors.New("DeleteUserApp error")
	}
	// 删除app表内数据
	_, err = db.Exec("DELETE FROM `app` WHERE `appId` = ?", appId)
	if err != nil {
		log.Printf("[ERROR] DeleteUserApp error 2: %s", err)
		_ = db.Rollback()
		return false, errors.New("system error")
	}
	// 删除openId表内数据
	_, err = db.Exec("DELETE FROM `openId` WHERE `appId` = ?", appId)
	if err != nil {
		log.Printf("[ERROR] DeleteUserApp error 3: %s", err)
		_ = db.Rollback()
		return false, errors.New("system error")
	}
	_ = db.Commit()
	return true, nil
}

// GetUserAppList
// @description: 获取用户app列表
func GetUserAppList(userId, limit, offset int) ([]AppBaseStruct, error) {
	// 开始获取
	db, err := mysqlutil.D.Prepare("SELECT `id`,`appId`,`appName`,`time` FROM `app` WHERE `userId` = ? ORDER BY `id` DESC LIMIT ? OFFSET ?")
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
		err := row.Scan(&app.Id, &app.AppId, &app.AppName, &app.CreateTime)
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

func GetAppInfo(appId int) (AppFullInfoStruct, error) {
	var appInfo AppFullInfoStruct
	db, err := mysqlutil.D.Prepare("SELECT `id`,`userId`,`appId`,`appName`,`appSecret`,`appGateway`,`time` FROM `app` WHERE `appId` = ?")
	if err != nil {
		log.Printf("[ERROR] AppGetInfo error: %s", err)
		return appInfo, errors.New("server error")
	}
	err = db.QueryRow(appId).Scan(&appInfo.Id, &appInfo.AppUserId, &appInfo.AppId, &appInfo.AppName, &appInfo.AppSecret, &appInfo.AppGateway, &appInfo.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return appInfo, ErrAppNotExist
		} else {
			log.Printf("[ERROR] GetAppInfo error: %s", err.Error())
			return appInfo, errors.New("server error")
		}
	}
	return appInfo, nil
}

// CheckAppSecret
// @description: 检查appSecret
func CheckAppSecret(appId int, appSecret string) error {
	appInfo, err := GetAppInfo(appId)
	if err != nil {
		return err
	}
	if appInfo.AppSecret != appSecret {
		return ErrAppSecretNotMatch
	}
	return nil
}

// GenerateAppId
// 创建唯一的appid
func generateAppId() (string, error) {
	timeUnix := time.Now().Unix()
	Tp := strconv.FormatInt(timeUnix, 10)
	// 随机数种子
	rand.Seed(time.Now().UnixNano())
	appId := time.Now().Format("20060102") + Tp[len(Tp)-4:] + strconv.Itoa(toolutil.RandInt(4))
	if exists, err := checkAppIdExists(appId); err != nil {
		return "", err
	} else {
		if exists {
			return generateAppId()
		}
		return appId, nil
	}
}

// CheckIfUserApp
// 判断是否为该用户的app
func CheckIfUserApp(appId, userId int) (bool, error) {
	db, err := mysqlutil.D.Prepare("SELECT `userId` FROM `app` WHERE `appId` = ?")
	if err != nil {
		log.Printf("[ERROR] CheckIfUserApp error: %s", err)
		return false, errors.New("CheckIfUserApp error")
	}
	var userIds int
	err = db.QueryRow(appId).Scan(&userIds)
	if err != nil {
		// 无数据
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("[ERROR] CheckIfUserApp error: %s", err)
		return false, errors.New("CheckIfUserApp error")
	}
	if userIds != userId {
		return false, nil
	}
	return true, nil
}

// GenerateAppSecret
// 创建唯一的appSecret
func generateAppSecret() string {
	a := toolutil.Md5(time.Now().Format("20060102"))[:16]
	b := toolutil.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))[:16]
	c := toolutil.RandStr(16)
	d := toolutil.RandStr(16)
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
