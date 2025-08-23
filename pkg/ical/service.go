package ical

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/LMBishop/confplanner/pkg/conference"
	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/LMBishop/confplanner/pkg/favourites"
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
	conferenceService conference.Service
}

func NewService(
	favouritesService favourites.Service,
	conferenceService conference.Service,
) Service {
	return &service{
		favouritesService: favouritesService,
		conferenceService: conferenceService,
	}
}

func (s *service) GenerateIcalForCalendar(calendar sqlc.Calendar) (string, error) {
	favourites, err := s.favouritesService.GetAllFavouritesForUser(calendar.UserID)
	if err != nil {
		return "", err
	}

	events := make([]conference.Event, 0)
	for _, favourite := range *favourites {
		event, err := s.conferenceService.GetEventByID(favourite.ConferenceID, favourite.EventID.Int32)
		if err != nil {
			continue
		}
		events = append(events, *event)
	}

	now := time.Now()
	counter := 0

	// https://www.rfc-editor.org/rfc/rfc5545.html

	ret := "BEGIN:VCALENDAR\r\n"
	ret += "PRODID:-//LMBishop//confplanner//EN\r\n"
	ret += "VERSION:2.0\r\n"
	ret += "METHOD:PUBLISH\r\n"
	ret += "X-WR-CALNAME:confplanner calendar\r\n"
	for _, event := range events {
		utcStart := event.Start.UTC()
		utcEnd := event.End.UTC()

		ret += "BEGIN:VEVENT\r\n"
		ret += "SUMMARY:" + event.Title + "\r\n"
		ret += "UID:" + now.Format("20060102T150405Z") + "-" + strconv.Itoa(counter) + "\r\n"
		ret += "DTSTAMP:" + now.Format("20060102T150405Z") + "\r\n"
		ret += "DTSTART:" + utcStart.Format("20060102T150405Z") + "\r\n"
		ret += "DTEND:" + utcEnd.Format("20060102T150405Z") + "\r\n"
		ret += "LOCATION:" + event.Room + "\r\n"
		ret += "DESCRIPTION;ENCODING=QUOTED-PRINTABLE:" + bluemonday.StrictPolicy().Sanitize(strings.Replace(event.Abstract, "\n", "\\n\\n", -1)) + "\\n\\nconfplanner: last synchronised: " + now.Format(time.RFC1123) + "\r\n"

		ret += "BEGIN:VALARM\r\n"
		ret += "TRIGGER:-PT10M\r\n"
		ret += "ACTION:AUDIO\r\n"
		ret += "END:VALARM\r\n"

		ret += "END:VEVENT\r\n"

		counter++
	}
	ret += "END:VCALENDAR\r\n"

	return ret, nil
}
