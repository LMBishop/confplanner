package handlers

import (
	"errors"
	"time"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Register(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.RegisterRequest
		if err := readBody(c, &request); err != nil {
			return err
		}

		createdUser, err := service.CreateUser(request.Username, request.Password)
		if err != nil {
			if errors.Is(err, user.ErrUserExists) {
				return &dto.ErrorResponse{
					Code:    fiber.StatusConflict,
					Message: "User with that username already exists",
				}
			} else if errors.Is(err, user.ErrNotAcceptingRegistrations) {
				return &dto.ErrorResponse{
					Code:    fiber.StatusForbidden,
					Message: "This service is not currently accepting registrations",
				}
			}

			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusCreated,
			Data: &dto.RegisterResponse{
				ID: createdUser.ID,
			},
		}
	}
}

func Login(service user.Service, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.LoginRequest
		if err := readBody(c, &request); err != nil {
			return err
		}

		user, err := service.Authenticate(request.Username, request.Password)
		if err != nil {
			return err
		}

		if user == nil {
			return &dto.ErrorResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Username and password combination not found",
			}
		}

		s, err := store.Get(c)
		if err != nil {
			return err
		}

		if s.Fresh() {
			uid := user.ID
			sid := s.ID()

			s.Set("uid", uid)
			s.Set("sid", sid)
			s.Set("ip", c.Context().RemoteIP().String())
			s.Set("login", time.Unix(time.Now().Unix(), 0).UTC().String())
			s.Set("ua", string(c.Request().Header.UserAgent()))

			err = s.Save()
			if err != nil {
				return err
			}
		}

		return &dto.OkResponse{
			Code: fiber.StatusOK,
			Data: &dto.LoginResponse{
				ID:       user.ID,
				Username: user.Username,
			},
		}
	}
}

func Logout(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		s, err := store.Get(c)
		if err != nil {
			return err
		}

		err = s.Destroy()
		if err != nil {
			return err
		}

		return &dto.OkResponse{
			Code: fiber.StatusNoContent,
		}
	}
}
