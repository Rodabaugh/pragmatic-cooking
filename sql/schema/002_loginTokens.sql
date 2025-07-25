-- +goose Up
CREATE TABLE login_tokens (
    token UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    expire_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE login_tokens;
