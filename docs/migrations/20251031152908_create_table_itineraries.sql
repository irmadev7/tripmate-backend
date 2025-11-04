-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS itineraries (
	id BIGSERIAL PRIMARY KEY,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	title varchar(255) NOT NULL,
	description text NOT NULL,
	start_date date NOT NULL,
	end_date date NOT NULL,
	user_id bigint NOT NULL,
    CONSTRAINT fk_users_itineraries FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE itineraries;
-- +goose StatementEnd
