package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

func MustAuthenticate(service user.Service, store session.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var sessionToken string
			for _, cookie := range r.Cookies() {
				if cookie.Name == "confplanner_session" {
					sessionToken = cookie.Value
					break
				}
			}

			s := store.GetByToken(sessionToken)
			if s == nil {
				dto.WriteDto(w, r, &dto.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized",
				})
				return
			}

			_, err := service.GetUserByID(s.UserID)
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

			ctx := context.WithValue(r.Context(), "session", s)

			next(w, r.WithContext(ctx))
		}
	}
}
