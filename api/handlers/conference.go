package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/conference"
	"github.com/golang-cz/nilslice"
)

func GetSchedule(service conference.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		conferenceID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad conference ID",
			}
		}

		schedule, lastUpdated, err := service.GetSchedule(int32(conferenceID))
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: &dto.GetScheduleResponse{
				Schedule:    nilslice.Initialize(*schedule),
				LastUpdated: lastUpdated,
			},
		}
	})
}

func GetConferences(service conference.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		conferences, err := service.GetConferences()
		if err != nil {
			return err
		}

		var conferencesResponse []*dto.ConferenceResponse
		for _, c := range conferences {
			conference := &dto.ConferenceResponse{}
			conference.Scan(c)
			conferencesResponse = append(conferencesResponse, conference)
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: conferencesResponse,
		}
	})
}

func CreateConference(service conference.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.CreateConferenceRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		createdConference, err := service.CreateConference(request.URL)
		if err != nil {
			if errors.Is(err, conference.ErrScheduleFetch) {
				return &dto.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Could not fetch schedule from URL (is it a valid pentabarf XML file?)",
				}
			}
			return err
		}

		var response dto.ConferenceResponse
		response.Scan(*createdConference)
		return &dto.OkResponse{
			Code: http.StatusCreated,
			Data: response,
		}
	})
}

func DeleteConference(service conference.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.DeleteConferenceRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		err := service.DeleteConference(request.ID)
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: nil,
		}
	})
}
