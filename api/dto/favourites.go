package dto

import "github.com/LMBishop/confplanner/pkg/database/sqlc"

type CreateFavouritesRequest struct {
	ConferenceID int32   `json:"conferenceID" validate:"required"`
	GUID         *string `json:"eventGuid"`
	EventID      *int32  `json:"eventId"`
}

type CreateFavouritesResponse struct {
	ID int32 `json:"id"`
}

type GetFavouritesResponse struct {
	ID      int32   `json:"id"`
	GUID    *string `json:"eventGuid,omitempty"`
	EventID *int32  `json:"eventId,omitempty"`
}

func (dst *GetFavouritesResponse) Scan(src sqlc.Favourite) {
	dst.ID = src.ID
	if src.EventGuid.Valid {
		strGuid := src.EventGuid.String()
		dst.GUID = &strGuid
	}
	if src.EventID.Valid {
		dst.EventID = &src.EventID.Int32
	}
}

type DeleteFavouritesRequest struct {
	ConferenceID int32   `json:"conferenceID" validate:"required"`
	GUID         *string `json:"eventGuid"`
	EventID      *int32  `json:"eventId"`
}
