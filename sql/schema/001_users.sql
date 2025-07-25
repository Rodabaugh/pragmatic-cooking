-- +goose Up
CREATE TABLE users (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        name TEXT NOT NULL,
                        email_addr TEXT NOT NULL UNIQUE,
						email_verified BOOLEAN not NULL);

-- +goose Down
DROP TABLE users;
