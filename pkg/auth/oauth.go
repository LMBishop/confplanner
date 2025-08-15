package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/LMBishop/confplanner/pkg/user"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
)

type OIDCAuthProvider struct {
	name                     string
	userService              user.Service
	oauthConfig              *oauth2.Config
	oidcProvider             *oidc.Provider
	oidcVerifier             *oidc.IDTokenVerifier
	loginFilter              string
	loginFilterAllowedValues []string
	userSyncFilter           string
	states                   map[string]*oidcState
	lock                     sync.RWMutex
}

type oidcState struct {
	expiry    time.Time
	ip        string
	userAgent string
}

var (
	ErrStateVerificationFailed = errors.New("state verification failed")
	ErrInvalidState            = errors.New("invalid state")
	ErrMissingIDToken          = errors.New("missing ID token")
	ErrNotAuthorised           = errors.New("not authorised")
	ErrUserSyncFailed          = errors.New("user sync failed")
)

func NewOIDCAuthProvider(userService user.Service, name, clientID, clientSecret, endpoint, callbackURL, loginFilter, userSyncFilter string, loginFilterAllowedValues []string) (AuthProvider, error) {
	provider, err := oidc.NewProvider(context.Background(), endpoint)
	if err != nil {
		return nil, err
	}

	return &OIDCAuthProvider{
		name:        name,
		userService: userService,
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  callbackURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		oidcProvider:             provider,
		oidcVerifier:             provider.Verifier(&oidc.Config{ClientID: clientID}),
		loginFilter:              loginFilter,
		loginFilterAllowedValues: loginFilterAllowedValues,
		userSyncFilter:           userSyncFilter,
		states:                   make(map[string]*oidcState),
	}, nil
}

func (p *OIDCAuthProvider) StartJourney(ip string, userAgent string) (string, error) {
	b := make([]byte, 50)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)

	p.lock.Lock()
	defer p.lock.Unlock()

	p.states[state] = &oidcState{
		expiry:    time.Now().Add(time.Minute * 5),
		ip:        ip,
		userAgent: userAgent,
	}

	return p.oauthConfig.AuthCodeURL(state), nil
}

func (p *OIDCAuthProvider) CompleteJourney(ctx context.Context, authCode string, state string, ip string, userAgent string) (*sqlc.User, error) {
	var s *oidcState

	p.lock.Lock()
	s = p.states[state]
	delete(p.states, state)
	p.lock.Unlock()

	if s == nil {
		return nil, ErrInvalidState
	}

	//if time.Now().After(s.expiry) || s.ip != ip || s.userAgent != userAgent {
	//	return nil, ErrStateVerificationFailed
	//}
	if time.Now().After(s.expiry) || s.userAgent != userAgent {
		return nil, ErrStateVerificationFailed
	}

	oauth2Token, err := p.oauthConfig.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, ErrMissingIDToken
	}

	_, err = p.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	claims, err := getRawClaims(rawIDToken)
	if err != nil {
		return nil, err
	}

	if p.loginFilter != "" {
		rolesClaim := gjson.Get(claims, p.loginFilter)
		if !rolesClaim.Exists() {
			return nil, fmt.Errorf("cannot verify authorisation as '%s' is missing from claims", p.loginFilter)
		}
		roles := rolesClaim.Array()
		var authorisation bool
	out:
		for _, allowedRole := range p.loginFilterAllowedValues {
			for _, role := range roles {
				if role.Str == allowedRole {
					authorisation = true
					break out
				}
			}
		}
		if !authorisation {
			return nil, ErrNotAuthorised
		}
	}

	usernameClaim := gjson.Get(claims, p.userSyncFilter)
	if !usernameClaim.Exists() {
		return nil, fmt.Errorf("cannot sync user as '%s' is missing from claims", p.userSyncFilter)
	}
	username := usernameClaim.Str

	u, err := p.userService.GetUserByName(username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			u, err = p.userService.CreateUser(username, "")
			if err != nil {
				return nil, errors.Join(ErrUserSyncFailed, err)
			}
		} else {
			return nil, errors.Join(ErrUserSyncFailed, err)
		}
	}

	return u, nil
}

func (p *OIDCAuthProvider) Name() string {
	return p.name
}

func (p *OIDCAuthProvider) Type() string {
	return "oidc"
}

func getRawClaims(p string) (string, error) {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("malformed jwt payload: %w", err)
	}
	return string(payload[:]), nil
}
