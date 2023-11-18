package config

type Config struct {
	ServerConfig    `yaml:"Server"`
	RedisConfig     `yaml:"Redis"`
	MysqlConfig     `yaml:"Mysql"`
	SmtpConfig      `yaml:"Smtp"`
	AliyunConfig    `yaml:"Aliyun"`
	JwtConfig       `yaml:"Jwt"`
	DeveloperConfig `yaml:"Developer"`
}
type ServerConfig struct {
	Addr     string `yaml:"Address"`
	Debug    bool   `yaml:"Debug"`
	Log      bool   `yaml:"Log"`
	Title    string `yaml:"Title"`
	Name     string `yaml:"ServerName"`
	FrontUrl string `yaml:"FrontUrl"`
}

type RedisConfig struct {
	Addr       string `yaml:"Address"`
	Pwd        string `yaml:"Password"`
	Db         int    `yaml:"Database"`
	Prefix     string `yaml:"Prefix"`
	MinIdle    int    `yaml:"MinIdle"`
	MaxIdle    int    `yaml:"MaxIdle"`
	MaxActive  int    `yaml:"MaxActive"`
	MaxRetries int    `yaml:"MaxRetries"`
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

type SmtpConfig struct {
	Host   string `yaml:"Host"`
	Port   int    `yaml:"Port"`
	Secure bool   `yaml:"Secure"`
	User   string `yaml:"Username"`
	Pwd    string `yaml:"Password"`
}

type AliyunConfig struct {
	Domain       string `yaml:"Domain"`
	Region       string `yaml:"Region"`
	Version      string `yaml:"Version"`
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
