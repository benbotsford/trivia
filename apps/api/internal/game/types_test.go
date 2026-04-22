package game

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/benbotsford/trivia/internal/store"
)

// ts wraps a time.Time in a valid pgtype.Timestamptz, matching what sqlc
// returns when a NOT NULL timestamp column is scanned.
func ts(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// --- nullText -------------------------------------------------------------------

func TestNullText(t *testing.T) {
	t.Run("nil input produces invalid pgtype.Text", func(t *testing.T) {
		got := nullText(nil)
		if got.Valid {
			t.Errorf("expected Valid=false for nil input, got Valid=true")
		}
	})

	t.Run("non-nil input produces valid pgtype.Text with correct string", func(t *testing.T) {
		s := "hello"
		got := nullText(&s)
		if !got.Valid {
			t.Errorf("expected Valid=true for non-nil input, got Valid=false")
		}
		if got.String != s {
			t.Errorf("expected String=%q, got %q", s, got.String)
		}
	})

	t.Run("empty string is valid (not null)", func(t *testing.T) {
		s := ""
		got := nullText(&s)
		if !got.Valid {
			t.Errorf("expected Valid=true for empty string, got Valid=false")
		}
	})
}

// --- bankFromStore --------------------------------------------------------------

func TestBankFromStore(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New()
	ownerID := uuid.New()

	t.Run("without description", func(t *testing.T) {
		bank := store.QuestionBank{
			ID:        id,
			OwnerID:   ownerID,
			Name:      "Pub Quiz Classics",
			CreatedAt: ts(now),
			UpdatedAt: ts(now),
		}
		got := bankFromStore(bank)

		if got.ID != id {
			t.Errorf("ID mismatch: got %v, want %v", got.ID, id)
		}
		if got.Name != "Pub Quiz Classics" {
			t.Errorf("Name mismatch: got %q, want %q", got.Name, "Pub Quiz Classics")
		}
		if got.Description != nil {
			t.Errorf("expected Description=nil, got %v", got.Description)
		}
		if !got.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt mismatch: got %v, want %v", got.CreatedAt, now)
		}
	})

	t.Run("with description", func(t *testing.T) {
		bank := store.QuestionBank{
			ID:          id,
			OwnerID:     ownerID,
			Name:        "Science",
			Description: pgtype.Text{String: "Hard science questions", Valid: true},
			CreatedAt:   ts(now),
			UpdatedAt:   ts(now),
		}
		got := bankFromStore(bank)

		if got.Description == nil {
			t.Fatal("expected non-nil Description")
		}
		if *got.Description != "Hard science questions" {
			t.Errorf("Description mismatch: got %q, want %q", *got.Description, "Hard science questions")
		}
	})
}

// --- questionFromStore ----------------------------------------------------------

func TestQuestionFromStore(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	qID := uuid.New()
	bankID := uuid.New()

	t.Run("text question with accepted_answers", func(t *testing.T) {
		answers := []string{"paris", "Paris", "PARIS"}
		raw, _ := json.Marshal(answers)

		q := store.Question{
			ID:        qID,
			BankID:    bankID,
			Type:      store.QuestionTypeText,
			Prompt:    "What is the capital of France?",
			Points:    1000,
			Position:  0,
			Choices:   raw,
			CreatedAt: ts(now),
			UpdatedAt: ts(now),
		}
		got := questionFromStore(q)

		if got.Type != "text" {
			t.Errorf("Type mismatch: got %q, want %q", got.Type, "text")
		}
		if len(got.AcceptedAnswers) != 3 {
			t.Errorf("expected 3 AcceptedAnswers, got %d", len(got.AcceptedAnswers))
		}
		if got.AcceptedAnswers[0] != "paris" {
			t.Errorf("AcceptedAnswers[0] mismatch: got %q, want %q", got.AcceptedAnswers[0], "paris")
		}
		if got.Choices != nil {
			t.Errorf("expected nil Choices for text question, got %v", got.Choices)
		}
	})

	t.Run("text question falls back to CorrectAnswer when Choices is empty", func(t *testing.T) {
		q := store.Question{
			ID:            qID,
			BankID:        bankID,
			Type:          store.QuestionTypeText,
			Prompt:        "Legacy question",
			Points:        500,
			CorrectAnswer: "legacy answer",
			Choices:       nil,
			CreatedAt:     ts(now),
			UpdatedAt:     ts(now),
		}
		got := questionFromStore(q)

		if len(got.AcceptedAnswers) != 1 {
			t.Errorf("expected 1 AcceptedAnswer from fallback, got %d", len(got.AcceptedAnswers))
		}
		if got.AcceptedAnswers[0] != "legacy answer" {
			t.Errorf("fallback answer mismatch: got %q, want %q", got.AcceptedAnswers[0], "legacy answer")
		}
	})

	t.Run("multiple_choice question", func(t *testing.T) {
		choices := []mcChoice{
			{Text: "Paris", Correct: true},
			{Text: "Berlin", Correct: false},
			{Text: "Madrid", Correct: false},
		}
		raw, _ := json.Marshal(choices)

		q := store.Question{
			ID:        qID,
			BankID:    bankID,
			Type:      store.QuestionTypeMultipleChoice,
			Prompt:    "Capital of France?",
			Points:    500,
			Choices:   raw,
			CreatedAt: ts(now),
			UpdatedAt: ts(now),
		}
		got := questionFromStore(q)

		if got.Type != "multiple_choice" {
			t.Errorf("Type mismatch: got %q, want %q", got.Type, "multiple_choice")
		}
		if len(got.Choices) != 3 {
			t.Errorf("expected 3 Choices, got %d", len(got.Choices))
		}
		if !got.Choices[0].Correct {
			t.Errorf("expected first choice to be correct")
		}
		if got.AcceptedAnswers != nil {
			t.Errorf("expected nil AcceptedAnswers for MC question, got %v", got.AcceptedAnswers)
		}
	})
}

