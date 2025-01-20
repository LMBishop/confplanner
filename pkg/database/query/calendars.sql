-- name: CreateCalendar :one
INSERT INTO calendars (
  user_id, name, key
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetCalendarForUser :one
SELECT * FROM calendars
WHERE user_id = $1 LIMIT 1;

-- name: GetCalendarByName :one
SELECT * FROM calendars
WHERE name = $1 LIMIT 1;

-- name: DeleteCalendar :execrows
DELETE FROM calendars
WHERE user_id = $1;

-- name: DeleteCalendarByName :execrows
DELETE FROM calendars
WHERE name = $1;
