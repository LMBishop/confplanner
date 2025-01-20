// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: calendars.sql

package sqlc

import (
	"context"
)

const createCalendar = `-- name: CreateCalendar :one
INSERT INTO calendars (
  user_id, name, key
) VALUES (
  $1, $2, $3
)
RETURNING id, user_id, name, key
`

type CreateCalendarParams struct {
	UserID int32  `json:"user_id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
}

func (q *Queries) CreateCalendar(ctx context.Context, arg CreateCalendarParams) (Calendar, error) {
	row := q.db.QueryRow(ctx, createCalendar, arg.UserID, arg.Name, arg.Key)
	var i Calendar
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const deleteCalendar = `-- name: DeleteCalendar :execrows
DELETE FROM calendars
WHERE user_id = $1
`

func (q *Queries) DeleteCalendar(ctx context.Context, userID int32) (int64, error) {
	result, err := q.db.Exec(ctx, deleteCalendar, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteCalendarByName = `-- name: DeleteCalendarByName :execrows
DELETE FROM calendars
WHERE name = $1
`

func (q *Queries) DeleteCalendarByName(ctx context.Context, name string) (int64, error) {
	result, err := q.db.Exec(ctx, deleteCalendarByName, name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getCalendarByName = `-- name: GetCalendarByName :one
SELECT id, user_id, name, key FROM calendars
WHERE name = $1 LIMIT 1
`

func (q *Queries) GetCalendarByName(ctx context.Context, name string) (Calendar, error) {
	row := q.db.QueryRow(ctx, getCalendarByName, name)
	var i Calendar
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const getCalendarForUser = `-- name: GetCalendarForUser :one
SELECT id, user_id, name, key FROM calendars
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetCalendarForUser(ctx context.Context, userID int32) (Calendar, error) {
	row := q.db.QueryRow(ctx, getCalendarForUser, userID)
	var i Calendar
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Key,
	)
	return i, err
}
