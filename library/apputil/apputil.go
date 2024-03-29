package apputil

import (
	"errors"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/dbutil"
	"gorm.io/gorm"
	"html"
	"log"
	"math/rand"
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
	//log.Printf("[apputil] check gateway: %s", gateway)
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
		// 是否存在空格
		if strings.Contains(part, " ") {
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

	if result := dbutil.D.Create(&model.App{
		UserId:     userId,
		AppId:      appId,
		AppName:    appName,
		AppSecret:  generateAppSecret(),
		AppGateway: "",
	}); result.Error != nil || result.RowsAffected == 0 {
		log.Printf("[apputil] create app failed: %s", result.Error.Error())
		return false, errors.New("创建应用时发生错误, 请稍后再试")
	}

	return true, nil
}

// DeleteUserApp
// 删除用户App
func DeleteUserApp(appId string) (bool, error) {
	// 开启 事物
	err := dbutil.D.Transaction(func(tx *gorm.DB) error {
		var App model.App
		var OpenId model.OpenId

		// 删除app表内数据
		err := tx.Where(model.App{AppId: appId}).Delete(&App).Error
		if err != nil {
			return errors.New("system error")
		}
		// 删除openId表内数据
		err = tx.Where(model.OpenId{AppId: appId}).Delete(&OpenId).Error
		if err != nil {
			return errors.New("system error")
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] DeleteUserApp error: %s", err)
		return false, err
	}
	return true, nil
}

// GetUserAppList
// @description: 获取用户app列表
func GetUserAppList(userId, limit, offset int) ([]AppBaseStruct, error) {
	// 开始获取
	var appList []AppBaseStruct
	var appListRaw []model.App
	err := dbutil.D.Model(model.App{}).Select("id, app_id, app_name, create_at").Where(model.App{UserId: userId}).Order("id desc").Limit(limit).Offset(offset).Find(&appListRaw).Error
	if err != nil {
		log.Printf("[ERROR] GetUserAppList error: %s", err)
		return nil, errors.New("GetUserAppList error")
	}
	for _, app := range appListRaw {
		appList = append(appList, AppBaseStruct{
			Id:       app.ID,
			AppId:    app.AppId,
			AppName:  app.AppName,
			CreateAt: app.CreateAt,
		})
	}
	return appList, nil
}

// GetUserAppCount
// 获取用户的app数量
func GetUserAppCount(userId int) (int, error) {
	var count int64
	err := dbutil.D.Model(&model.App{}).Where(model.App{UserId: userId}).Count(&count).Error
	if err != nil {
		log.Printf("[ERROR] GetUserAppCount error: %s", err)
		return 0, errors.New("GetUserAppCount error")
	}

	countInt := int(count)
	return countInt, nil
}

func GetAppInfo(appId string) (AppFullInfoStruct, error) {
	var appInfo AppFullInfoStruct
	var appInfoRaw model.App

	err := dbutil.D.Model(&model.App{}).Select("id, user_id, app_id, app_name, app_secret, app_gateway, create_at").Where(model.App{AppId: appId}).Take(&appInfoRaw).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appInfo, ErrAppNotExist
	} else if err != nil {
		log.Printf("[ERROR] GetAppInfo error: %s", err)
		return appInfo, errors.New("server error")
	}
	appInfo = AppFullInfoStruct{
		Id:         appInfoRaw.ID,
		AppUserId:  appInfoRaw.UserId,
		AppId:      appInfoRaw.AppId,
		AppName:    appInfoRaw.AppName,
		AppSecret:  appInfoRaw.AppSecret,
		AppGateway: appInfoRaw.AppGateway,
		CreateAt:   appInfoRaw.CreateAt,
	}
	return appInfo, nil
}

// CheckAppSecret
// @description: 检查appSecret
func CheckAppSecret(appId string, appSecret string) error {
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
	} else if exists {
		return generateAppId()
	}
	return appId, nil
}

// CheckIfUserApp
// 判断是否为该用户的app
func CheckIfUserApp(appId string, userId int) (bool, error) {
	var appUserId int
	err := dbutil.D.Model(&model.App{}).Select("user_id").Where(model.App{AppId: appId}).Take(&appUserId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[ERROR] CheckIfUserApp error: %s", err)

		return false, errors.New("app not exist")
	} else if err != nil {
		log.Printf("[ERROR] CheckIfUserApp error: %s", err)

		return false, errors.New("server error")
	}

	if appUserId != userId {
		return false, nil
	}
	return true, nil
}

// ReGenerateSecret
// 重新生成新的 appSecret
func ReGenerateSecret(appId string) (string, error) {
	appSecret := generateAppSecret()
	err := dbutil.D.Model(&model.App{}).
		Where(model.App{AppId: appId}).
		Updates(model.App{AppSecret: appSecret}).Error
	if err != nil {
		log.Printf("[ERROR] ReGenerateAppSecret error: %s", err)
		return "", errors.New("server error")
	}
	return appSecret, nil
}

// generateAppSecret
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
	var ID int64
	err := dbutil.D.Model(model.App{}).Select("id").Where(model.App{AppId: appid}).Take(&ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		log.Printf("[ERROR] CheckAppIdExists: %s", err)
		return false, errors.New("system error")
	}
	return true, nil
}

// CheckRedirectUriIsMatchUserGateway
// @description: 检测回调地址是否匹配用户网关
func CheckRedirectUriIsMatchUserGateway(redirectUri string, GateWay string) bool {
	for _, gateway := range strings.Split(GateWay, ",") {
		if redirectUri == gateway {
			return true
		}
	}

	return false
}
