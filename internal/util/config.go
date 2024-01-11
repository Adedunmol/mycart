package util

import (
	"log"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
}

func LoadConfig(path string) (*Config, error) {

	viper.AddConfigPath(path)
	// viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
