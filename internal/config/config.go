package config

import "github.com/spf13/viper"

type Config struct {
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	DBConn       string `mapstructure:"DB_CONN"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
}

func LoadConfig(path string) (cfg *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&cfg)
	return
}
