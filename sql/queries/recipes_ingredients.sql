-- name: CreateRecipeIngredient :one
INSERT INTO recipes_ingredients (recipe_id, ingredient_id, quantity, created_at, updated_at)
VALUES (
    $1, $2, $3, NOW(), NOW()
)
RETURNING *;

-- name: GetAllRecipeIngredients :many
SELECT * FROM recipes_ingredients;

-- name: GetIngredientsByRecipe :many
SELECT
	ri.recipe_id,
	ri.ingredient_id,
    i.name AS ingredient_name,
    i.unit,
    ri.quantity
FROM
    ingredients i
JOIN
    recipes_ingredients ri ON i.id = ri.ingredient_id
WHERE
    ri.recipe_id = $1;

-- name: GetRecipesByIngredient :many
SELECT
	r.id,
    r.name AS recipe_name,
    r.description,
	r.link
FROM
    recipes r
JOIN
    recipes_ingredients ri ON r.id = ri.recipe_id
WHERE
    ri.ingredient_id = $1;

-- name: GetIngredientDependantCount :one
SELECT COUNT(*) FROM recipes_ingredients WHERE ingredient_id = $1;

-- name: DeleteRecipeIngredient :exec
DELETE FROM
    recipes_ingredients
WHERE
    recipe_id = $1 AND ingredient_id = $2;
