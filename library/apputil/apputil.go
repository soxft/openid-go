package apputil

import (
	"errors"
	"log"
	"openid/process/mysqlutil"
)

func CheckName(name string) bool {
	return false
}

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
