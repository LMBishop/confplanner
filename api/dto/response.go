package dto

import "fmt"

type OkResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *OkResponse) Error() string {
	return fmt.Sprintf("HTTP status %d", r.Code)
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("HTTP status %d: %s", r.Code, r.Message)
}
