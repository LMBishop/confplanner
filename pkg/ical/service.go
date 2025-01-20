package ical

import (
	"errors"
	"strings"
	"time"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/microcosm-cc/bluemonday"
)

type Service interface {
	GenerateIcalForCalendar(calendar sqlc.Calendar) (string, error)
}

var (
	ErrImproperType = errors.New("improper type")
	ErrNotFound     = errors.New("not found")
)

type service struct {
	favouritesService favourites.Service
	scheduleService   schedule.Service
}

func NewService(
	favouritesService favourites.Service,
	scheduleService schedule.Service,
) Service {
	return &service{
		favouritesService: favouritesService,
		scheduleService:   scheduleService,
	}
}

func (s *service) GenerateIcalForCalendar(calendar sqlc.Calendar) (string, error) {
	favourites, err := s.favouritesService.GetFavouritesForUser(calendar.UserID)
	if err != nil {
		return "", err
	}

	events := make([]schedule.Event, 0)
	for _, favourite := range *favourites {
		event := s.scheduleService.GetEventByID(favourite.EventID.Int32)
		events = append(events, *event)
	}

	now := time.Now()

	// https://www.rfc-editor.org/rfc/rfc5545.html

	ret := "BEGIN:VCALENDAR\n"
	ret += "VERSION:2.0\n"
	ret += "METHOD:PUBLISH\n"
	ret += "X-WR-CALNAME:confplanner calendar\n"
	for _, event := range events {
		utcStart := event.Start.UTC()
		utcEnd := event.End.UTC()

		ret += "BEGIN:VEVENT\n"
		ret += "SUMMARY:" + event.Title + "\n"
		ret += "DTSTART:" + utcStart.Format("20060102T150405Z") + "\n"
		ret += "DTEND:" + utcEnd.Format("20060102T150405Z") + "\n"
		ret += "LOCATION:" + event.Room + "\n"
		ret += "DESCRIPTION;ENCODING=QUOTED-PRINTABLE:" + bluemonday.StrictPolicy().Sanitize(strings.Replace(event.Abstract, "\n", "\\n", -1)) + "\\n\\nconfplanner: last synchronised: " + now.Format(time.RFC1123) + "\n"
		ret += "END:VEVENT\n"
	}
	ret += "END:VCALENDAR\n"

	return ret, nil
}
