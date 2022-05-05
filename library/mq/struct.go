package mq

import (
	"github.com/gomodule/redigo/redis"
)

type MessageQueue interface {
	Publish(topic string, msg string, delay int) error
	Subscribe(topic string, handler func(msg string))
}

type QueueArgs struct {
	redis      *redis.Pool
	maxRetries int
}

type MsgArgs struct {
	Msg   string `json:"msg"`
	Retry int    `json:"retry"`
	Delay int    `json:"delay"`
}
