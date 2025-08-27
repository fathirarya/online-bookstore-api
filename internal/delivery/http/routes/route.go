package routes

import (
	"github.com/fathirarya/online-bookstore-api/internal/delivery/http/handler"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	User           *handler.UserHandler
	AuthMiddleware fiber.Handler
	Category       *handler.CategoryHandler
	Book           *handler.BookHandler
	Order          *handler.OrderHandler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoutes()

}

func (c *RouteConfig) SetupGuestRoutes() {
	// Login & Register
	apiV1 := c.App.Group("/api")
	apiV1.Post("/register", c.User.Register)
	apiV1.Post("/login", c.User.Login)

	apiV1.Use(c.AuthMiddleware)
	// Categories
	apiV1.Post("/categories", c.Category.Create)
	apiV1.Get("/categories", c.Category.List)
	apiV1.Put("/categories/:id", c.Category.Update)
	apiV1.Delete("/categories/:id", c.Category.Delete)

	// Books
	apiV1.Post("/books", c.Book.Create)
	apiV1.Get("/books", c.Book.List)
	apiV1.Get("/books/:id", c.Book.GetByID)
	apiV1.Put("/books/:id", c.Book.Update)
	apiV1.Delete("/books/:id", c.Book.Delete)

	// Orders
	apiV1.Post("/orders", c.Order.Create)
	apiV1.Post("/orders/:id/pay", c.Order.Pay)
	apiV1.Get("/orders", c.Order.List)

	// Statistics
	apiV1.Get("/books/stats/total", c.Book.GetTotalBooks)
	apiV1.Get("/books/stats/price", c.Book.GetBookPriceStats)
}
