package game

// Handler tests for game CRUD endpoints.
//
// Strategy: for each handler under test, we build a minimal Service whose
// stub fields return canned data, wrap the handler in a chi router (so that
// URL parameters like {gameID} are available via chi.URLParam), and fire
// requests through httptest.  No database or network required.
//
// stubQuerier embeds the querier interface so that any method NOT overridden
// with a function field panics immediately — this makes accidental "extra
// call" bugs loud rather than silent.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/benbotsford/trivia/internal/auth"
	"github.com/benbotsford/trivia/internal/billing"
	"github.com/benbotsford/trivia/internal/store"
)

// ---- stubs -------------------------------------------------------------------

// stubQuerier implements the querier interface with per-method function fields.
// Embed the interface so unimplemented methods panic rather than silently no-op.
type stubQuerier struct {
	querier // embedded interface — nil, panics on any unimplemented method call

	listGamesByHostFn       func(ctx context.Context, arg store.ListGamesByHostParams) ([]store.Game, error)
	getGameByIDFn           func(ctx context.Context, id uuid.UUID) (store.Game, error)
	cancelGameFn            func(ctx context.Context, id uuid.UUID) (store.Game, error)
	listActivePlayersInGame func(ctx context.Context, gameID uuid.UUID) ([]store.GamePlayer, error)
	listAnswersForGameFn    func(ctx context.Context, gameID uuid.UUID) ([]store.Answer, error)
}

func (s *stubQuerier) ListGamesByHost(ctx context.Context, arg store.ListGamesByHostParams) ([]store.Game, error) {
	return s.listGamesByHostFn(ctx, arg)
}
func (s *stubQuerier) GetGameByID(ctx context.Context, id uuid.UUID) (store.Game, error) {
	return s.getGameByIDFn(ctx, id)
}
func (s *stubQuerier) CancelGame(ctx context.Context, id uuid.UUID) (store.Game, error) {
	return s.cancelGameFn(ctx, id)
}
func (s *stubQuerier) ListActivePlayersInGame(ctx context.Context, gameID uuid.UUID) ([]store.GamePlayer, error) {
	return s.listActivePlayersInGame(ctx, gameID)
}
func (s *stubQuerier) ListAnswersForGame(ctx context.Context, gameID uuid.UUID) ([]store.Answer, error) {
	return s.listAnswersForGameFn(ctx, gameID)
}

// stubUserResolver always returns the same fixed user.
type stubUserResolver struct {
	user store.User
	err  error
}

func (s *stubUserResolver) GetOrCreate(_ context.Context, _, _ string) (store.User, error) {
	return s.user, s.err
}

// ---- helpers -----------------------------------------------------------------

// newTestService builds a Service wired to the provided stubs.
// hub and entitlements are set to nil/NoopChecker — sufficient for the read
// handlers tested here.
func newTestService(q querier, u store.User) *Service {
	return &Service{
		q:            q,
		users:        &stubUserResolver{user: u},
		entitlements: billing.NoopChecker{},
		hub:          nil,
	}
}

// withAuth injects auth.Claims into the request context so that handlers can
// call currentUser() without a real JWT or middleware in the chain.
func withAuth(r *http.Request, sub, email string) *http.Request {
	ctx := auth.ContextWithClaims(r.Context(), auth.Claims{Sub: sub, Email: email})
	return r.WithContext(ctx)
}

// chiRequest builds an *http.Request routed through a minimal chi router that
// mounts the given pattern and handler. This populates chi.URLParam so that
// handlers can read path variables like {gameID}.
func chiRequest(method, pattern, url string, body io.Reader, handler http.HandlerFunc) *http.Request {
	r := httptest.NewRequest(method, url, body)
	rr := chi.NewRouter()
	rr.Method(method, pattern, handler)

	// Reuse chi's route context by running the router's Match against the request.
	// This populates chi.RouteContext so URLParam works in the handler.
	rctx := chi.NewRouteContext()
	rr.Match(rctx, method, url)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// decode unmarshals the response body into v.
func decode(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response body: %v\nbody: %s", err, rec.Body.String())
	}
}

