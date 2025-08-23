package middleware

import (
	"net/http"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

func MustAuthoriseAdmin(service user.Service, store session.Service) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value("session").(*session.UserSession)

			if !session.Admin {
				dto.WriteDto(w, r, &dto.ErrorResponse{
					Code:    http.StatusForbidden,
					Message: "Forbidden",
				})
				return
			}

			next(w, r)
		}
	}
}
