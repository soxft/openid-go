package dbutil

import (
	"fmt"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var D *gorm.DB

func Init() {
	m := config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", m.User, m.Pwd, m.Addr, m.Db, m.Charset)

	var logMode = logger.Warn
	if config.Server.Debug {
		logMode = logger.Info
	}

	sqlLogger := logger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200, // 慢 SQL 阈值
			LogLevel:                  logMode,                // 日志级别
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

	// 自动同步表结构
	_ = D.AutoMigrate(model.Account{})
	_ = D.AutoMigrate(model.App{})
	_ = D.AutoMigrate(model.OpenId{})
	_ = D.AutoMigrate(model.UniqueId{})
}
