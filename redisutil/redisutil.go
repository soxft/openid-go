package redisutil

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
)

var R *redis.Pool

func init() {
	RedisConfig := config.C.Redis
	R = &redis.Pool{
		MaxIdle:   RedisConfig.MaxIdle,
		MaxActive: RedisConfig.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisConfig.Addr,
				redis.DialPassword(RedisConfig.Pwd),
				redis.DialDatabase(RedisConfig.Db),
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
