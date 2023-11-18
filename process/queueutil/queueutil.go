package queueutil

import (
	"context"
	"github.com/soxft/openid-go/library/mq"
	"github.com/soxft/openid-go/process/redisutil"
	"log"
)

var Q mq.MessageQueue

// Init
// @desc golang消息队列
func Init() {
	log.Printf("[INFO] Queue initailizing...")

	// do nothing
	Q = mq.New(context.Background(), redisutil.RDB, 3)

	Q.Subscribe("mail", 2, Mail)

	log.Printf("[INFO] Queue initailize success")
}
