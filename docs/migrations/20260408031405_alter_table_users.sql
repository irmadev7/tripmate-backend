-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users 
    DROP COLUMN username,
    ADD COLUMN email TEXT UNIQUE DEFAULT NULL,
    ADD COLUMN name TEXT DEFAULT NULL,
    ADD CONSTRAINT unique_email UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users 
    DROP COLUMN email,
    DROP COLUMN name,
    DROP CONSTRAINT unique_email,
    ADD COLUMN username TEXT UNIQUE NOT NULL,
    ADD INDEX idx_users_username (username);
-- +goose StatementEnd
