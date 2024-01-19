package util

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
}

func LoadConfig(path string) (Config, error) {
	// var envConfig Config

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

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
