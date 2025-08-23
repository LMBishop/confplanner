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
	GetAllFavouritesForUser(id int32) (*[]sqlc.Favourite, error)
	GetFavouritesForUserConference(id int32, conference int32) (*[]sqlc.Favourite, error)
	CreateFavouriteForUser(id int32, eventGUID pgtype.UUID, eventID *int32, conferenceID int32) (*sqlc.Favourite, error)
	DeleteFavouriteForUserByEventDetails(id int32, eventGUID pgtype.UUID, eventID *int32, conferenceID int32) error
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

func (s *service) CreateFavouriteForUser(userID int32, eventGUID pgtype.UUID, eventID *int32, conferenceID int32) (*sqlc.Favourite, error) {
	queries := sqlc.New(s.pool)

	var pgEventID pgtype.Int4
	if eventID != nil {
		pgEventID = pgtype.Int4{
			Int32: *eventID,
			Valid: true,
		}
	}

	favourite, err := queries.CreateFavourite(context.Background(), sqlc.CreateFavouriteParams{
		UserID:       userID,
		EventGuid:    eventGUID,
		EventID:      pgEventID,
		ConferenceID: conferenceID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create favourite: %w", err)
	}

	return &favourite, nil
}

func (s *service) GetAllFavouritesForUser(userID int32) (*[]sqlc.Favourite, error) {
	queries := sqlc.New(s.pool)

	favourites, err := queries.GetFavouritesForUser(context.Background(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			empty := make([]sqlc.Favourite, 0)
			return &empty, nil
		}
		return nil, fmt.Errorf("could not fetch user: %w", err)
	}

	return &favourites, nil
}

func (s *service) GetFavouritesForUserConference(userID int32, conferenceID int32) (*[]sqlc.Favourite, error) {
	queries := sqlc.New(s.pool)

	favourites, err := queries.GetFavouritesForUserConference(context.Background(), sqlc.GetFavouritesForUserConferenceParams{
		UserID:       userID,
		ConferenceID: conferenceID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			empty := make([]sqlc.Favourite, 0)
			return &empty, nil
		}
		return nil, fmt.Errorf("could not fetch user: %w", err)
	}

	return &favourites, nil
}

func (s *service) DeleteFavouriteForUserByEventDetails(id int32, eventGUID pgtype.UUID, eventID *int32, conferenceID int32) error {
	queries := sqlc.New(s.pool)

	var pgEventID pgtype.Int4
	if eventID != nil {
		pgEventID = pgtype.Int4{
			Int32: *eventID,
			Valid: true,
		}
	}
	rowsAffected, err := queries.DeleteFavouriteByEventDetails(context.Background(), sqlc.DeleteFavouriteByEventDetailsParams{
		EventGuid:    eventGUID,
		EventID:      pgEventID,
		UserID:       id,
		ConferenceID: conferenceID,
	})
	if err != nil {
		return fmt.Errorf("could not delete favourite: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
