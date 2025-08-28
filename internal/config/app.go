package config

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/fathirarya/online-bookstore-api/db/migrations"
	"github.com/fathirarya/online-bookstore-api/internal/auth"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/handler"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/middleware"
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/routes"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
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
	log.Println("✅ Database migration completed")

	// setup repositories
	userRepository := repository.NewUserRepository(config.DB, config.Log)
	categoryRepository := repository.NewCategoryRepository(config.DB, config.Log)
	bookRepository := repository.NewBookRepository(config.DB, config.Log)
	orderRepository := repository.NewOrderRepository(config.DB, config.Log)

	// setup usecases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, userRepository)
	categoryUseCase := usecase.NewCategoryUseCase(config.DB, config.Log, categoryRepository)
	bookUseCase := usecase.NewBookUseCase(config.DB, config.Log, bookRepository, categoryRepository)
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, config.Validate, orderRepository, bookRepository)

	// setup JWT config & service
	jwtConfig := LoadJWTConfig()
	jwtService := auth.NewJWTService(jwtConfig)

	// setup handlers
	userHandler := handler.NewUserHandler(userUseCase, config.Log, jwtService, config.Validate)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase, config.Log, config.Validate)
	bookHandler := handler.NewBookHandler(bookUseCase, config.Log, config.Validate)
	orderHandler := handler.NewOrderHandler(orderUseCase, config.Log)

	// setup routes
	routeConfig := routes.RouteConfig{
		App:            config.App,
		User:           userHandler,
		AuthMiddleware: middleware.JWTProtected(jwtService),
		Category:       categoryHandler,
		Book:           bookHandler,
		Order:          orderHandler,
	}
	routeConfig.Setup()

	// setup cron job
	ctx := context.Background()
	scheduler := cron.New(cron.WithLocation(time.Local))
	orderCronjob := usecase.NewOrderCronJob(config.DB, config.Log, orderRepository)
	_, err := scheduler.AddFunc("*/2 * * * *", func() { orderCronjob.CheckingOrderPaymentStatus(ctx) })
	if err != nil {
		slog.Error("Failed to add cron job", "error", err.Error())
		os.Exit(1)
	}
	go scheduler.Start()
	slog.Info("Cron job started")
}
