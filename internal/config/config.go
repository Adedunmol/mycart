package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfig Config

type Config struct {
	DatabaseUrl     string `mapstructure:"DATABASE_URL"`
	TestDatabaseUrl string `mapstructure:"DATABASE_URL_TEST"`
	Environment     string `mapstructure:"ENVIRONMENT"`
	SecretKey       string `mapstructure:"SECRET_KEY"`
	EmailUsername   string `mapstructure:"EMAIL_USERNAME"`
	EmailSender     string `mapstructure:"EMAIL_SENDER"`
	EmailPassword   string `mapstructure:"EMAIL_PASSWORD"`
	Port            string `mapstructure:"PORT"`
	RedisAddress    string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig() (Config, error) {

	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return EnvConfig, err
	}

	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		log.Fatal(err)
		return EnvConfig, err
	}

	return EnvConfig, nil
}
