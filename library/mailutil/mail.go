package mailutil

func Send(mail Mail, platform MailPlatform) error {
	switch platform {
	case MailPlatformAliyun:
		return SendByAliyun(mail)
	default:
		return SendByAliyun(mail)
	}
}
