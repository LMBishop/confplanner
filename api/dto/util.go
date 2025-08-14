package dto

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ReadDto(r *http.Request, o interface{}) Response {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(o); err != nil {
		return &ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Errorf("Invalid request (%w)", err).Error(),
		}
	}

	if err := validate.Struct(o); err != nil {
		return &ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}

func WrapResponseFunc(dtoFunc func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteDto(w, r, dtoFunc(w, r))
	}
}

func WriteDto(w http.ResponseWriter, r *http.Request, o error) {
	if o, ok := o.(Response); ok {
		data, err := json.Marshal(o)
		if err != nil {
			w.WriteHeader(500)
			slog.Error("could not serialise JSON", "error", err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(o.Status())
		w.Write(data)
	} else {
		w.WriteHeader(500)
		slog.Error("internal server error handling request", "error", o)
	}
}
