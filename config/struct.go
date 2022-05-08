package config

type Config struct {
	Server struct {
		Addr  string `yaml:"Address"`
		Debug bool   `yaml:"Debug"`
		Log   bool   `yaml:"Log"`
		Title string `yaml:"Title"`
		Name  string `yaml:"ServerName"`
	} `yaml:"Server"`
	Redis struct {
		Addr      string `yaml:"Address"`
		Pwd       string `yaml:"Password"`
		Db        int    `yaml:"Database"`
		Prefix    string `yaml:"Prefix"`
		MaxIdle   int    `yaml:"MaxIdle"`
		MaxActive int    `yaml:"MaxActive"`
	} `yaml:"Redis"`
	Mysql struct {
		Addr        string `yaml:"Address"`
		User        string `yaml:"Username"`
		Pwd         string `yaml:"Password"`
		Db          string `yaml:"Database"`
		Charset     string `yaml:"Charset"`
		MaxOpen     int    `yaml:"MaxOpen"`
		MaxIdle     int    `yaml:"MaxIdle"`
		MaxLifetime int    `yaml:"MaxLifetime"`
	} `yaml:"Mysql"`
	Aliyun struct {
		AccessKey    string `yaml:"AccessKey"`
		AccessSecret string `yaml:"AccessSecret"`
		Email        string `yaml:"Email"`
	} `yaml:"Aliyun"`
	Jwt struct {
		Secret string `yaml:"Secret"`
	} `yaml:"Jwt"`
	Developer struct {
		AppLimit int `yaml:"AppLimit"`
	} `yaml:"Developer"`
}
