package middleware

import (
	"errors"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func RequireAuthenticated(service user.Service, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		s, err := store.Get(c)
		if err != nil {
			return err
		}

		if s.Fresh() || len(s.Keys()) == 0 {
			return &dto.ErrorResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
			}
		}

		uid := s.Get("uid").(int32)

		fetchedUser, err := service.GetUserByID(uid)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				s.Destroy()
				return &dto.ErrorResponse{
					Code:    fiber.StatusUnauthorized,
					Message: "Invalid session",
				}
			}

			return err
		}

		c.Locals("uid", uid)
		c.Locals("username", fetchedUser.Username)

		return c.Next()
	}
}
