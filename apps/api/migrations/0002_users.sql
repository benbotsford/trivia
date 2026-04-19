-- +goose Up
CREATE TABLE users (
    id              uuid PRIMARY KEY,
    auth0_sub       text        NOT NULL UNIQUE,
    email           citext      UNIQUE,
    display_name    text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS users;
