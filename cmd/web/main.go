package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fathirarya/online-bookstore-api/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("WEB_PORT")

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// signal handling di sini
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)
	<-sigTerm

	slog.Info("Shutting down gracefully...")
	_ = app.Shutdown()
}
