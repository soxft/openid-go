package mq

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type MessageQueue interface {
	Publish(topic string, msg string, delay int64) error
	Subscribe(topic string, processes int, handler func(msg string))
}

type QueueArgs struct {
	redis      *redis.Client
	maxRetries int
	ctx        context.Context
}

type MsgArgs struct {
	Msg     string `json:"msg"`
	Retry   int    `json:"retry"`
	DelayAt int64  `json:"delay_at"`
}
