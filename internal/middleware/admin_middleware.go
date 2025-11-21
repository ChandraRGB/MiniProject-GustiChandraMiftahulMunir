package middleware

import "github.com/gofiber/fiber/v2"

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("is_admin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  false,
				"message": "Forbidden",
				"errors":  []string{"admin only"},
				"data":    nil,
			})
		}

		return c.Next()
	}
}
