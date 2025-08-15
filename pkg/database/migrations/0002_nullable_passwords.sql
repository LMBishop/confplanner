-- +goose Up
ALTER TABLE users DROP CONSTRAINT valid_hash;
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
ALTER TABLE users ADD CONSTRAINT valid_hash CHECK (length(password) = 60 OR password IS NULL);
