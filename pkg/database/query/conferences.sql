-- name: CreateConference :one
INSERT INTO conferences (
  url, title, venue, city
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateConferenceDetails :one
UPDATE conferences SET (
  title, venue, city
) = ($2, $3, $4)
WHERE id = $1
RETURNING *;

-- name: GetConferences :many
SELECT * FROM conferences;

-- name: DeleteConference :exec
DELETE FROM conferences
WHERE id = $1;
