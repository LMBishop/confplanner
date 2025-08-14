package handlers

import (
	"errors"
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

func Register(service user.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.RegisterRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		createdUser, err := service.CreateUser(request.Username, request.Password)
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

func Login(service user.Service, store session.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var request dto.LoginRequest
		if err := dto.ReadDto(r, &request); err != nil {
			return err
		}

		user, err := service.Authenticate(request.Username, request.Password)
		if err != nil {
			return err
		}

		if user == nil {
			return &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Username and password combination not found",
			}
		}

		session, err := store.Create(user.ID, user.Username, r.RemoteAddr, r.UserAgent())
		if err != nil {
			return err
		}

		cookie := &http.Cookie{
			Name:  "confplanner_session",
			Value: session.Token,
		}
		http.SetCookie(w, cookie)

		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: &dto.LoginResponse{
				ID:       user.ID,
				Username: user.Username,
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
