-- +goose Up
CREATE TABLE question_banks (
    id          uuid PRIMARY KEY,
    owner_id    uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        text        NOT NULL,
    description text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT question_banks_name_len CHECK (char_length(name) BETWEEN 1 AND 120)
);

CREATE INDEX question_banks_owner_idx ON question_banks (owner_id);

CREATE TRIGGER question_banks_set_updated_at
    BEFORE UPDATE ON question_banks
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- +goose Down
DROP TABLE IF EXISTS question_banks;
