package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get session",
			})
		}

		userID := sess.Get("userID")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// You can make userID available to handlers
		c.Locals("userID", userID.(uuid.UUID))
		return c.Next()
	}
}
