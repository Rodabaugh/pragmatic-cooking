-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email_addr, email_verified)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email_addr = $1;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;
