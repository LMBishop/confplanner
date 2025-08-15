package auth

import (
	"fmt"
	"sync"
)

type Service interface {
	GetAuthProvider(string) AuthProvider
	GetAuthProviders() []string
	RegisterAuthProvider(string, AuthProvider) error
}

type AuthProvider interface {
	Name() string
	Type() string
}

type service struct {
	authProviders map[string]AuthProvider
	order         []string
	lock          sync.Mutex
}

func NewService() Service {
	return &service{
		authProviders: make(map[string]AuthProvider),
	}
}

func (s *service) GetAuthProvider(name string) AuthProvider {
	return s.authProviders[name]
}

func (s *service) GetAuthProviders() []string {
	return s.order
}

func (s *service) RegisterAuthProvider(name string, provider AuthProvider) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.authProviders[name]; ok {
		return fmt.Errorf("duplicate auth provider: %s", name)
	}
	s.order = append(s.order, name)
	s.authProviders[name] = provider
	return nil
}
