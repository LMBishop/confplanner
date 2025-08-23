package conference

import "time"

type Schedule struct {
	Conference Conference `json:"conference"`
	Tracks     []Track    `json:"tracks"`
	Days       []Day      `json:"days"`
}

type Conference struct {
	Title            string `json:"title"`
	Venue            string `json:"venue"`
	City             string `json:"city"`
	Start            string `json:"start"`
	End              string `json:"end"`
	Days             int    `json:"days"`
	DayChange        string `json:"dayChange"`
	TimeslotDuration string `json:"timeslotDuration"`
	BaseURL          string `json:"baseUrl"`
	TimeZoneName     string `json:"timeZoneName"`
}

type Track struct {
	Name string `json:"name"`
}

type Day struct {
	Date  string    `json:"date"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Rooms []Room    `json:"rooms"`
}

type Room struct {
	Name   string  `json:"name"`
	Events []Event `json:"events"`
}

type Event struct {
	ID          int32        `json:"id"`
	GUID        string       `json:"guid"`
	Date        string       `json:"date"`
	Start       time.Time    `json:"start"`
	End         time.Time    `json:"end"`
	Duration    int32        `json:"duration"`
	Room        string       `json:"room"`
	URL         string       `json:"url"`
	Track       string       `json:"track"`
	Type        string       `json:"type"`
	Title       string       `json:"title"`
	Abstract    string       `json:"abstract"`
	Persons     []Person     `json:"persons"`
	Attachments []Attachment `json:"attachments"`
	Links       []Link       `json:"links"`
}

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Attachment struct {
	Type string `json:"string"`
	Href string `json:"href"`
	Name string `json:"name"`
}

type Link struct {
	Href string `json:"href"`
	Name string `json:"name"`
}
