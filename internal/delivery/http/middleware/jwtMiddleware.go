package middleware

import (
	"strings"

	"github.com/fathirarya/online-bookstore-api/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing or invalid authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid or expired token",
			})
		}

		// Set user info ke context jika perlu
		c.Locals("user_id", claims.UserID)

		return c.Next()
	}
}
