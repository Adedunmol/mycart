package util

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfig Config

type Config struct {
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
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

	// EnvConfig = Config{DatabaseUrl: "hey"}

	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		log.Fatal(err)
		return EnvConfig, err
	}

	return EnvConfig, nil
}
