-- USERS

-- name: GetUser :one
SELECT id, username, password, created_at
FROM users
WHERE username = $1;

-- name: CreateUser :execlastid
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id, created_at;


-- MESSAGES
-- name: SelectSentMessages :many
SELECT id, sender_id, rec_id, body, created_at
FROM messages
WHERE sender_id = $1
ORDER BY created_at DESC;

-- name: SelectRecMessages :many
SELECT id, sender_id, rec_id, body, created_at
FROM messages
WHERE rec_id = $1
ORDER BY created_at DESC;


-- name: SelectConvMessages :many
SELECT id, sender_id, rec_id, body, created_at
FROM messages
WHERE (sender_id = $1 AND rec_id = $2) OR (sender_id = $2 AND rec_id = $1)
ORDER BY created_at DESC;


-- name: CreateMessage :execlastid
INSERT INTO messages (sender_id, rec_id, body)
VALUES ($1, $2, $3)
RETURNING id, created_at;
