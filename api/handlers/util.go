package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func readBody(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: fmt.Errorf("Invalid request (%w)", err).Error(),
		}
	}

	if err := validate.Struct(request); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}
