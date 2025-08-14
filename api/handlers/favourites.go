package handlers

import (
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateFavourite(service favourites.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.CreateFavouritesRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		if request.GUID == nil && request.ID == nil {
			return &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "One of event GUID or event ID must be specified",
			}
		}

		session := r.Context().Value("session").(*session.UserSession)
		var uuid pgtype.UUID
		if request.GUID != nil {
			if err := uuid.Scan(*request.GUID); err != nil {
				return &dto.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Bad event GUID",
				}
			}
		}

		createdFavourite, err := service.CreateFavouriteForUser(session.UserID, uuid, request.ID)
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusCreated,
			Data: &dto.CreateFavouritesResponse{
				ID: createdFavourite.ID,
			},
		}
	})
}

func GetFavourites(service favourites.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		session := r.Context().Value("session").(*session.UserSession)

		favourites, err := service.GetFavouritesForUser(session.UserID)
		if err != nil {
			return err
		}

		favouritesResponse := make([]dto.GetFavouritesResponse, 0)
		for _, favourite := range *favourites {
			var favouriteResponse dto.GetFavouritesResponse
			favouriteResponse.Scan(favourite)

			favouritesResponse = append(favouritesResponse, favouriteResponse)
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: favouritesResponse,
		}
	})
}

func DeleteFavourite(service favourites.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.DeleteFavouritesRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		if request.GUID == nil && request.ID == nil {
			return &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "One of event GUID or event ID must be specified",
			}
		}

		session := r.Context().Value("session").(*session.UserSession)
		var err error
		var uuid pgtype.UUID
		if err := uuid.Scan(*request.GUID); err != nil {
			return &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Bad event GUID",
			}
		}

		err = service.DeleteFavouriteForUserByEventDetails(session.UserID, uuid, request.ID)
		if err != nil {
			if err == favourites.ErrNotFound {
				return &dto.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Favourite not found",
				}
			}
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusOK,
		}
	})
}
