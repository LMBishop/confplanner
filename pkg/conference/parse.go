package conference

import (
	"encoding/xml"
	"fmt"
	"time"
)

type schedule struct {
	XMLName    xml.Name   `xml:"schedule"`
	Conference conference `xml:"conference"`
	Tracks     []track    `xml:"tracks>track"`
	Days       []day      `xml:"day"`
}

type conference struct {
	Title            string `xml:"title"`
	Venue            string `xml:"venue"`
	City             string `xml:"city"`
	Start            string `xml:"start"`
	End              string `xml:"end"`
	Days             int    `xml:"days"`
	DayChange        string `xml:"day_change"`
	TimeslotDuration string `xml:"timeslot_duration"`
	BaseURL          string `xml:"base_url"`
	TimeZoneName     string `xml:"time_zone_name"`
}

type track struct {
	Name string `xml:",chardata"`
}

type day struct {
	Date  string `xml:"date,attr"`
	Start string `xml:"start,attr"`
	End   string `xml:"end,attr"`
	Rooms []room `xml:"room"`
}

type room struct {
	Name   string  `xml:"name,attr"`
	Events []event `xml:"event"`
}

type event struct {
	ID          int32        `xml:"id,attr"`
	GUID        string       `xml:"guid,attr"`
	Date        string       `xml:"date"`
	Start       string       `xml:"start"`
	Duration    string       `xml:"duration"`
	Room        string       `xml:"room"`
	URL         string       `xml:"url"`
	Track       string       `xml:"track"`
	Type        string       `xml:"type"`
	Title       string       `xml:"title"`
	Abstract    string       `xml:"abstract"`
	Persons     []person     `xml:"persons>person"`
	Attachments []attachment `xml:"attachments>attachment"`
	Links       []link       `xml:"links>link"`
}

type person struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type attachment struct {
	Type string `xml:"id,attr"`
	Href string `xml:"href,attr"`
	Name string `xml:",chardata"`
}

type link struct {
	Href string `xml:"href,attr"`
	Name string `xml:",chardata"`
}

func (dst *Schedule) Scan(src schedule) error {
	dst.Conference.Scan(src.Conference)

	dst.Tracks = make([]Track, len(src.Tracks))
	for i := range src.Tracks {
		dst.Tracks[i].Scan(src.Tracks[i])
	}
	dst.Days = make([]Day, len(src.Days))
	for i := range src.Days {
		if err := dst.Days[i].Scan(src.Days[i]); err != nil {
			return fmt.Errorf("failed to scan day: %w", err)
		}
	}
	return nil
}

func (dst *Conference) Scan(src conference) {
	dst.Title = src.Title
	dst.Venue = src.Venue
	dst.City = src.City
	dst.Start = src.Start
	dst.End = src.End
	dst.Days = src.Days
	dst.DayChange = src.DayChange
	dst.TimeslotDuration = src.TimeslotDuration
	dst.BaseURL = src.BaseURL
	dst.TimeZoneName = src.TimeZoneName
}

func (dst *Track) Scan(src track) {
	dst.Name = src.Name
}

func (dst *Day) Scan(src day) error {
	dst.Date = src.Date

	start, err := time.Parse(time.RFC3339, src.Start)
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}
	end, err := time.Parse(time.RFC3339, src.End)
	if err != nil {
		return fmt.Errorf("failed to parse end time: %w", err)
	}

	dst.Start = start
	dst.End = end

	dst.Rooms = make([]Room, len(src.Rooms))
	for i := range src.Rooms {
		dst.Rooms[i].Scan(src.Rooms[i])
	}
	return nil
}

func (dst *Room) Scan(src room) {
	dst.Name = src.Name

	dst.Events = make([]Event, len(src.Events))
	for i := range src.Events {
		dst.Events[i].Scan(src.Events[i])
	}
}

func (dst *Event) Scan(src event) error {
	dst.ID = src.ID
	dst.GUID = src.GUID
	dst.Date = src.Date

	duration, err := parseDuration(src.Duration)
	if err != nil {
		return err
	}
	start, err := time.Parse(time.RFC3339, src.Date)
	if err != nil {
		start = time.Unix(0, 0)
	}
	dst.Start = start
	dst.End = start.Add(time.Minute * time.Duration(duration))

	dst.Room = src.Room
	dst.URL = src.URL
	dst.Track = src.Track
	dst.Type = src.Type
	dst.Title = src.Title
	dst.Abstract = src.Abstract

	dst.Persons = make([]Person, len(src.Persons))
	for i := range src.Persons {
		dst.Persons[i].Scan(src.Persons[i])
	}

	dst.Attachments = make([]Attachment, len(src.Attachments))
	for i := range src.Attachments {
		dst.Attachments[i].Scan(src.Attachments[i])
	}

	dst.Links = make([]Link, len(src.Links))
	for i := range src.Links {
		dst.Links[i].Scan(src.Links[i])
	}

	return nil
}

func (dst *Person) Scan(src person) {
	dst.ID = src.ID
	dst.Name = src.Name
}

func (dst *Attachment) Scan(src attachment) {
	dst.Type = src.Type
	dst.Href = src.Href
	dst.Name = src.Name
}

func (dst *Link) Scan(src link) {
	dst.Href = src.Href
	dst.Name = src.Name
}

func parseDuration(duration string) (int32, error) {
	d, err := time.Parse("15:04", duration)
	if err != nil {
		return 0, err
	}
	return int32(d.Minute() + d.Hour()*60), nil
}
