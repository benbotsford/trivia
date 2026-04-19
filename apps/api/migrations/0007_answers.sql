-- +goose Up
CREATE TABLE answers (
    id              uuid PRIMARY KEY,
    game_id         uuid        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    question_id     uuid        NOT NULL REFERENCES questions(id),
    player_id       uuid        NOT NULL REFERENCES game_players(id) ON DELETE CASCADE,
    answer          text        NOT NULL,
    is_correct      boolean     NOT NULL,
    points_awarded  integer     NOT NULL DEFAULT 0,
    submitted_at    timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT answers_points_nonneg CHECK (points_awarded >= 0)
);

-- One answer per player per question.
CREATE UNIQUE INDEX answers_one_per_player_question
    ON answers (game_id, question_id, player_id);

CREATE INDEX answers_game_idx     ON answers (game_id);
CREATE INDEX answers_player_idx   ON answers (player_id);
CREATE INDEX answers_question_idx ON answers (question_id);

-- +goose Down
DROP TABLE IF EXISTS answers;
