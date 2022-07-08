package mysqlutil

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("mysql error: %v", err)
	}

	D, err := db.DB()
	D.SetMaxOpenConns(m.MaxOpen)
	D.SetMaxIdleConns(m.MaxIdle)
	D.SetConnMaxLifetime(time.Duration(m.MaxLifetime) * time.Second)
	if err := D.Ping(); err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}
}
