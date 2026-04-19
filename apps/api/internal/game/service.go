// Package game manages quiz question banks and game lifecycle.
package game

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/btbots1994/trivia/internal/billing"
	"github.com/btbots1994/trivia/internal/store"
)

const codeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // unambiguous chars

// Service handles question bank and game CRUD.
type Service struct {
	q            *store.Queries
	entitlements billing.EntitlementChecker
}

// New creates a Service.
func New(q *store.Queries, ent billing.EntitlementChecker) *Service {
	return &Service{q: q, entitlements: ent}
}

// RegisterRoutes mounts all game-related routes onto r.
// All routes under /games and /banks require an authenticated host (auth
// middleware must already be applied on the parent router).
func (s *Service) RegisterRoutes(r chi.Router) {
	r.Route("/banks", func(r chi.Router) {
		r.Get("/", s.listBanks)
		r.Post("/", s.createBank)
		r.Route("/{bankID}", func(r chi.Router) {
			r.Get("/", s.getBank)
			r.Delete("/", s.deleteBank)
			r.Get("/questions", s.listQuestions)
			r.Post("/questions", s.createQuestion)
		})
	})

	r.Route("/games", func(r chi.Router) {
		r.Post("/", s.createGame)
		r.Get("/", s.listGames)
		r.Route("/{gameID}", func(r chi.Router) {
			r.Get("/", s.getGame)
			r.Post("/start", s.startGame)
			r.Post("/next", s.nextQuestion)
			r.Post("/end", s.endGame)
		})
	})
}

// --- Question Banks ---

func (s *Service) listBanks(w http.ResponseWriter, r *http.Request) {
	// TODO: extract userID from auth claims, query s.q.ListQuestionBanksByOwner
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) createBank(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) getBank(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) deleteBank(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) listQuestions(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) createQuestion(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// --- Games ---

func (s *Service) createGame(w http.ResponseWriter, r *http.Request) {
	// TODO: parse body, check entitlements, create game row
	_ = chi.URLParam(r, "gameID")
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) listGames(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) getGame(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) startGame(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) nextQuestion(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) endGame(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// --- Helpers ---

// generateCode returns a random 6-character game code.
func generateCode(ctx context.Context) (string, error) {
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeChars))))
		if err != nil {
			return "", fmt.Errorf("generate code: %w", err)
		}
		code[i] = codeChars[n.Int64()]
	}
	slog.DebugContext(ctx, "generated game code", "code", string(code))
	return string(code), nil
}

// mustParseUUID is a helper that parses a Chi URL param as a UUID.
func mustParseUUID(r *http.Request, param string) (uuid.UUID, error) {
	raw := chi.URLParam(r, param)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s %q: %w", param, raw, err)
	}
	return id, nil
}
