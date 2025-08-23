package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type memoryStore struct {
	sessionsByToken map[string]*UserSession
	sessionsBySID   map[uint]*UserSession
	nextSessionId   uint
	lock            sync.RWMutex
}

func NewMemoryStore() Service {
	return &memoryStore{
		sessionsByToken: make(map[string]*UserSession),
		sessionsBySID:   make(map[uint]*UserSession),
	}
}

func (s *memoryStore) GetByToken(token string) *UserSession {
	if token == "" {
		return nil
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.sessionsByToken[token]
}

func (s *memoryStore) GetBySID(sid uint) *UserSession {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.sessionsBySID[sid]
}

func (s *memoryStore) Create(uid int32, username string, ip string, ua string, admin bool) (*UserSession, error) {
	token := generateSessionToken()

	s.lock.Lock()
	defer s.lock.Unlock()

	sessionId := s.nextSessionId
	s.nextSessionId++

	_, sidExists := s.sessionsBySID[sessionId]
	_, tokenExists := s.sessionsByToken[token]

	// should never realistically happen but still theoretically possible
	if sidExists || tokenExists {
		return nil, fmt.Errorf("session conflict")
	}

	session := &UserSession{
		UserID:    uid,
		SessionID: sessionId,
		Token:     token,
		Username:  username,
		IP:        ip,
		UserAgent: ua,
		LoginTime: time.Now(),
		Admin:     admin,
	}
	s.sessionsByToken[token] = session
	s.sessionsBySID[sessionId] = session

	return session, nil
}

func (s *memoryStore) Destroy(sid uint) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	session := s.sessionsBySID[sid]
	if session == nil {
		return fmt.Errorf("session does not exist")
	}

	delete(s.sessionsBySID, sid)
	delete(s.sessionsByToken, session.Token)
	return nil
}

func generateSessionToken() string {
	b := make([]byte, 100)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
