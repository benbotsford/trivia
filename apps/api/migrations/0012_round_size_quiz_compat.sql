-- +goose Up

-- round_size = 0 is now the sentinel for "quiz-based game" (no fixed round size).
-- Migration 0009 added a CHECK (round_size >= 1) which blocks quiz-based game creation
-- because the handler intentionally passes 0 for quiz_id-driven games.
ALTER TABLE games DROP CONSTRAINT IF EXISTS games_round_size_positive;
ALTER TABLE games ADD CONSTRAINT games_round_size_nonneg CHECK (round_size >= 0);

-- +goose Down
ALTER TABLE games DROP CONSTRAINT IF EXISTS games_round_size_nonneg;
ALTER TABLE games ADD CONSTRAINT games_round_size_positive CHECK (round_size >= 1);
