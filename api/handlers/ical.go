package handlers

import (
	"crypto/subtle"
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/ical"
)

func GetIcal(icalService ical.Service, calendarService calendar.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		key := r.URL.Query().Get("key")

		if name == "" || key == "" {
			dto.WriteDto(w, r, &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Both name and key must be specified",
			})
			return
		}

		calendar, err := calendarService.GetCalendarByName(name)
		if err != nil {
			dto.WriteDto(w, r, err)
			return
		}

		if subtle.ConstantTimeCompare([]byte(key), []byte(calendar.Key)) != 1 {
			dto.WriteDto(w, r, &dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid key",
			})
			return
		}

		ical, err := icalService.GenerateIcalForCalendar(*calendar)
		if err != nil {
			dto.WriteDto(w, r, err)
			return
		}

		w.Header().Add("Content-Type", "text/calendar")
		w.Write([]byte(ical))
	}
}