// --- gameFromStore --------------------------------------------------------------

func TestGameFromStore(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	gameID := uuid.New()
	hostID := uuid.New()

	t.Run("game with neither bank nor quiz ID", func(t *testing.T) {
		g := store.Game{
			ID:        gameID,
			Code:      "ABC123",
			Status:    store.GameStatusCompleted,
			HostID:    hostID,
			CreatedAt: ts(now),
			UpdatedAt: ts(now),
		}
		got := gameFromStore(g)

		if got.ID != gameID {
			t.Errorf("ID mismatch: got %v, want %v", got.ID, gameID)
		}
		if got.Code != "ABC123" {
			t.Errorf("Code mismatch: got %q, want %q", got.Code, "ABC123")
		}
		if got.Status != "completed" {
			t.Errorf("Status mismatch: got %q, want %q", got.Status, "completed")
		}
		if got.BankID != nil {
			t.Errorf("expected nil BankID, got %v", got.BankID)
		}
		if got.QuizID != nil {
			t.Errorf("expected nil QuizID, got %v", got.QuizID)
		}
		if !got.CreatedAt.Equal(now) {
			t.Errorf("CreatedAt mismatch: got %v, want %v", got.CreatedAt, now)
		}
	})

	t.Run("game with a bank ID", func(t *testing.T) {
		bankID := uuid.New()
		g := store.Game{
			ID:     gameID,
			Code:   "XYZ789",
			Status: store.GameStatusInProgress,
			HostID: hostID,
			BankID: pgtype.UUID{Bytes: bankID, Valid: true},
			CreatedAt: ts(now),
			UpdatedAt: ts(now),
		}
		got := gameFromStore(g)

		if got.BankID == nil {
			t.Fatal("expected non-nil BankID")
		}
		if *got.BankID != bankID {
			t.Errorf("BankID mismatch: got %v, want %v", *got.BankID, bankID)
		}
		if got.QuizID != nil {
			t.Errorf("expected nil QuizID, got %v", got.QuizID)
		}
	})

	t.Run("game with a quiz ID", func(t *testing.T) {
		quizID := uuid.New()
		g := store.Game{
			ID:     gameID,
			Code:   "QQQ111",
			Status: store.GameStatusLobby,
			HostID: hostID,
			QuizID: pgtype.UUID{Bytes: quizID, Valid: true},
			CreatedAt: ts(now),
			UpdatedAt: ts(now),
		}
		got := gameFromStore(g)

		if got.QuizID == nil {
			t.Fatal("expected non-nil QuizID")
		}
		if *got.QuizID != quizID {
			t.Errorf("QuizID mismatch: got %v, want %v", *got.QuizID, quizID)
		}
		if got.BankID != nil {
			t.Errorf("expected nil BankID, got %v", got.BankID)
		}
	})
}