// fixture data
var (
	hostID  = uuid.New()
	hostSub = "auth0|testhostid"

	testUser = store.User{
		ID:       hostID,
		Auth0Sub: hostSub,
	}

	gameID   = uuid.New()
	otherID  = uuid.New() // a different host — used for ownership checks

	testGame = store.Game{
		ID:        gameID,
		Code:      "TST001",
		Status:    store.GameStatusCompleted,
		HostID:    hostID,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
)

// ---- listGames ---------------------------------------------------------------

func TestListGames(t *testing.T) {
	t.Run("returns 200 with games for the authenticated host", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			listGamesByHostFn: func(_ context.Context, arg store.ListGamesByHostParams) ([]store.Game, error) {
				if arg.HostID != hostID {
					t.Errorf("ListGamesByHost called with wrong HostID: %v", arg.HostID)
				}
				return []store.Game{testGame}, nil
			},
		}, testUser)

		r := withAuth(httptest.NewRequest(http.MethodGet, "/games", nil), hostSub, "test@example.com")
		w := httptest.NewRecorder()
		svc.listGames(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}
		var games []gameResponse
		decode(t, w, &games)
		if len(games) != 1 {
			t.Errorf("expected 1 game, got %d", len(games))
		}
		if games[0].Code != "TST001" {
			t.Errorf("Code mismatch: got %q, want %q", games[0].Code, "TST001")
		}
	})

	t.Run("returns 200 with empty array when host has no games", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			listGamesByHostFn: func(_ context.Context, _ store.ListGamesByHostParams) ([]store.Game, error) {
				return []store.Game{}, nil
			},
		}, testUser)

		r := withAuth(httptest.NewRequest(http.MethodGet, "/games", nil), hostSub, "")
		w := httptest.NewRecorder()
		svc.listGames(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}
		var games []gameResponse
		decode(t, w, &games)
		if len(games) != 0 {
			t.Errorf("expected empty game list, got %d games", len(games))
		}
	})

	t.Run("returns 401 when no auth claims in context", func(t *testing.T) {
		svc := newTestService(&stubQuerier{}, testUser)

		r := httptest.NewRequest(http.MethodGet, "/games", nil) // no auth
		w := httptest.NewRecorder()
		svc.listGames(w, r)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

// ---- getGame -----------------------------------------------------------------

func TestGetGame(t *testing.T) {
	t.Run("returns 200 with game data", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, id uuid.UUID) (store.Game, error) {
				if id != gameID {
					t.Errorf("GetGameByID called with wrong id: %v", id)
				}
				return testGame, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}", fmt.Sprintf("/games/%s", gameID), nil, svc.getGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.getGame(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}
		var g gameResponse
		decode(t, w, &g)
		if g.ID != gameID {
			t.Errorf("ID mismatch: got %v, want %v", g.ID, gameID)
		}
		if g.Code != "TST001" {
			t.Errorf("Code mismatch: got %q, want %q", g.Code, "TST001")
		}
	})

	t.Run("returns 404 when game does not exist", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return store.Game{}, errors.New("not found")
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}", fmt.Sprintf("/games/%s", gameID), nil, svc.getGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.getGame(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("returns 403 when game belongs to a different host", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				g := testGame
				g.HostID = otherID // belongs to someone else
				return g, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}", fmt.Sprintf("/games/%s", gameID), nil, svc.getGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.getGame(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("returns 400 on malformed game ID", func(t *testing.T) {
		svc := newTestService(&stubQuerier{}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}", "/games/not-a-uuid", nil, svc.getGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.getGame(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

// ---- cancelGame --------------------------------------------------------------

func TestCancelGame(t *testing.T) {
	t.Run("returns 204 on successful cancellation", func(t *testing.T) {
		cancelCalled := false
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return testGame, nil
			},
			cancelGameFn: func(_ context.Context, id uuid.UUID) (store.Game, error) {
				cancelCalled = true
				if id != gameID {
					t.Errorf("CancelGame called with wrong id: %v", id)
				}
				return testGame, nil
			},
		}, testUser)

		r := chiRequest(http.MethodDelete, "/games/{gameID}", fmt.Sprintf("/games/%s", gameID), nil, svc.cancelGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.cancelGame(w, r)

		if w.Code != http.StatusNoContent {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNoContent)
		}
		if !cancelCalled {
			t.Error("expected CancelGame to be called, but it wasn't")
		}
	})

	t.Run("returns 403 when game belongs to a different host", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				g := testGame
				g.HostID = otherID
				return g, nil
			},
		}, testUser)

		r := chiRequest(http.MethodDelete, "/games/{gameID}", fmt.Sprintf("/games/%s", gameID), nil, svc.cancelGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.cancelGame(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("returns 409 when game is already completed or cancelled", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return testGame, nil
			},
			cancelGameFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				// CancelGame only matches lobby/in_progress — returning an error
				// simulates the "no rows updated" case.
				return store.Game{}, errors.New("no rows affected")
			},
		}, testUser)

		r := chiRequest(http.MethodDelete, "/games/{gameID}", fmt.Sprintf("/games/%s", gameID), nil, svc.cancelGame)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.cancelGame(w, r)

		if w.Code != http.StatusConflict {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusConflict)
		}
	})
}

// ---- listPlayers -------------------------------------------------------------

