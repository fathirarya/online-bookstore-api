package routes

import (
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/handler"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UseHandler     *handler.UserHandler
	AuthMiddleware fiber.Handler
	Category       *handler.CategoryHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoutes()
	c.SetupAuthRoutes()

}

func (c *RouteConfig) SetupGuestRoutes() {
	// Hanya untuk login & register
	c.App.Post("/api/register", c.UseHandler.Register)
	c.App.Post("/api/login", c.UseHandler.Login)
}

func (c *RouteConfig) SetupAuthRoutes() {
	c.App.Use(c.AuthMiddleware)

	// Categories
	c.App.Post("/api/categories", c.Category.Create)
	c.App.Get("/api/categories", c.Category.List)
	c.App.Put("/api/categories/:id", c.Category.Update)
	c.App.Delete("/api/categories/:id", c.Category.Delete)

	// // Books
	// c.App.Get("/api/books", c.BookHandler.List)
	// c.App.Get("/api/books/:id", c.BookHandler.Get)
	// c.App.Post("/api/books", c.BookHandler.Create)
	// c.App.Put("/api/books/:id", c.BookHandler.Update)
	// c.App.Delete("/api/books/:id", c.BookHandler.Delete)

	// // Orders
	// c.App.Post("/api/orders", c.OrderHandler.Create)
	// c.App.Post("/api/orders/:id/pay", c.OrderHandler.Pay)
	// c.App.Get("/api/orders", c.OrderHandler.List)

	// // Statistics
	// c.App.Get("/api/books/stats/total", c.StatsHandler.TotalBooks)
	// c.App.Get("/api/books/stats/price", c.StatsHandler.PriceStats)
}
