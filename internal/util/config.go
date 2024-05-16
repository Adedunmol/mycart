package util

import (
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/spf13/viper"
)

var EnvConfig Config

type Config struct {
	DatabaseUrl     string `mapstructure:"DATABASE_URL"`
	TestDatabaseUrl string `mapstructure:"TEST_DATABASE_URL"`
	Environment     string `mapstructure:"ENVIRONMENT"`
	SecretKey       string `mapstructure:"SECRET_KEY"`
}

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error.Println(err)
		logger.Error.Fatal("Could not load env file")
		return err
	}

	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		logger.Error.Fatal("Could not unmarshal env file")
		return err
	}

	return err
}
