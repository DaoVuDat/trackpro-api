-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE profile
    ADD COLUMN about text,
    ADD COLUMN phone text;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE project
    DROP paid,
    DROP phone;