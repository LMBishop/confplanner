package handlers

import (
	"errors"
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/auth"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

func Register(userService user.Service, authService auth.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.RegisterRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		basicAuthProvider := authService.GetAuthProvider("basic")
		if _, ok := basicAuthProvider.(*auth.BasicAuthProvider); !ok {
			return &dto.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "Registrations are only accepted via an identity provider",
			}
		}

		createdUser, err := userService.CreateUser(request.Username, request.Password)
		if err != nil {
			if errors.Is(err, user.ErrUserExists) {
				return &dto.ErrorResponse{
					Code:    http.StatusConflict,
					Message: "User with that username already exists",
				}
			} else if errors.Is(err, user.ErrNotAcceptingRegistrations) {
				return &dto.ErrorResponse{
					Code:    http.StatusForbidden,
					Message: "This service is not currently accepting registrations",
				}
			}

			return err
		}

		return &dto.OkResponse{
			Code: http.StatusCreated,
			Data: &dto.RegisterResponse{
				ID: createdUser.ID,
			},
		}
	})
}

func Logout(store session.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		session := r.Context().Value("session").(*session.UserSession)

		err := store.Destroy(session.SessionID)
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: http.StatusNoContent,
		}
	})
}
