package dbutil

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"openid/config"
	"time"
)

var D *gorm.DB

func init() {
	m := config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", m.User, m.Pwd, m.Addr, m.Db, m.Charset)
	D, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("mysql error: %v", err)
	}

	sqlDb, err := D.DB()
	sqlDb.SetMaxOpenConns(m.MaxOpen)
	sqlDb.SetMaxIdleConns(m.MaxIdle)
	sqlDb.SetConnMaxLifetime(time.Duration(m.MaxLifetime) * time.Second)
	if err := sqlDb.Ping(); err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}
}
