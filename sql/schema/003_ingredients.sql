-- +goose Up
CREATE TABLE ingredients (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        name TEXT NOT NULL,
                        unit TEXT NOT NULL,
						owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE ingredients;
