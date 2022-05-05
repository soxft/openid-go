package redisutil

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
)

var R *redis.Pool

func init() {
	RedisC := config.C.Redis
	R = &redis.Pool{
		MaxIdle:   RedisC.MaxIdle,
		MaxActive: RedisC.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisC.Addr,
				redis.DialPassword(RedisC.Pwd),
				redis.DialDatabase(RedisC.Db),
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
