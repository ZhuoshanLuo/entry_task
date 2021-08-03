package model

type Config struct {
	Db struct {
		Driver   string `yaml:"driver"`
		SqlUser  string `yaml:"sqlUser"`
		Passwd   string `yaml:"passwd"`
		Host     string `yaml:"host"`
		Database string `yaml:"database"`
	}
}

