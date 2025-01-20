package schedule

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Service interface {
	GetSchedule() (*Schedule, time.Time, error)
	GetEventByID(id int32) *Event
}

type service struct {
	pentabarfUrl string

	schedule    *Schedule
	eventsById  map[int32]Event
	lastUpdated time.Time
	accessLock  sync.RWMutex
	updateLock  sync.Mutex
}

// TODO: Create a service implementation that persists to DB
// and isn't in memory
func NewService(pentabarfUrl string) (Service, error) {
	service := &service{
		pentabarfUrl: pentabarfUrl,
		lastUpdated:  time.Unix(0, 0),
	}

	err := service.updateSchedule()
	if err != nil {
		return nil, fmt.Errorf("could not read schedule from '%s' (is it a valid pentabarf XML file?): %w", pentabarfUrl, err)
	}
	return service, nil
}

func (s *service) GetSchedule() (*Schedule, time.Time, error) {
	err := s.updateSchedule()
	if err != nil {
		return nil, time.Time{}, err
	}

	s.accessLock.RLock()
	defer s.accessLock.RUnlock()

	return s.schedule, s.lastUpdated, nil
}

func (s *service) GetEventByID(id int32) *Event {
	s.accessLock.RLock()
	defer s.accessLock.RUnlock()

	event := s.eventsById[id]

	return &event
}

func (s *service) hasScheduleExpired() bool {
	expire := s.lastUpdated.Add(15 * time.Minute)
	return time.Now().After(expire)
}

func (s *service) updateSchedule() error {
	if !s.hasScheduleExpired() {
		return nil
	}

	if !s.updateLock.TryLock() {
		// don't block if another goroutine is already fetching
		return nil
	}
	defer s.updateLock.Unlock()

	res, err := http.Get(s.pentabarfUrl)
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

	s.accessLock.Lock()
	defer s.accessLock.Unlock()

	s.schedule = &newSchedule
	s.lastUpdated = time.Now()

	s.eventsById = make(map[int32]Event)

	for _, day := range newSchedule.Days {
		for _, room := range day.Rooms {
			for _, event := range room.Events {
				s.eventsById[event.ID] = event
			}
		}
	}

	return nil
}
