-- +goose Up
CREATE TABLE recipes_ingredients (PRIMARY KEY (recipe_id, ingredient_id),
						recipe_id UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
						ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
						quantity NUMERIC(10, 2) NOT NULL,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL);

-- +goose Down
DROP TABLE recipes_ingredients;
