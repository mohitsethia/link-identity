package config

type DatabaseConfig struct {
	Host     string `mapstructure:"db_host"`
	Port     string `mapstructure:"db_port"`
	Username string `mapstructure:"db_username"`
	Name     string `mapstructure:"db_name"`
	Password string `mapstructure:"db_pass"`
}
