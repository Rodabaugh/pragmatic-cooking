-- name: CreateIngredient :one
INSERT INTO ingredients (id, created_at, updated_at, name, unit, owner_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetAllIngredients :many
SELECT * FROM ingredients;

-- name: GetIngredientByID :one
SELECT * FROM ingredients WHERE id = $1;

-- name: GetIngredientsByName :many
SELECT * FROM ingredients WHERE name = $1;

-- name: DeleteIngrendientByID :exec
DELETE FROM ingredients WHERE id = $1;
