package dto

type GetCalendarResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	URL  string `json:"url"`
}

type CreateCalendarResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
	URL  string `json:"url"`
}
