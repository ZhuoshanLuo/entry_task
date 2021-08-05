package model

type Config struct {
	Database DatabaseConf
}

type DatabaseConf struct {
	Driver   string `yaml:"driver"`
	SqlUser  string `yaml:"sqlUser"`
	Passwd   string `yaml:"passwd"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}
