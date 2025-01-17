-- name: GetFavouritesForUser :many
SELECT * FROM favourites
WHERE user_id = $1;

-- name: CreateFavourite :one
INSERT INTO favourites (
  user_id, event_guid, event_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteFavourite :exec
DELETE FROM favourites
WHERE id = $1;

-- name: DeleteFavouriteByEventDetails :execrows
DELETE FROM favourites
WHERE (event_guid = $1 OR event_id = $2) AND user_id = $3;