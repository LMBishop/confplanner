package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/auth"
	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/LMBishop/confplanner/pkg/session"
)

func Login(authService auth.Service, store session.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		provider := authService.GetAuthProvider(r.PathValue("provider"))

		var user *sqlc.User
		var err error
		switch p := provider.(type) {
		case *auth.BasicAuthProvider:
			user, err = doBasicAuth(r, p)
		case *auth.OIDCAuthProvider:
			user, err = doOIDCAuthJourney(r, p)
		default:
			return &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Unknown auth provider",
			}
		}

		if err != nil {
			return err
		}

		// TODO X-Forwarded-For
		session, err := store.Create(user.ID, user.Username, r.RemoteAddr, r.UserAgent())
		if err != nil {
			return err
		}

		cookie := &http.Cookie{
			Name:  "confplanner_session",
			Value: session.Token,
			Path:  "/api",
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

func GetLoginOptions(authService auth.Service) http.HandlerFunc {
	return dto.WrapResponseFunc(func(w http.ResponseWriter, r *http.Request) error {
		var loginOptions []dto.LoginOption

		for _, identifier := range authService.GetAuthProviders() {
			provider := authService.GetAuthProvider(identifier)
			loginOptions = append(loginOptions, dto.LoginOption{
				Name:       provider.Name(),
				Identifier: identifier,
				Type:       provider.Type(),
			})
		}
		return &dto.OkResponse{
			Code: http.StatusOK,
			Data: &dto.LoginOptionsResponse{
				Options: loginOptions,
			},
		}
	})
}

func doBasicAuth(r *http.Request, p *auth.BasicAuthProvider) (*sqlc.User, error) {
	var request dto.LoginBasicRequest
	if err := dto.ReadDto(r, &request); err != nil {
		return nil, err
	}

	user, err := p.Authenticate(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Username and password combination not found",
		}
	}

	return user, nil
}

func doOIDCAuthJourney(r *http.Request, p *auth.OIDCAuthProvider) (*sqlc.User, error) {
	var request dto.LoginOAuthCallbackRequest
	if err := dto.ReadDto(r, &request); err != nil {
		url, err := p.StartJourney(r.RemoteAddr, r.UserAgent())
		if err != nil {
			return nil, &dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Could not start OAuth journey",
			}
		}

		return nil, &dto.OkResponse{
			Code: http.StatusTemporaryRedirect,
			Data: &dto.LoginOAuthOutboundResponse{
				URL: url,
			},
		}
	}

	user, err := p.CompleteJourney(r.Context(), request.Code, request.State, r.RemoteAddr, r.UserAgent())
	if err != nil {
		if errors.Is(err, auth.ErrNotAuthorised) {
			return nil, &dto.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "You are not authorised to use this service",
			}
		} else if errors.Is(err, auth.ErrInvalidState) {
			return nil, &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid state",
			}
		} else if errors.Is(err, auth.ErrStateVerificationFailed) {
			return nil, &dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "State verification failed",
			}
		} else if errors.Is(err, auth.ErrUserSyncFailed) {
			return nil, &dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "User sync failed",
			}
		}
		slog.Error("error completing oidc journey", "error", err, "ip", r.RemoteAddr)
		return nil, err
	}

	return user, nil
}
