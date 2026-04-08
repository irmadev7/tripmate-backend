-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO users (email, name, password) VALUES (
    'test@test.com', 'Test User', '$2a$12$yG.5dY7Rni27fVBAy5Cs3.Z0ep8u3kl4wYvrHWe214eN6L.EfeXPi'), --password
    ('test2@test.com', 'Test User 2', '$2a$12$1MYmrqAVDTC8mobLc.UYeefP56OX8nScwtkzEPTq3A.Tkki.H0Z.2'), --password2
    ('test3@test.com', 'Test User 3', '$2a$12$xmSQZ4WGDvxr3h0nSqChJOAZzCPUhKN88z7WF4WTBjluSu2ZyNl6O'); --password3
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM users WHERE email IN ('test@test.com', 'test2@test.com', 'test3@test.com');
-- +goose StatementEnd
