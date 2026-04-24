-- +goose Up

-- The answers.question_id FK was created without an ON DELETE action, so it
-- defaulted to RESTRICT.  This means deleting a question bank (which cascades
-- to its questions) is blocked whenever any answer row references one of those
-- questions — resulting in a 500 on DELETE /banks/:id.
--
-- Answer rows are historical game records; they're not meaningful without the
-- source question, so CASCADE is the correct behavior here.
ALTER TABLE answers DROP CONSTRAINT IF EXISTS answers_question_id_fkey;
ALTER TABLE answers
    ADD CONSTRAINT answers_question_id_fkey
        FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE;

-- +goose Down

ALTER TABLE answers DROP CONSTRAINT IF EXISTS answers_question_id_fkey;
ALTER TABLE answers
    ADD CONSTRAINT answers_question_id_fkey
        FOREIGN KEY (question_id) REFERENCES questions(id);
