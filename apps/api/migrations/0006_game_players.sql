-- +goose Up
CREATE TABLE game_players (
    id              uuid PRIMARY KEY,
    game_id         uuid        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    display_name    text        NOT NULL,
    score           integer     NOT NULL DEFAULT 0,
    session_token   text        NOT NULL,
    joined_at       timestamptz NOT NULL DEFAULT now(),
    left_at         timestamptz,
    CONSTRAINT game_players_display_name_len CHECK (char_length(display_name) BETWEEN 1 AND 32),
    CONSTRAINT game_players_left_after_join  CHECK (left_at IS NULL OR left_at >= joined_at)
);

-- Display names unique per game, case-insensitive. Prevents "Alice" and "alice" collisions.
CREATE UNIQUE INDEX game_players_display_name_unique
    ON game_players (game_id, lower(display_name));

-- Session tokens globally unique for reconnect lookups.
CREATE UNIQUE INDEX game_players_session_token_unique
    ON game_players (session_token);

CREATE INDEX game_players_game_idx ON game_players (game_id);

-- +goose Down
DROP TABLE IF EXISTS game_players;
