package auth

import (
	"errors"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/LMBishop/confplanner/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthProvider struct {
	userService user.Service
}

func NewBasicAuthProvider(userService user.Service) AuthProvider {
	return &BasicAuthProvider{
		userService: userService,
	}
}

func (p *BasicAuthProvider) Authenticate(username string, password string) (*sqlc.User, error) {
	random, err := bcrypt.GenerateFromPassword([]byte("00000000"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := p.userService.GetUserByName(username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			bcrypt.CompareHashAndPassword(random, []byte(password))
			return nil, nil
		}
		return nil, err
	}
	if !u.Password.Valid {
		bcrypt.CompareHashAndPassword(random, []byte(password))
		return nil, nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (p *BasicAuthProvider) Name() string {
	return "Basic"
}

func (p *BasicAuthProvider) Type() string {
	return "basic"
}
