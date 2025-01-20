package handlers

import (
	"crypto/subtle"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/ical"
	"github.com/gofiber/fiber/v2"
)

func GetIcal(icalService ical.Service, calendarService calendar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Query("name")
		key := c.Query("key")

		if name == "" || key == "" {
			return &dto.ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Both name and key must be specified",
			}
		}

		calendar, err := calendarService.GetCalendarByName(name)
		if err != nil {
			return err
		}

		if subtle.ConstantTimeCompare([]byte(key), []byte(calendar.Key)) != 1 {
			return &dto.ErrorResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid key",
			}
		}

		ical, err := icalService.GenerateIcalForCalendar(*calendar)
		if err != nil {
			return err
		}

		return c.SendString(ical)
	}
}
