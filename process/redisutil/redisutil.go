package redisutil

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/soxft/openid-go/config"
	"log"
)

var R *redis.Client

func Init() {
	log.Printf("[INFO] Redis trying connect to tcp://%s/%d", config.Redis.Addr, config.Redis.Db)

	r := config.Redis

	R := redis.NewClient(&redis.Options{
		Addr:           r.Addr,
		Password:       r.Pwd, // no password set
		DB:             r.Db,  // use default DB
		MaxIdleConns:   r.MaxIdle,
		MaxActiveConns: r.MaxActive,
		MaxRetries:     r.MaxRetries,
	})

	if err := R.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("[ERROR] Redis connect error: %s", err)
	}

	log.Printf("[INFO] Redis connect success")
}
