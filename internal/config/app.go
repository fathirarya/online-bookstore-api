package config

import (
	"log"

	"github.com/fathirarya/online-bookstore-api/db/migrations"
	"github.com/fathirarya/online-bookstore-api/internal/auth"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/handler"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/middleware"
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
	userRepository := repository.NewUserRepository(config.DB, config.Log)
	categoryRepository := repository.NewCategoryRepository(config.DB, config.Log)
	bookRepository := repository.NewBookRepository(config.DB, config.Log)
	orderRepository := repository.NewOrderRepository(config.DB, config.Log)

	// setup usecases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, config.Validate, categoryRepository)
	bookUseCase := usecase.NewBookUseCase(config.DB, config.Log, config.Validate, bookRepository, categoryRepository)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validate, orderRepository, bookRepository)

	// setup JWT config & service
	jwtConfig := LoadJWTConfig()
	jwtService := auth.NewJWTService(jwtConfig)

	// setup handlers
	userHandler := handler.NewUserHandler(userUseCase, config.Log, jwtService)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase, config.Log)
	bookHandler := handler.NewBookHandler(bookUseCase, config.Log)
	orderHandler := handler.NewOrderHandler(orderUseCase, config.Log)

	routeConfig := routes.RouteConfig{
		App:            config.App,
		User:           userHandler,
		AuthMiddleware: middleware.JWTProtected(jwtService),
		Category:       categoryHandler,
		Book:           bookHandler,
		Order:          orderHandler,
	}
	routeConfig.Setup()
}
