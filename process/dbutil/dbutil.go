package dbutil

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"openid/config"
	"os"
	"time"
)

var D *gorm.DB

func init() {
	m := config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", m.User, m.Pwd, m.Addr, m.Db, m.Charset)

	sqlLogger := logger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Millisecond * 200, // 慢 SQL 阈值
			LogLevel:                  logger.Info,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,                   // 禁用彩色打印
		},
	)

	var err error
	D, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})
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
