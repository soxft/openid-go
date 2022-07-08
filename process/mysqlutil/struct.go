package mysqlutil

type Account struct {
	ID       uint
	Username string
	Password string
	Salt     string
	Email    string
	RegTime  int64
	RegIp    int64
	LastTime int64
	LastIp   int64
}

type App struct {
	ID         uint
	UserId     uint
	AppId      int64
	AppName    string
	AppSecret  string
	AppGateway string
	CreateAt   int64 `gorm:"autoCreateTime"`
}

type OpenID struct {
	ID       uint
	UserId   uint
	AppId    int64
	OpenId   string
	CreateAt int64 `gorm:"autoCreateTime"`
}

type UniqueId struct {
	ID        uint
	UserId    uint
	DevUserId uint
	UniqueId  string
	CreateAt  int64 `gorm:"autoCreateTime"`
}
