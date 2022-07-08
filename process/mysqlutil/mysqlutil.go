package mysqlutil

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"openid/config"
	"time"
)

var D *sql.DB

func init() {
	m := config.Mysql
	var err error
	D, err = sql.Open("mysql", m.User+":"+m.Pwd+"@tcp("+m.Addr+")/"+m.Db+"?charset="+m.Charset)
	if err != nil {
		log.Fatalf("mysql error: %v", err)
	}
	D.SetMaxOpenConns(m.MaxOpen)
	D.SetMaxIdleConns(m.MaxIdle)
	D.SetConnMaxLifetime(time.Duration(m.MaxLifetime) * time.Second)
	if err := D.Ping(); err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}
}
