package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/benbotsford/trivia/internal/store"
)

// --- Bank request/response ---

type createBankRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type updateBankRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type bankResponse struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func bankFromStore(b store.QuestionBank) bankResponse {
	resp := bankResponse{
		ID:        b.ID,
		OwnerID:   b.OwnerID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}
	if b.Description.Valid {
		resp.Description = &b.Description.String
	}
	return resp
}

// nullText converts an optional string pointer to a pgtype.Text.
func nullText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}