func TestListPlayers(t *testing.T) {
	players := []store.GamePlayer{
		{ID: uuid.New(), GameID: gameID, DisplayName: "Alice", Score: 3000},
		{ID: uuid.New(), GameID: gameID, DisplayName: "Bob", Score: 1500},
	}

	t.Run("returns 200 with player list", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return testGame, nil
			},
			listActivePlayersInGame: func(_ context.Context, _ uuid.UUID) ([]store.GamePlayer, error) {
				return players, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}/players", fmt.Sprintf("/games/%s/players", gameID), nil, svc.listPlayers)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.listPlayers(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}
		var got []map[string]any
		decode(t, w, &got)
		if len(got) != 2 {
			t.Errorf("expected 2 players, got %d", len(got))
		}
	})

	t.Run("returns 403 when game belongs to a different host", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				g := testGame
				g.HostID = otherID
				return g, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}/players", fmt.Sprintf("/games/%s/players", gameID), nil, svc.listPlayers)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.listPlayers(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusForbidden)
		}
	})
}

// ---- gameResults -------------------------------------------------------------

func TestGameResults(t *testing.T) {
	playerID := uuid.New()
	questionID := uuid.New()

	players := []store.GamePlayer{
		{ID: playerID, GameID: gameID, DisplayName: "Alice", Score: 2000},
	}
	answers := []store.Answer{
		{
			ID:            uuid.New(),
			GameID:        gameID,
			QuestionID:    questionID,
			PlayerID:      playerID,
			Answer:        "paris",
			IsCorrect:     true,
			PointsAwarded: 1000,
			SubmittedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	t.Run("returns 200 with per-player results", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return testGame, nil
			},
			listActivePlayersInGame: func(_ context.Context, _ uuid.UUID) ([]store.GamePlayer, error) {
				return players, nil
			},
			listAnswersForGameFn: func(_ context.Context, _ uuid.UUID) ([]store.Answer, error) {
				return answers, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}/results", fmt.Sprintf("/games/%s/results", gameID), nil, svc.gameResults)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.gameResults(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d — body: %s", w.Code, http.StatusOK, w.Body.String())
		}
		var got gameResultsResponse
		decode(t, w, &got)

		if got.Code != "TST001" {
			t.Errorf("Code mismatch: got %q, want %q", got.Code, "TST001")
		}
		if len(got.Players) != 1 {
			t.Fatalf("expected 1 player, got %d", len(got.Players))
		}
		p := got.Players[0]
		if p.DisplayName != "Alice" {
			t.Errorf("DisplayName mismatch: got %q, want %q", p.DisplayName, "Alice")
		}
		if p.TotalScore != 2000 {
			t.Errorf("TotalScore mismatch: got %d, want %d", p.TotalScore, 2000)
		}
		if len(p.Answers) != 1 {
			t.Fatalf("expected 1 answer, got %d", len(p.Answers))
		}
		if !p.Answers[0].IsCorrect {
			t.Errorf("expected answer to be correct")
		}
		if p.Answers[0].PointsAwarded != 1000 {
			t.Errorf("PointsAwarded mismatch: got %d, want %d", p.Answers[0].PointsAwarded, 1000)
		}
	})

	t.Run("returns empty players array when no answers recorded", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return testGame, nil
			},
			listActivePlayersInGame: func(_ context.Context, _ uuid.UUID) ([]store.GamePlayer, error) {
				return []store.GamePlayer{}, nil
			},
			listAnswersForGameFn: func(_ context.Context, _ uuid.UUID) ([]store.Answer, error) {
				return []store.Answer{}, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}/results", fmt.Sprintf("/games/%s/results", gameID), nil, svc.gameResults)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.gameResults(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
		}
		var got gameResultsResponse
		decode(t, w, &got)
		if len(got.Players) != 0 {
			t.Errorf("expected empty players, got %d", len(got.Players))
		}
	})

	t.Run("returns 403 when game belongs to a different host", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				g := testGame
				g.HostID = otherID
				return g, nil
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}/results", fmt.Sprintf("/games/%s/results", gameID), nil, svc.gameResults)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.gameResults(w, r)

		if w.Code != http.StatusForbidden {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("returns 404 when game does not exist", func(t *testing.T) {
		svc := newTestService(&stubQuerier{
			getGameByIDFn: func(_ context.Context, _ uuid.UUID) (store.Game, error) {
				return store.Game{}, errors.New("not found")
			},
		}, testUser)

		r := chiRequest(http.MethodGet, "/games/{gameID}/results", fmt.Sprintf("/games/%s/results", gameID), nil, svc.gameResults)
		r = withAuth(r, hostSub, "")
		w := httptest.NewRecorder()
		svc.gameResults(w, r)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
