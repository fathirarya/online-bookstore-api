package main

import (
	"fmt"
	"log"

	"github.com/fathirarya/technical-test-backend/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	db := config.NewDatabase(viperConfig)
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB: db,
		// App:    *app,
		Config: viperConfig,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
