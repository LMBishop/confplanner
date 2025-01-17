package handlers

import (
	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateFavourite(service favourites.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.CreateFavouritesRequest
		if err := readBody(c, &request); err != nil {
			return err
		}

		if request.GUID == nil && request.ID == nil {
			return &dto.ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "One of event GUID or event ID must be specified",
			}
		}

		uid := c.Locals("uid").(int32)
		var uuid pgtype.UUID
		if request.GUID != nil {
			if err := uuid.Scan(*request.GUID); err != nil {
				return &dto.ErrorResponse{
					Code:    fiber.StatusBadRequest,
					Message: "Bad event GUID",
				}
			}
		}

		createdFavourite, err := service.CreateFavouriteForUser(uid, uuid, request.ID)
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusCreated,
			Data: &dto.CreateFavouritesResponse{
				ID: createdFavourite.ID,
			},
		}
	}
}

func GetFavourites(service favourites.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid := c.Locals("uid").(int32)

		favourites, err := service.GetFavouritesForUser(uid)
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
			Code: fiber.StatusOK,
			Data: favouritesResponse,
		}
	}
}

func DeleteFavourite(service favourites.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.DeleteFavouritesRequest
		if err := readBody(c, &request); err != nil {
			return err
		}

		if request.GUID == nil && request.ID == nil {
			return &dto.ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "One of event GUID or event ID must be specified",
			}
		}

		uid := c.Locals("uid").(int32)
		var err error
		var uuid pgtype.UUID
		if err := uuid.Scan(*request.GUID); err != nil {
			return &dto.ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Bad event GUID",
			}
		}

		err = service.DeleteFavouriteForUserByEventDetails(uid, uuid, request.ID)
		if err != nil {
			if err == favourites.ErrNotFound {
				return &dto.ErrorResponse{
					Code:    fiber.StatusNotFound,
					Message: "Favourite not found",
				}
			}
			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusOK,
		}
	}
}
