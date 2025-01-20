package calendar

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetCalendarForUser(id int32) (*sqlc.Calendar, error)
	GetCalendarByName(name string) (*sqlc.Calendar, error)
	CreateCalendarForUser(id int32) (*sqlc.Calendar, error)
	DeleteCalendarForUser(id int32) error
}

var (
	ErrCalendarNotFound = errors.New("calendar not found")
)

type service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) Service {
	return &service{
		pool: pool,
	}
}

func (s *service) GetCalendarForUser(id int32) (*sqlc.Calendar, error) {
	queries := sqlc.New(s.pool)

	calendar, err := queries.GetCalendarForUser(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCalendarNotFound
		}
		return nil, err
	}

	return &calendar, nil
}

func (s *service) GetCalendarByName(name string) (*sqlc.Calendar, error) {
	queries := sqlc.New(s.pool)

	calendar, err := queries.GetCalendarByName(context.Background(), name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCalendarNotFound
		}
		return nil, err
	}

	return &calendar, nil
}

func (s *service) CreateCalendarForUser(id int32) (*sqlc.Calendar, error) {
	queries := sqlc.New(s.pool)

	name, err := randomString(16)
	if err != nil {
		return nil, fmt.Errorf("could not generate random string: %w", err)
	}

	key, err := randomString(32)
	if err != nil {
		return nil, fmt.Errorf("could not generate random string: %w", err)
	}

	calendar, err := queries.CreateCalendar(context.Background(), sqlc.CreateCalendarParams{
		UserID: id,
		Name:   name,
		Key:    key,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create calendar: %w", err)
	}

	return &calendar, nil
}

func (s *service) DeleteCalendarForUser(id int32) error {
	queries := sqlc.New(s.pool)

	_, err := queries.DeleteCalendar(context.Background(), id)
	if err != nil {
		return fmt.Errorf("could not delete calendar: %w", err)
	}

	return nil
}

func randomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
