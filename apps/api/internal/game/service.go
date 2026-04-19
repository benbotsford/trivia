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

	"github.com/benbotsford/trivia/internal/auth"
	"github.com/benbotsford/trivia/internal/billing"
	"github.com/benbotsford/trivia/internal/store"
	"github.com/benbotsford/trivia/internal/user"
)

const codeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // unambiguous chars

// Service handles question bank and game CRUD.
type Service struct {
	q            *store.Queries
	users        *user.Service
	entitlements billing.EntitlementChecker
}

// New creates a Service.
func New(q *store.Queries, users *user.Service, ent billing.EntitlementChecker) *Service {
	return &Service{q: q, users: users, entitlements: ent}
}

// RegisterRoutes mounts all game-related routes onto r.
// All routes require an authenticated host (auth middleware must already be
// applied on the parent router).
func (s *Service) RegisterRoutes(r chi.Router) {
	r.Route("/banks", func(r chi.Router) {
		r.Get("/", s.listBanks)
		r.Post("/", s.createBank)
		r.Route("/{bankID}", func(r chi.Router) {
			r.Get("/", s.getBank)
			r.Put("/", s.updateBank)
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
	u, err := s.currentUser(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	banks, err := s.q.ListQuestionBanksByOwner(r.Context(), u.ID)
	if err != nil {
		slog.ErrorContext(r.Context(), "listBanks: query failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := make([]bankResponse, len(banks))
	for i, b := range banks {
		resp[i] = bankFromStore(b)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Service) createBank(w http.ResponseWriter, r *http.Request) {
	u, err := s.currentUser(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createBankRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}

	bank, err := s.q.CreateQuestionBank(r.Context(), store.CreateQuestionBankParams{
		ID:          uuid.New(),
		OwnerID:     u.ID,
		Name:        req.Name,
		Description: nullText(req.Description),
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "createBank: insert failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, bankFromStore(bank))
}

func (s *Service) getBank(w http.ResponseWriter, r *http.Request) {
	u, err := s.currentUser(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bankID, err := mustParseUUID(r, "bankID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bank, err := s.q.GetQuestionBank(r.Context(), bankID)
	if err != nil {
		writeError(w, http.StatusNotFound, "bank not found")
		return
	}
	if bank.OwnerID != u.ID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	writeJSON(w, http.StatusOK, bankFromStore(bank))
}

func (s *Service) updateBank(w http.ResponseWriter, r *http.Request) {
	u, err := s.currentUser(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bankID, err := mustParseUUID(r, "bankID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bank, err := s.q.GetQuestionBank(r.Context(), bankID)
	if err != nil {
		writeError(w, http.StatusNotFound, "bank not found")
		return
	}
	if bank.OwnerID != u.ID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req updateBankRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusUnprocessableEntity, "name is required")
		return
	}

	updated, err := s.q.UpdateQuestionBank(r.Context(), store.UpdateQuestionBankParams{
		ID:          bankID,
		Name:        req.Name,
		Description: nullText(req.Description),
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "updateBank: update failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, bankFromStore(updated))
}

func (s *Service) deleteBank(w http.ResponseWriter, r *http.Request) {
	u, err := s.currentUser(r.Context())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bankID, err := mustParseUUID(r, "bankID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bank, err := s.q.GetQuestionBank(r.Context(), bankID)
	if err != nil {
		writeError(w, http.StatusNotFound, "bank not found")
		return
	}
	if bank.OwnerID != u.ID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := s.q.DeleteQuestionBank(r.Context(), bankID); err != nil {
		slog.ErrorContext(r.Context(), "deleteBank: delete failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) listQuestions(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Service) createQuestion(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// --- Games ---

func (s *Service) createGame(w http.ResponseWriter, r *http.Request) {
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

// currentUser resolves the authenticated user from the request context,
// creating a DB record on first login if needed.
func (s *Service) currentUser(ctx context.Context) (store.User, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return store.User{}, fmt.Errorf("no auth claims in context")
	}
	return s.users.GetOrCreate(ctx, claims.Sub, claims.Email)
}

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

// mustParseUUID parses a Chi URL param as a UUID.
func mustParseUUID(r *http.Request, param string) (uuid.UUID, error) {
	raw := chi.URLParam(r, param)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s %q: %w", param, raw, err)
	}
	return id, nil
}
