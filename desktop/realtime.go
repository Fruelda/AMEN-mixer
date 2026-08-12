package main

import (
	"net/http"
	"sync"

	"desktop/backend/protocol"

	"github.com/gorilla/websocket"
)

// ============================================================
// CLIENT
// ============================================================

type WSClient struct {
	conn *websocket.Conn
	mu   sync.Mutex

	ID   string
	Name string
	Type string
}

// ============================================================
// REALTIME SERVER
// ============================================================

type RealtimeServer struct {
	clients map[*WSClient]bool

	channels map[int]protocol.Channel

	mu sync.Mutex

	// Handler untuk update yang datang dari
	// HP / browser / client remote.
	//
	// Handler ini nanti dihubungkan ke audio.Manager
	// melalui app.go.
	onChannelUpdate func(
		protocol.ChannelUpdate,
	) error
}

// ============================================================
// WEBSOCKET CONFIG
// ============================================================

var upgrader = websocket.Upgrader{
	CheckOrigin: func(
		r *http.Request,
	) bool {
		return true
	},
}

// ============================================================
// CREATE SERVER
// ============================================================

func NewRealtimeServer() *RealtimeServer {

	return &RealtimeServer{
		clients: make(
			map[*WSClient]bool,
		),

		channels: make(
			map[int]protocol.Channel,
		),
	}
}

// ============================================================
// CHANNEL UPDATE HANDLER
// ============================================================

func (s *RealtimeServer) SetChannelUpdateHandler(
	handler func(
		protocol.ChannelUpdate,
	) error,
) {

	s.onChannelUpdate =
		handler
}
