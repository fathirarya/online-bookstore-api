package config

import (
	"time"

	"github.com/spf13/viper"
)

type JWTConfig struct {
	SecretKey      string
	Issuer         string
	ExpireDuration time.Duration
}

// LoadJWTConfig loads JWT config from environment variables using Viper
func LoadJWTConfig() *JWTConfig {
	// Set default values if needed
	viper.SetDefault("JWT_SECRET_KEY", "your_default_secret")
	viper.SetDefault("JWT_ISSUER", "your_app_name")
	viper.SetDefault("JWT_EXPIRE_DURATION", "60") // in minutes

	secret := viper.GetString("JWT_SECRET_KEY")
	issuer := viper.GetString("JWT_ISSUER")
	expireMinutes := viper.GetInt("JWT_EXPIRE_DURATION")

	return &JWTConfig{
		SecretKey:      secret,
		Issuer:         issuer,
		ExpireDuration: time.Duration(expireMinutes) * time.Minute,
	}
}
