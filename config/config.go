package config

import (
	"github.com/spf13/viper"
	"log"
)

var Conf *viper.Viper

func LoadConfig(fileName string) {
	c := viper.New()
	c.SetConfigFile(fileName)
	if err := c.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	Conf = c
}
