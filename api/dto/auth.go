package dto

type LoginBasicRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginOAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type LoginOAuthOutboundResponse struct {
	URL string `json:"url"`
}

type LoginResponse struct {
	ID       int32  `json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

type LoginOptionsResponse struct {
	Options []LoginOption `json:"options"`
}

type LoginOption struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}
