package config

import (
	"log"

	"github.com/fathirarya/online-bookstore-api/db/migrations"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/handler"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/routes"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// Run AutoMigrate
	migrations.Migrate(config.DB)
	log.Println("✅ Datbase migration completed")

	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)

	// setup usecases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)

	// setup handlers
	userHandler := handler.NewUserHandler(userUseCase, config.Log)
	routeConfig := routes.RouteConfig{
		App:        config.App,
		UseHandler: userHandler,
	}
	routeConfig.Setup()
}
