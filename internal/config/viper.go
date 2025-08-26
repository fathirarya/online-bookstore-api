package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	// Read from .env file
	config.SetConfigFile(".env")
	config.SetConfigType("env")

	// Validate if the config file is found
	err := config.ReadInConfig()

	// Panic if the config file is not found
	if err != nil {
		panic(fmt.Errorf("Fatal error cinfig file: %w \n", err))
	}

	return config

}
