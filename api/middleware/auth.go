package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

func MustAuthenticate(service user.Service, store session.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token, err := extractBearerToken(authHeader)
			if err != nil {
				dto.WriteDto(w, r, &dto.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized",
				})
				return
			}

			s := store.GetByToken(token)
			if s == nil {
				dto.WriteDto(w, r, &dto.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized",
				})
				return
			}

			u, err := service.GetUserByID(s.UserID)
			if err != nil {
				if errors.Is(err, user.ErrUserNotFound) {
					store.Destroy(s.SessionID)
					dto.WriteDto(w, r, &dto.ErrorResponse{
						Code:    http.StatusForbidden,
						Message: "Invalid session",
					})
					return
				}

				return
			}

			s.Username = u.Username
			s.Admin = u.Admin

			ctx := context.WithValue(r.Context(), "session", s)

			next(w, r.WithContext(ctx))
		}
	}
}

func extractBearerToken(header string) (string, error) {
	const prefix = "Bearer "
	if header == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	if !strings.HasPrefix(header, prefix) {
		return "", fmt.Errorf("invalid authorization scheme")
	}
	token := strings.TrimSpace(header[len(prefix):])
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return token, nil
}
