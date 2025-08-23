package session

import "time"

type Service interface {
	GetByToken(token string) *UserSession
	GetBySID(sid uint) *UserSession
	Create(uid int32, username string, ip string, ua string, admin bool) (*UserSession, error)
	Destroy(sid uint) error
}

type UserSession struct {
	UserID    int32
	SessionID uint
	Token     string
	Username  string
	IP        string
	LoginTime time.Time
	UserAgent string
	Admin     bool
}
