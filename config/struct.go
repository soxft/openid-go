package config

type Config struct {
	ServerConfig    `yaml:"Server"`
	RedisConfig     `yaml:"Redis"`
	MysqlConfig     `yaml:"Mysql"`
	AliyunConfig    `yaml:"Aliyun"`
	JwtConfig       `yaml:"Jwt"`
	DeveloperConfig `yaml:"Developer"`
	GithubConfig    `yaml:"Github"`
}

type ServerConfig struct {
	Addr  string `yaml:"Address"`
	Debug bool   `yaml:"Debug"`
	Log   bool   `yaml:"Log"`
	Title string `yaml:"Title"`
	Name  string `yaml:"ServerName"`
	Url   string `yaml:"Url"`
}

type RedisConfig struct {
	Addr      string `yaml:"Address"`
	Pwd       string `yaml:"Password"`
	Db        int    `yaml:"Database"`
	Prefix    string `yaml:"Prefix"`
	MaxIdle   int    `yaml:"MaxIdle"`
	MaxActive int    `yaml:"MaxActive"`
}

type MysqlConfig struct {
	Addr        string `yaml:"Address"`
	User        string `yaml:"Username"`
	Pwd         string `yaml:"Password"`
	Db          string `yaml:"Database"`
	Charset     string `yaml:"Charset"`
	MaxOpen     int    `yaml:"MaxOpen"`
	MaxIdle     int    `yaml:"MaxIdle"`
	MaxLifetime int    `yaml:"MaxLifetime"`
}

type AliyunConfig struct {
	AccessKey    string `yaml:"AccessKey"`
	AccessSecret string `yaml:"AccessSecret"`
	Email        string `yaml:"Email"`
}

type JwtConfig struct {
	Secret string `yaml:"Secret"`
}

type DeveloperConfig struct {
	AppLimit int `yaml:"AppLimit"`
}

type GithubConfig struct {
	ClientID     string `yaml:"ClientID"`
	ClientSecret string `yaml:"ClientSecret"`
}
