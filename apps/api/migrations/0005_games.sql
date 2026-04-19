-- +goose Up
CREATE TYPE game_status AS ENUM ('lobby', 'in_progress', 'completed', 'cancelled');

CREATE TABLE games (
    id                    uuid PRIMARY KEY,
    code                  text        NOT NULL,
    status                game_status NOT NULL DEFAULT 'lobby',
    host_id               uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_id               uuid        NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    current_question_idx  integer     NOT NULL DEFAULT 0,
    started_at            timestamptz,
    ended_at              timestamptz,
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT games_code_format            CHECK (code ~ '^[A-Z0-9]{6}$'),
    CONSTRAINT games_current_question_nonneg CHECK (current_question_idx >= 0),
    CONSTRAINT games_timeline CHECK (
        (started_at IS NULL OR started_at >= created_at)
        AND (ended_at IS NULL OR started_at IS NOT NULL)
        AND (ended_at IS NULL OR ended_at >= started_at)
    )
);

-- Only enforce uniqueness among active games. Completed/cancelled games keep their
-- historical code, but the string can be recycled once they're no longer live.
CREATE UNIQUE INDEX games_code_active_unique
    ON games (code)
    WHERE status IN ('lobby', 'in_progress');

CREATE INDEX games_host_idx ON games (host_id);
CREATE INDEX games_status_idx ON games (status);

CREATE TRIGGER games_set_updated_at
    BEFORE UPDATE ON games
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS games;
DROP TYPE IF EXISTS game_status;
