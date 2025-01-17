package favourites

import (
	"context"
	"errors"
	"fmt"

	"github.com/LMBishop/confplanner/pkg/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetFavouritesForUser(id int32) (*[]sqlc.Favourite, error)
	CreateFavouriteForUser(id int32, eventGuid pgtype.UUID, eventId *int32) (*sqlc.Favourite, error)
	DeleteFavouriteForUserByEventDetails(id int32, eventGuid pgtype.UUID, eventId *int32) error
}

var (
	ErrImproperType = errors.New("improper type")
	ErrNotFound     = errors.New("not found")
)

type service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) Service {
	return &service{
		pool: pool,
	}
}

func (s *service) CreateFavouriteForUser(id int32, eventGuid pgtype.UUID, eventId *int32) (*sqlc.Favourite, error) {
	queries := sqlc.New(s.pool)

	var pgEventId pgtype.Int4
	if eventId != nil {
		pgEventId = pgtype.Int4{
			Int32: *eventId,
			Valid: true,
		}
	}

	favourite, err := queries.CreateFavourite(context.Background(), sqlc.CreateFavouriteParams{
		UserID:    id,
		EventGuid: eventGuid,
		EventID:   pgEventId,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create favourite: %w", err)
	}

	return &favourite, nil
}

func (s *service) GetFavouritesForUser(id int32) (*[]sqlc.Favourite, error) {
	queries := sqlc.New(s.pool)

	favourites, err := queries.GetFavouritesForUser(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			empty := make([]sqlc.Favourite, 0)
			return &empty, nil
		}
		return nil, fmt.Errorf("could not fetch user: %w", err)
	}

	return &favourites, nil
}

func (s *service) DeleteFavouriteForUserByEventDetails(id int32, eventGuid pgtype.UUID, eventId *int32) error {
	queries := sqlc.New(s.pool)

	var pgEventId pgtype.Int4
	if eventId != nil {
		pgEventId = pgtype.Int4{
			Int32: *eventId,
			Valid: true,
		}
	}
	rowsAffected, err := queries.DeleteFavouriteByEventDetails(context.Background(), sqlc.DeleteFavouriteByEventDetailsParams{
		EventGuid: eventGuid,
		EventID:   pgEventId,
		UserID:    id,
	})
	if err != nil {
		return fmt.Errorf("could not delete favourite: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
