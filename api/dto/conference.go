package dto

import (
	"time"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
)

type ConferenceResponse struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Venue string `json:"venue"`
	City  string `json:"city"`
}

func (dst *ConferenceResponse) Scan(src sqlc.Conference) {
	dst.ID = src.ID
	dst.Title = src.Title.String
	dst.URL = src.Url
	dst.Venue = src.Venue.String
	dst.City = src.City.String
}

type GetScheduleResponse struct {
	Schedule    interface{} `json:"schedule"`
	LastUpdated time.Time   `json:"lastUpdated"`
}

type CreateConferenceRequest struct {
	URL string `json:"url" validate:"required"`
}

type DeleteConferenceRequest struct {
	ID int32 `json:"id"`
}
