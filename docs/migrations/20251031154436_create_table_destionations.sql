-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS destinations (
	id BIGSERIAL PRIMARY KEY,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	itinerary_id bigint NOT NULL,
	name varchar(255) NOT NULL,
	note text NOT NULL,
	day int NOT NULL,
	"order" int NOT NULL,
    CONSTRAINT fk_itineraries_destinations FOREIGN KEY (itinerary_id) REFERENCES itineraries(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE destinations;
-- +goose StatementEnd
