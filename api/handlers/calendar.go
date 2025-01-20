package handlers

import (
	"errors"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/gofiber/fiber/v2"
)

func GetCalendar(calendarService calendar.Service, baseURL string) fiber.Handler {
	// TODO create config service
	return func(c *fiber.Ctx) error {
		uid := c.Locals("uid").(int32)

		cal, err := calendarService.GetCalendarForUser(uid)
		if err != nil {
			if errors.Is(err, calendar.ErrCalendarNotFound) {
				return &dto.ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Calendar not found",
				}
			}

			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusOK,
			Data: &dto.GetCalendarResponse{
				ID:   cal.ID,
				Name: cal.Name,
				Key:  cal.Key,
				URL:  baseURL + "/calendar/ical?name=" + cal.Name + "&key=" + cal.Key,
			},
		}
	}
}

func CreateCalendar(calendarService calendar.Service, baseURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid := c.Locals("uid").(int32)

		cal, err := calendarService.CreateCalendarForUser(uid)
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusCreated,
			Data: &dto.CreateCalendarResponse{
				ID:   cal.ID,
				Name: cal.Name,
				Key:  cal.Key,
				URL:  baseURL + "/calendar/ical?name=" + cal.Name + "&key=" + cal.Key,
			},
		}
	}
}

func DeleteCalendar(calendarService calendar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid := c.Locals("uid").(int32)

		err := calendarService.DeleteCalendarForUser(uid)
		if err != nil {
			if errors.Is(err, calendar.ErrCalendarNotFound) {
				return &dto.ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Calendar not found",
				}
			}

			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusOK,
		}
	}
}
