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
	GetSchedule() (*Schedule, *time.Time, error)
}

type Schedule struct {
	XMLName    xml.Name   `json:"-" xml:"schedule"`
	Conference Conference `json:"conference" xml:"conference"`
	Tracks     []Track    `json:"tracks" xml:"tracks>track"`
	Days       []Day      `json:"days" xml:"day"`
}

type Conference struct {
	Title            string `json:"title" xml:"title"`
	Venue            string `json:"venue" xml:"venue"`
	City             string `json:"city" xml:"city"`
	Start            string `json:"start" xml:"start"`
	End              string `json:"end" xml:"end"`
	Days             int    `json:"days" xml:"days"`
	DayChange        string `json:"dayChange" xml:"day_change"`
	TimeslotDuration string `json:"timeslotDuration" xml:"timeslot_duration"`
	BaseURL          string `json:"baseUrl" xml:"base_url"`
	TimeZoneName     string `json:"timeZoneName" xml:"time_zone_name"`
}

type Track struct {
	Name string `json:"name" xml:",chardata"`
}

type Day struct {
	Date  string `json:"date" xml:"date,attr"`
	Start string `json:"start" xml:"start"`
	End   string `json:"end" xml:"end"`
	Rooms []Room `json:"rooms" xml:"room"`
}

type Room struct {
	Name   string  `json:"name" xml:"name,attr"`
	Events []Event `json:"events" xml:"event"`
}

type Event struct {
	ID          int          `json:"id" xml:"id,attr"`
	GUID        string       `json:"guid" xml:"guid,attr"`
	Date        string       `json:"date" xml:"date"`
	Start       string       `json:"start" xml:"start"`
	Duration    string       `json:"duration" xml:"duration"`
	Room        string       `json:"room" xml:"room"`
	URL         string       `json:"url" xml:"url"`
	Track       string       `json:"track" xml:"track"`
	Type        string       `json:"type" xml:"type"`
	Title       string       `json:"title" xml:"title"`
	Abstract    string       `json:"abstract" xml:"abstract"`
	Persons     []Person     `json:"persons" xml:"persons>person"`
	Attachments []Attachment `json:"attachments" xml:"attachments>attachment"`
	Links       []Link       `json:"links" xml:"links>link"`
}

type Person struct {
	ID   int    `json:"id" xml:"id,attr"`
	Name string `json:"name" xml:",chardata"`
}

type Attachment struct {
	Type string `json:"string" xml:"id,attr"`
	Href string `json:"href" xml:"href,attr"`
	Name string `json:"name" xml:",chardata"`
}

type Link struct {
	Href string `json:"href" xml:"href,attr"`
	Name string `json:"name" xml:",chardata"`
}

type service struct {
	schedule     *Schedule
	pentabarfUrl string
	lastUpdated  time.Time
	lock         sync.Mutex
}

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

func (s *service) GetSchedule() (*Schedule, *time.Time, error) {
	if s.hasScheduleExpired() {
		err := s.updateSchedule()
		if err != nil {
			return nil, nil, err
		}
	}

	return s.schedule, &s.lastUpdated, nil
}

func (s *service) hasScheduleExpired() bool {
	expire := s.lastUpdated.Add(15 * time.Minute)
	return time.Now().After(expire)
}

func (s *service) updateSchedule() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.hasScheduleExpired() {
		return nil
	}

	res, err := http.Get(s.pentabarfUrl)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(res.Body)

	var schedule Schedule

	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&schedule); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	s.schedule = &schedule
	s.lastUpdated = time.Now()
	return nil
}
