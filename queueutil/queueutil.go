package queueutil

import (
	"encoding/json"
	"log"
	"openid/library/mail"
	"openid/library/mq"
	"openid/redisutil"
)

var Q mq.MessageQueue

// Init
// @desc golang消息队列
func Init() {
	// do nothing
	Q = mq.New(redisutil.R, 3)

	Q.Subscribe("mail", 2, func(msg string) {
		var mailMsg mail.Mail
		if err := json.Unmarshal([]byte(msg), &mailMsg); err != nil {
			log.Panic(err)
		}
		if err := mail.Send(mailMsg); err != nil {
			log.Panic(err)
		}
	})
}
