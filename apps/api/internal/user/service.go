// Package user manages host accounts.
package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/benbotsford/trivia/internal/store"
)

// Service provides operations on host user accounts.
type Service struct {
	q *store.Queries
}

// New creates a Service backed by the given sqlc Queries.
func New(q *store.Queries) *Service {
	return &Service{q: q}
}

// GetOrCreate returns the user matching auth0Sub, creating one if needed.
// email is optional — pass an empty string if not available from the token.
func (s *Service) GetOrCreate(ctx context.Context, auth0Sub, email string) (store.User, error) {
	u, err := s.q.GetUserByAuth0Sub(ctx, auth0Sub)
	if err == nil {
		return u, nil
	}

	slog.Info("creating new user", "auth0_sub", auth0Sub)
	return s.q.CreateUser(ctx, store.CreateUserParams{
		ID:       uuid.New(),
		Auth0Sub: auth0Sub,
		Email:    pgtype.Text{String: email, Valid: email != ""},
	})
}
