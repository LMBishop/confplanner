package conference

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	CreateConference(url string) (*sqlc.Conference, error)
	DeleteConference(id int32) error
	GetConferences() ([]sqlc.Conference, error)
	GetSchedule(id int32) (*Schedule, time.Time, error)
	GetEventByID(conferenceID, eventID int32) (*Event, error)
}

type loadedConference struct {
	pentabarfUrl string
	schedule     *Schedule
	eventsById   map[int32]Event
	lastUpdated  time.Time
	lock         sync.RWMutex
}

var (
	ErrConferenceNotFound = errors.New("conference not found")
	ErrScheduleFetch      = errors.New("could not fetch schedule")
)

type service struct {
	conferences map[int32]*loadedConference
	lock        sync.RWMutex
	pool        *pgxpool.Pool
}

// TODO: Create a service implementation that persists to DB
// and isn't in memory
func NewService(pool *pgxpool.Pool) (Service, error) {
	service := &service{
		pool:        pool,
		conferences: make(map[int32]*loadedConference),
	}

	queries := sqlc.New(pool)
	conferences, err := queries.GetConferences(context.Background())
	if err != nil {
		return nil, err
	}

	for _, conference := range conferences {
		c := &loadedConference{
			pentabarfUrl: conference.Url,
			lastUpdated:  time.Unix(0, 0),
		}
		service.conferences[conference.ID] = c
	}

	return service, nil
}

func (s *service) CreateConference(url string) (*sqlc.Conference, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	c := &loadedConference{
		pentabarfUrl: url,
		lastUpdated:  time.Unix(0, 0),
	}
	err := c.updateSchedule()
	if err != nil {
		return nil, errors.Join(ErrScheduleFetch, err)
	}

	queries := sqlc.New(s.pool)

	conference, err := queries.CreateConference(context.Background(), sqlc.CreateConferenceParams{
		Url:   url,
		Title: pgtype.Text{String: c.schedule.Conference.Title, Valid: true},
		Venue: pgtype.Text{String: c.schedule.Conference.Venue, Valid: true},
		City:  pgtype.Text{String: c.schedule.Conference.City, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("could not create conference: %w", err)
	}

	s.conferences[conference.ID] = c

	return &conference, nil
}

func (s *service) DeleteConference(id int32) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	queries := sqlc.New(s.pool)
	err := queries.DeleteConference(context.Background(), id)
	if err != nil {
		return fmt.Errorf("could not delete conference: %w", err)
	}

	delete(s.conferences, id)
	return nil
}

func (s *service) GetConferences() ([]sqlc.Conference, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	queries := sqlc.New(s.pool)
	return queries.GetConferences(context.Background())
}

func (s *service) GetSchedule(id int32) (*Schedule, time.Time, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	c, ok := s.conferences[id]
	if !ok {
		return nil, time.Time{}, ErrConferenceNotFound
	}

	if err := c.updateSchedule(); err != nil {
		return nil, time.Time{}, err
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	queries := sqlc.New(s.pool)
	if _, err := queries.UpdateConferenceDetails(context.Background(), sqlc.UpdateConferenceDetailsParams{
		ID:    id,
		Title: pgtype.Text{String: c.schedule.Conference.Title, Valid: true},
		Venue: pgtype.Text{String: c.schedule.Conference.Venue, Valid: true},
		City:  pgtype.Text{String: c.schedule.Conference.City, Valid: true},
	}); err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to update cached conference details: %w", err)
	}

	return c.schedule, c.lastUpdated, nil
}

func (s *service) GetEventByID(conferenceID, eventID int32) (*Event, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	c, ok := s.conferences[conferenceID]
	if !ok {
		return nil, ErrConferenceNotFound
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	event := c.eventsById[eventID]

	return &event, nil
}

func (c *loadedConference) hasScheduleExpired() bool {
	expire := c.lastUpdated.Add(15 * time.Minute)
	return time.Now().After(expire)
}

func (c *loadedConference) updateSchedule() error {
	if !c.hasScheduleExpired() {
		return nil
	}

	if !c.lock.TryLock() {
		// don't block if another goroutine is already fetching
		return nil
	}
	defer c.lock.Unlock()

	res, err := http.Get(c.pentabarfUrl)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(res.Body)

	var schedule schedule

	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&schedule); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	var newSchedule Schedule
	err = newSchedule.Scan(schedule)
	if err != nil {
		return fmt.Errorf("failed to scan schedule: %w", err)
	}

	c.schedule = &newSchedule
	c.lastUpdated = time.Now()

	c.eventsById = make(map[int32]Event)

	for _, day := range newSchedule.Days {
		for _, room := range day.Rooms {
			for _, event := range room.Events {
				c.eventsById[event.ID] = event
			}
		}
	}

	return nil
}
