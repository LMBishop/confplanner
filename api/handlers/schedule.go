package handlers

import (
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/golang-cz/nilslice"
)

func GetSchedule(service schedule.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		schedule, lastUpdated, err := service.GetSchedule()
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: &dto.GetScheduleResponse{
				Schedule:    nilslice.Initialize(*schedule),
				LastUpdated: lastUpdated,
			},
		}
	})
}
