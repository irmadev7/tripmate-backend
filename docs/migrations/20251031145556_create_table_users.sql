-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	username text UNIQUE NOT NULL,
	password text NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
