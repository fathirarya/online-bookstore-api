package config

import (
	"log"

	"github.com/fathirarya/technical-test-backend/db/migrations"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      fiber.App
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// Run AutoMigrate
	migrations.Migrate(config.DB)
	log.Println("✅ Datbase migration completed")
	// routeConfig := route.RouteConfig{
	// 	App: config.App,
	// }
	// routeConfig.Setup()
}
