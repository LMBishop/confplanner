package dto

import (
	"time"
)

type GetScheduleResponse struct {
	Schedule    interface{} `json:"schedule"`
	LastUpdated time.Time   `json:"lastUpdated"`
}
