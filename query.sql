-- USERS

-- name: GetUser :one
SELECT id, username, password, created_at
FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id, created_at;

