package queueutil

import (
	"log"
	"openid/library/mq"
	"openid/redisutil"
)

var Q mq.MessageQueue

// Init
// @desc golang消息队列
func Init() {
	// do nothing
	Q = mq.New(redisutil.R, 3)

	Q.Subscribe("test", 2, func(msg string) {
		log.Println(msg)
	})
}
