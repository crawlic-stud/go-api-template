-- name: CreateUser :exec
INSERT INTO users (username, hashed_password) VALUES ($1, $2);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;