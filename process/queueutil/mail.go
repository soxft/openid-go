package queueutil

import (
	"encoding/json"
	"github.com/soxft/openid/library/mailutil"
	"log"
)

// Mail
// @description: 邮件发送相关
func Mail(msg string) {
	var mailMsg mailutil.Mail
	if err := json.Unmarshal([]byte(msg), &mailMsg); err != nil {
		log.Panic(err)
	}
	if mailMsg.ToAddress == "" {
		log.Printf("[ERROR] Mail(%s) 空收件人", mailMsg.Typ)
		return
	}
	log.Printf("[INFO] Mail(%s) %s", mailMsg.Typ, mailMsg.ToAddress)

	// get mail platform
	var platform mailutil.MailPlatform
	switch mailMsg.Typ {
	case "register":
		platform = mailutil.MailplatformAliyun
	default:
		platform = mailutil.MailplatformAliyun
	}
	// send mail
	if err := mailutil.Send(mailMsg, platform); err != nil {
		log.Panic(err)
	}
}
