package config

type DatabaseConfig struct {
	Type        string `mapstructure:"type"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	SSLmode     string `mapstructure:"sslmode"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
}
