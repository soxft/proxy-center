package config

import (
	"github.com/spf13/viper"
	"log"
)

func Parse() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("[config] Error reading config file, %s", err)
	}

}
