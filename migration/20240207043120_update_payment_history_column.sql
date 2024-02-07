-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE payment_history
    ADD user_id uuid NOT NULL,
    ADD CONSTRAINT payment_history_user_id_fk
        FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE payment_history
    DROP CONSTRAINT payment_history_user_id_fk,
    DROP user_id;