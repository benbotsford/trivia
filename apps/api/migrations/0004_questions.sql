-- +goose Up
CREATE TYPE question_type AS ENUM ('text', 'multiple_choice');

CREATE TABLE questions (
    id              uuid PRIMARY KEY,
    bank_id         uuid          NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    type            question_type NOT NULL,
    prompt          text          NOT NULL,
    correct_answer  text          NOT NULL,
    choices         jsonb,
    points          integer       NOT NULL DEFAULT 1000,
    position        integer       NOT NULL DEFAULT 0,
    created_at      timestamptz   NOT NULL DEFAULT now(),
    updated_at      timestamptz   NOT NULL DEFAULT now(),
    CONSTRAINT questions_prompt_len       CHECK (char_length(prompt) BETWEEN 1 AND 2000),
    CONSTRAINT questions_points_positive  CHECK (points > 0),
    CONSTRAINT questions_position_nonneg  CHECK (position >= 0),
    CONSTRAINT questions_choices_shape CHECK (
        (type <> 'multiple_choice')
        OR (choices IS NOT NULL AND jsonb_typeof(choices) = 'array' AND jsonb_array_length(choices) BETWEEN 2 AND 10)
    )
);

CREATE INDEX questions_bank_position_idx ON questions (bank_id, position);

CREATE TRIGGER questions_set_updated_at
    BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS questions;
DROP TYPE IF EXISTS question_type;
