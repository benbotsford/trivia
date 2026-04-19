// Package realtime manages WebSocket connections and in-process room state.
//
// Architecture (single-replica phase):
//   - One Hub per process, holding all active rooms keyed by game code.
//   - Scale-out path: replace the in-process broadcast with Redis Pub/Sub fan-out.
//     All call sites (Broadcast, JoinRoom, LeaveRoom) stay the same.
package realtime

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-chi/chi/v5"
)

// MessageType identifies the kind of realtime event being sent or received.
type MessageType string

const (
	// Server → client
	MsgQuestionRevealed  MessageType = "question_revealed"
	MsgTimerTick         MessageType = "timer_tick"
	MsgAnswerAccepted    MessageType = "answer_accepted"
	MsgScoreboardUpdate  MessageType = "scoreboard_update"
	MsgGameEnded         MessageType = "game_ended"

	// Client → server (host)
	MsgStartGame    MessageType = "start_game"
	MsgAdvanceQuestion MessageType = "advance_question"
	MsgEndGame      MessageType = "end_game"

	// Client → server (player)
	MsgJoin         MessageType = "join"
	MsgSubmitAnswer MessageType = "submit_answer"
)

// Message is the wire format for all WebSocket messages.
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// client represents a single connected WebSocket peer.
type client struct {
	conn *websocket.Conn
	send chan Message
}

// room holds all clients connected to a single game.
type room struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

func newRoom() *room {
	return &room{clients: make(map[*client]struct{})}
}

func (r *room) add(c *client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (r *room) remove(c *client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

func (r *room) broadcast(msg Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for c := range r.clients {
		select {
		case c.send <- msg:
		default:
			// Slow client — drop the message rather than blocking the broadcast.
			slog.Warn("dropped message to slow client")
		}
	}
}

// Hub owns all active rooms and the HTTP handler for WebSocket upgrades.
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*room // keyed by game code
}

// New creates a Hub.
func New() *Hub {
	return &Hub{rooms: make(map[string]*room)}
}

// RegisterRoutes mounts the WebSocket upgrade endpoint on r.
// Expected URL: /ws/{gameCode}
func (h *Hub) RegisterRoutes(r chi.Router) {
	r.Get("/ws/{gameCode}", h.HandleWebSocket)
}

// HandleWebSocket upgrades an HTTP connection to WebSocket and registers the
// client in the appropriate room. It blocks until the connection closes.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	gameCode := chi.URLParam(r, "gameCode")
	if gameCode == "" {
		http.Error(w, "missing game code", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// TODO: tighten OriginPatterns before going to production.
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		slog.Error("websocket accept failed", "err", err)
		return
	}

	c := &client{conn: conn, send: make(chan Message, 64)}
	rm := h.getOrCreateRoom(gameCode)
	rm.add(c)
	defer func() {
		rm.remove(c)
		conn.Close(websocket.StatusNormalClosure, "bye")
	}()

	ctx := r.Context()

	// Writer goroutine — drains c.send and writes to the wire.
	go func() {
		for msg := range c.send {
			if err := wsjson.Write(ctx, conn, msg); err != nil {
				slog.Debug("ws write error", "err", err)
				return
			}
		}
	}()

	// Reader loop — blocks until the client disconnects.
	h.readLoop(ctx, c, gameCode)
}

func (h *Hub) readLoop(ctx context.Context, c *client, gameCode string) {
	for {
		var msg Message
		if err := wsjson.Read(ctx, c.conn, &msg); err != nil {
			slog.Debug("ws read closed", "gameCode", gameCode, "err", err)
			return
		}
		h.handleMessage(ctx, c, gameCode, msg)
	}
}

// handleMessage dispatches an inbound message to the appropriate handler.
// TODO: wire up game service calls here.
func (h *Hub) handleMessage(_ context.Context, _ *client, gameCode string, msg Message) {
	slog.Debug("ws message received", "gameCode", gameCode, "type", msg.Type)
	switch msg.Type {
	case MsgJoin:
		// TODO: validate display name, register player in DB, broadcast roster update
	case MsgSubmitAnswer:
		// TODO: validate answer window, persist to DB, ack to sender
	case MsgStartGame:
		// TODO: check host claim, update game status, broadcast first question
	case MsgAdvanceQuestion:
		// TODO: increment question index, broadcast next question
	case MsgEndGame:
		// TODO: flush scores to Postgres, broadcast game_ended, clean up Redis
	default:
		slog.Warn("unknown message type", "type", msg.Type)
	}
}

// Broadcast sends a message to every client in the named room.
func (h *Hub) Broadcast(gameCode string, msg Message) {
	h.mu.RLock()
	rm, ok := h.rooms[gameCode]
	h.mu.RUnlock()
	if !ok {
		return
	}
	rm.broadcast(msg)
}

func (h *Hub) getOrCreateRoom(code string) *room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if r, ok := h.rooms[code]; ok {
		return r
	}
	r := newRoom()
	h.rooms[code] = r
	return r
}
