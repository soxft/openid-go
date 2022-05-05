package config

type Config struct {
	Server struct {
		Addr  string `yaml:"Address"`
		Debug bool   `yaml:"Debug"`
		Log   bool   `yaml:"Log"`
	} `yaml:"Server"`
	Redis struct {
		Addr   string `yaml:"Address"`
		Pwd    string `yaml:"Password"`
		Db     int    `yaml:"Database"`
		Prefix string `yaml:"Prefix"`
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
}
