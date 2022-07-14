package redisutil

import (
	"github.com/gomodule/redigo/redis"
	"github.com/soxft/openid/config"
	"log"
)

var R *redis.Pool

func init() {
	r := config.Redis
	R = &redis.Pool{
		MaxIdle:   r.MaxIdle,
		MaxActive: r.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Addr,
				redis.DialPassword(r.Pwd),
				redis.DialDatabase(r.Db),
			)
			if err != nil {
				log.Fatalf(err.Error())
			}
			return c, err
		},
	}
	if _, err := R.Get().Do("PING"); err != nil {
		log.Fatalf("redis connect error: %s", err.Error())
	}
}
