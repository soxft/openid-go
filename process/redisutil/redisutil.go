package redisutil

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/soxft/openid-go/config"
	"log"
	"time"
)

var RDB *redis.Client

func Init() {
	log.Printf("[INFO] Redis trying connect to tcp://%s/%d", config.Redis.Addr, config.Redis.Db)

	r := config.Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:            r.Addr,
		Password:        r.Pwd, // no password set
		DB:              r.Db,  // use default DB
		MinIdleConns:    r.MinIdle,
		MaxIdleConns:    r.MaxIdle,
		MaxRetries:      r.MaxRetries,
		ConnMaxLifetime: 5 * time.Minute,
		MaxActiveConns:  r.MaxActive,
	})

	RDB = rdb

	// test redis
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("[ERROR] Redis ping failed: %v", err)
	}

	log.Printf("[INFO] Redis init success, pong: %s ", pong)
}
