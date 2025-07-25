-- name: CreateLoginToken :one
INSERT INTO login_tokens (token, user_id, created_at, expire_at)
VALUES (
    gen_random_uuid(), $1, NOW(), NOW() + interval '1 hour'
)
RETURNING *;

-- name: GetLoginByToken :one
SELECT * FROM login_tokens WHERE token = $1;

-- name: GetTokenByUserID :one
SELECT * FROM login_tokens WHERE user_id = $1;

-- name: DeleteToken :exec
DELETE FROM login_tokens WHERE token = $1;
