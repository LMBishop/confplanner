package handlers

import (
	"errors"
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/session"
)

func GetCalendar(calendarService calendar.Service, baseURL string) http.HandlerFunc {
	// TODO create config service
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		session := r.Context().Value("session").(*session.UserSession)

		cal, err := calendarService.GetCalendarForUser(session.UserID)
		if err != nil {
			if errors.Is(err, calendar.ErrCalendarNotFound) {
				return &dto.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Calendar not found",
				}
			}

			return err
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: &dto.GetCalendarResponse{
				ID:   cal.ID,
				Name: cal.Name,
				Key:  cal.Key,
				URL:  baseURL + "/api/calendar/ical?name=" + cal.Name + "&key=" + cal.Key,
			},
		}
	})
}

func CreateCalendar(calendarService calendar.Service, baseURL string) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		session := r.Context().Value("session").(*session.UserSession)

		cal, err := calendarService.CreateCalendarForUser(session.UserID)
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusCreated,
			Data: &dto.CreateCalendarResponse{
				ID:   cal.ID,
				Name: cal.Name,
				Key:  cal.Key,
				URL:  baseURL + "/calendar/ical?name=" + cal.Name + "&key=" + cal.Key,
			},
		}
	})
}

func DeleteCalendar(calendarService calendar.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		session := r.Context().Value("session").(*session.UserSession)

		err := calendarService.DeleteCalendarForUser(session.UserID)
		if err != nil {
			if errors.Is(err, calendar.ErrCalendarNotFound) {
				return &dto.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Calendar not found",
				}
			}

			return err
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
		}
	})
}
