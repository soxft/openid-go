package dbutil

type Account struct {
	ID       int
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
	ID         int
	UserId     int
	AppId      string
	AppName    string
	AppSecret  string
	AppGateway string
	CreateAt   int64 `gorm:"autoCreateTime"`
}

type OpenID struct {
	ID       int
	UserId   int
	AppId    string
	OpenId   string
	CreateAt int64 `gorm:"autoCreateTime"`
}

type UniqueId struct {
	ID        int
	UserId    int
	DevUserId int
	UniqueId  string
	CreateAt  int64 `gorm:"autoCreateTime"`
}

func (UniqueId) TableName() string {
	return "unique_id"
}

func (OpenID) TableName() string {
	return "open_id"
}
