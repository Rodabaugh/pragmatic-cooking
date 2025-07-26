-- name: CreateRecipe :one
INSERT INTO recipes (id, created_at, updated_at, name, description, link, owner_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4
)
RETURNING *;

-- name: GetAllRecipes :many
SELECT * FROM recipes;

-- name: GetRecipeByID :one
SELECT * FROM recipes WHERE id = $1;

-- name: GetRecipesByName :many
SELECT * FROM recipes WHERE name = $1;

-- name: DeleteRecipeByID :exec
DELETE FROM recipes WHERE id = $1;
