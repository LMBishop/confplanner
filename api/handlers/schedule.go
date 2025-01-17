package handlers

import (
	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-cz/nilslice"
)

func GetSchedule(service schedule.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schedule, lastUpdated, err := service.GetSchedule()
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusOK,
			Data: &dto.GetScheduleResponse{
				Schedule:    nilslice.Initialize(*schedule),
				LastUpdated: *lastUpdated,
			},
		}
	}
}
