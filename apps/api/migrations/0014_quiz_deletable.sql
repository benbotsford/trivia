-- +goose Up

-- Allow quizzes to be deleted even when past games reference them.
--
-- Previous behavior (migration 0013): ON DELETE RESTRICT — quiz deletion was
-- blocked at the DB level if any game row had quiz_id set.
--
-- The reason RESTRICT was chosen over SET NULL was the games_bank_or_quiz
-- CHECK constraint: a quiz-only game (bank_id IS NULL) would violate the
-- check if its quiz_id were nulled out.
--
-- New behavior:
--   • FK changes to ON DELETE SET NULL — quiz deletion nulls quiz_id on all
--     referencing game rows.
--   • CHECK constraint is relaxed to allow both ids to be NULL when the game
--     is already completed or cancelled (historical records no longer need a
--     live quiz). Active games (lobby / in_progress) still require a source.
--   • The application layer no longer blocks deletion based on game count.

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_quiz_id_fkey;
ALTER TABLE games
    ADD CONSTRAINT games_quiz_id_fkey
        FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE SET NULL;

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_bank_or_quiz;
ALTER TABLE games
    ADD CONSTRAINT games_bank_or_quiz CHECK (
        bank_id IS NOT NULL
        OR quiz_id IS NOT NULL
        OR status IN ('completed', 'cancelled')
    );

-- +goose Down

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_quiz_id_fkey;
ALTER TABLE games
    ADD CONSTRAINT games_quiz_id_fkey
        FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE RESTRICT;

ALTER TABLE games DROP CONSTRAINT IF EXISTS games_bank_or_quiz;
ALTER TABLE games
    ADD CONSTRAINT games_bank_or_quiz CHECK (
        bank_id IS NOT NULL OR quiz_id IS NOT NULL
    );
