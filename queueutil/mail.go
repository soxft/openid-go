package queueutil

import (
	"encoding/json"
	"log"
	"openid/library/mail"
)

// Mail
// @description: 邮件发送相关
func Mail(msg string) {
	var mailMsg mail.Mail
	if err := json.Unmarshal([]byte(msg), &mailMsg); err != nil {
		log.Panic(err)
	}
	log.Printf("send mail to %s", mailMsg.ToAddress)
	if err := mail.Send(mailMsg); err != nil {
		log.Panic(err)
	}
}
