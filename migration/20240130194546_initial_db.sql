-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "account_type" AS ENUM (
    'client',
    'admin'
    );

CREATE TYPE "account_status" AS ENUM (
    'pending',
    'activated'
    );

CREATE TYPE "project_status" AS ENUM (
    'registering',
    'progressing',
    'finished'
    );

CREATE TABLE IF NOT EXISTS account
(
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   text UNIQUE              NOT NULL,
    type       account_type             NOT NULL DEFAULT 'client',
    status     account_status           NOT NULL DEFAULT 'pending',
    created_at timestamp with time zone NOT NULL DEFAULT (now()),
    updated_at timestamp with time zone NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS password
(
    user_id       uuid UNIQUE NOT NULL,
    hash_password text        NOT NULL
);

ALTER TABLE password
    ADD CONSTRAINT password_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;

CREATE INDEX password_user_id_idx ON password (user_id);

CREATE TABLE IF NOT EXISTS profile
(
    user_id    uuid UNIQUE              NOT NULL,
    first_name text                     NOT NULL,
    last_name  text                     NOT NULL,
    image_url  text,
    created_at timestamp with time zone NOT NULL DEFAULT (now()),
    updated_at timestamp with time zone NOT NULL DEFAULT (now())
);

ALTER TABLE profile
    ADD CONSTRAINT profile_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;

CREATE INDEX profile_user_id_idx ON profile (user_id);

CREATE TABLE IF NOT EXISTS project
(
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     uuid                     NOT NULL,
    name        text,
    description text,
    price       integer,
    paid        integer,
    status      project_status           NOT NULL DEFAULT 'registering',
    start_time  timestamp with time zone,
    end_time    timestamp with time zone,
    created_at  timestamp with time zone NOT NULL DEFAULT (now()),
    updated_at  timestamp with time zone NOT NULL DEFAULT (now())
);

ALTER TABLE project
    ADD CONSTRAINT project_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE;

COMMENT ON COLUMN project.description IS 'Which technologies and algorithms used';

COMMENT ON COLUMN project.price IS 'Price of the project';

COMMENT ON COLUMN project.paid IS 'How much money client paid';

CREATE TABLE IF NOT EXISTS payment_history
(
    id         BIGSERIAL PRIMARY KEY,
    project_id uuid                     NOT NULL,
    amount     integer                  NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT (now()),
    updated_at timestamp with time zone NOT NULL DEFAULT (now())
);

ALTER TABLE payment_history
    ADD CONSTRAINT payment_history_project_id_fkey
        FOREIGN KEY (project_id) REFERENCES project (id) ON DELETE CASCADE;


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE payment_history
    DROP CONSTRAINT payment_history_project_id_fkey;

DROP TABLE payment_history;

ALTER TABLE project
    DROP CONSTRAINT project_user_id_fkey;

DROP TABLE project;

ALTER TABLE profile
    DROP CONSTRAINT profile_user_id_fkey;

DROP TABLE profile;

ALTER TABLE password
    DROP CONSTRAINT password_user_id_fkey;

DROP INDEX password_user_id_idx;

DROP TABLE password;

DROP TABLE account;

DROP TYPE account_type;

DROP TYPE account_status;

DROP TYPE project_status;

DROP EXTENSION IF EXISTS "uuid-ossp";