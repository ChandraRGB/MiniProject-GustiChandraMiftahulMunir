package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/helper"
)

// JWTMiddleware validates JWT from the `token` header and injects
// user information into the request context.
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Unauthorized",
				"errors":  []string{"missing token"},
				"data":    nil,
			})
		}

		claims, err := helper.ParseJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Unauthorized",
				"errors":  []string{"invalid token"},
				"data":    nil,
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}
