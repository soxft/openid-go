package mailutil

type Mail struct {
	Subject   string
	Content   string
	ToAddress string
	Typ       string // 邮件类型
}

type MailPlatform string

const (
	MailPlatformAliyun MailPlatform = "aliyun"
)
