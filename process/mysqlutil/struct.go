package mysqlutil

type Account struct {
	ID       uint
	Username string
	Password string
	Salt     string
	Email    string
	RegTime  int64
	RegIp    string
	LastTime int64
	LastIp   string
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

func (UniqueId) TableName() string {
	return "unique_id"
}

func (OpenID) TableName() string {
	return "open_id"
}
