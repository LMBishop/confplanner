package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(username string, password string) (*sqlc.User, error)
	GetUserByName(username string) (*sqlc.User, error)
	GetUserByID(id int32) (*sqlc.User, error)
}

var (
	ErrUserExists                = errors.New("user already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrNotAcceptingRegistrations = errors.New("not currently accepting registrations")
)

type service struct {
	pool                   *pgxpool.Pool
	acceptingRegistrations bool
}

func NewService(pool *pgxpool.Pool, acceptingRegistrations bool) Service {
	return &service{
		pool:                   pool,
		acceptingRegistrations: acceptingRegistrations,
	}
}

func (s *service) CreateUser(username string, password string) (*sqlc.User, error) {
	if !s.acceptingRegistrations {
		return nil, ErrNotAcceptingRegistrations
	}

	var passwordHash pgtype.Text
	queries := sqlc.New(s.pool)

	if password != "" {
		var passwordBytes = []byte(password)

		hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("could not hash password: %w", err)
		}

		passwordHash = pgtype.Text{
			String: string(hash),
			Valid:  true,
		}
	} else {
		passwordHash = pgtype.Text{
			Valid: false,
		}
	}

	user, err := queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Username: strings.ToLower(username),
		Password: passwordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	return &user, nil
}

func (s *service) GetUserByName(username string) (*sqlc.User, error) {
	queries := sqlc.New(s.pool)

	user, err := queries.GetUserByName(context.Background(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("could not fetch user: %w", err)
	}

	return &user, nil
}

func (s *service) GetUserByID(id int32) (*sqlc.User, error) {
	queries := sqlc.New(s.pool)

	user, err := queries.GetUserByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("could not fetch user: %w", err)
	}

	return &user, nil
}
