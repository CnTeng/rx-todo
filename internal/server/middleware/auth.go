package middleware

import (
	"strings"

	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/registry" || c.Path() == "/token" {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is required"})
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header format must be Bearer {token}"})
		}

		token := authParts[1]

		userID, err := db.GetUserIDByToken(&token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		c.Locals("user_id", userID)

		return c.Next()
	}
}
