package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"course-project/internal/service"
)

func AuthMiddleware(svc *service.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		claims, err := svc.ParseToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("login", claims.Login)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole := c.Locals("role").(string)
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "insufficient permissions",
		})
	}
}
