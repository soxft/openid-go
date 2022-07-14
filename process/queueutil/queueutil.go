package queueutil

import (
	"github.com/soxft/openid/library/mq"
	"github.com/soxft/openid/process/redisutil"
)

var Q mq.MessageQueue

// Init
// @desc golang消息队列
func Init() {
	// do nothing
	Q = mq.New(redisutil.R, 3)

	Q.Subscribe("mail", 2, Mail)
}
