package config

type RiskBird struct {
	DB  RiskBirdDB  `mapstructure:"db" json:"db" yaml:"db"`
	API RiskBirdAPI `mapstructure:"api" json:"api" yaml:"api"`
}

type RiskBirdDB struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
}

type RiskBirdAPI struct {
	BaseUrl string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
}
