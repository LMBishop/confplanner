package dto

import "fmt"

type Response interface {
	Error() string
	Status() int
}

type OkResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

var _ Response = (*OkResponse)(nil)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var _ Response = (*ErrorResponse)(nil)

func (r *OkResponse) Status() int {
	return r.Code
}

func (r *OkResponse) Error() string {
	return fmt.Sprintf("HTTP status %d", r.Code)
}

func (r *ErrorResponse) Status() int {
	return r.Code
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("HTTP status %d: %s", r.Code, r.Message)
}
