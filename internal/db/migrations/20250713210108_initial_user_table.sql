-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY default gen_random_uuid(),
    username TEXT NOT NULL,
    hashed_password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_username ON users (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
