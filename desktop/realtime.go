package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================================
// MESSAGE
// ============================================================

type RealtimeMessage struct {
	Type    string      `json:"type"`
	Channel interface{} `json:"channel,omitempty"`
}

// ============================================================
// CLIENT
// ============================================================

type WSClient struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// ============================================================
// SERVER
// ============================================================

type RealtimeServer struct {
	clients map[*WSClient]bool
	mu      sync.Mutex
}

// ============================================================
// WEBSOCKET CONFIG
// ============================================================

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Development LAN.
		// Nanti kalau production sebaiknya dibatasi.
		return true
	},
}

// ============================================================
// NEW SERVER
// ============================================================

func NewRealtimeServer() *RealtimeServer {
	return &RealtimeServer{
		clients: make(map[*WSClient]bool),
	}
}

// ============================================================
// WEBSOCKET HANDLER
// ============================================================

func (s *RealtimeServer) HandleWebSocket(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Println("===================================")
	log.Println("[WS] Incoming request")
	log.Println("[WS] Method:", r.Method)
	log.Println("[WS] URL:", r.URL.String())
	log.Println("[WS] Remote:", r.RemoteAddr)
	log.Println("[WS] Connection:", r.Header.Get("Connection"))
	log.Println("[WS] Upgrade:", r.Header.Get("Upgrade"))
	log.Println("===================================")

	// --------------------------------------------------------
	// UPGRADE HTTP -> WEBSOCKET
	// --------------------------------------------------------

	conn, err := upgrader.Upgrade(
		w,
		r,
		nil,
	)

	if err != nil {
		log.Println(
			"[WS] Upgrade error:",
			err,
		)

		return
	}

	client := &WSClient{
		conn: conn,
	}

	// --------------------------------------------------------
	// REGISTER CLIENT
	// --------------------------------------------------------

	s.mu.Lock()

	s.clients[client] = true

	clientCount := len(s.clients)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client connected. Total: %d\n",
		clientCount,
	)

	// --------------------------------------------------------
	// CLEANUP
	// --------------------------------------------------------

	defer func() {

		s.removeClient(client)

		_ = conn.Close()

	}()

	// --------------------------------------------------------
	// READ LOOP
	// --------------------------------------------------------

	for {

		messageType, message, err :=
			conn.ReadMessage()

		if err != nil {

			log.Printf(
				"[WS] Client disconnected: %v\n",
				err,
			)

			break
		}

		// ----------------------------------------------------
		// ONLY ACCEPT TEXT
		// ----------------------------------------------------

		if messageType != websocket.TextMessage {
			continue
		}

		log.Println(
			"[WS] Received:",
			string(message),
		)

		// ----------------------------------------------------
		// BROADCAST TO ALL CLIENTS
		// ----------------------------------------------------

		s.broadcast(message)
	}
}

// ============================================================
// REMOVE CLIENT
// ============================================================

func (s *RealtimeServer) removeClient(
	client *WSClient,
) {
	s.mu.Lock()

	_, exists := s.clients[client]

	if exists {
		delete(
			s.clients,
			client,
		)
	}

	clientCount := len(s.clients)

	s.mu.Unlock()

	if exists {

		log.Printf(
			"[WS] Client disconnected. Total: %d\n",
			clientCount,
		)
	}
}

// ============================================================
// BROADCAST RAW MESSAGE
// ============================================================

func (s *RealtimeServer) broadcast(
	message []byte,
) {

	log.Println(
		"[WS] BROADCAST TO CLIENTS:",
		string(message),
	)

	// ========================================================
	// COPY CLIENT LIST
	// ========================================================

	s.mu.Lock()

	clients := make(
		[]*WSClient,
		0,
		len(s.clients),
	)

	for client := range s.clients {
		clients = append(
			clients,
			client,
		)
	}

	s.mu.Unlock()

	log.Printf(
		"[WS] Broadcasting to %d clients\n",
		len(clients),
	)

	// ========================================================
	// SEND TO ALL CLIENTS
	// ========================================================

	for _, client := range clients {

		client.mu.Lock()

		err := client.conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		client.mu.Unlock()

		if err != nil {

			log.Println(
				"[WS] Write error:",
				err,
			)

			s.removeClient(client)

			_ = client.conn.Close()

			continue
		}

		log.Println(
			"[WS] Message sent successfully",
		)
	}
}

// ============================================================
// BROADCAST JSON
// ============================================================

func (s *RealtimeServer) BroadcastJSON(
	message RealtimeMessage,
) {
	data, err :=
		json.Marshal(message)

	if err != nil {

		log.Println(
			"[WS] JSON marshal error:",
			err,
		)

		return
	}

	log.Println(
		"[WS] Broadcasting:",
		string(data),
	)

	s.broadcast(data)
}

// ============================================================
// CLIENT COUNT
// ============================================================

func (s *RealtimeServer) ClientCount() int {

	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.clients)
}

// ============================================================
// START SERVER
// ============================================================

func StartRealtimeServer(
	server *RealtimeServer,
) {

	mux := http.NewServeMux()

	// --------------------------------------------------------
	// WEBSOCKET ENDPOINT
	// --------------------------------------------------------

	mux.HandleFunc(
		"/ws",
		server.HandleWebSocket,
	)

	// --------------------------------------------------------
	// HEALTH CHECK
	// --------------------------------------------------------
	//
	// Buka:
	//
	// http://IP-MAC:8081/
	//
	// Kalau keluar "AMEN MIXER REALTIME SERVER OK"
	// berarti port 8081 memang kebuka.
	//

	mux.HandleFunc(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			w.Header().Set(
				"Content-Type",
				"text/plain",
			)

			w.WriteHeader(
				http.StatusOK,
			)

			_, _ = w.Write(
				[]byte(
					"AMEN MIXER REALTIME SERVER OK\n",
				),
			)
		},
	)

	// --------------------------------------------------------
	// START
	// --------------------------------------------------------

	go func() {

		log.Println(
			"[WS] ===================================",
		)

		log.Println(
			"[WS] REALTIME SERVER",
		)

		log.Println(
			"[WS] Listening on 0.0.0.0:8081",
		)

		log.Println(
			"[WS] WebSocket: ws://0.0.0.0:8081/ws",
		)

		log.Println(
			"[WS] ===================================",
		)

		server := &http.Server{
			Addr:              "0.0.0.0:8081",
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		err := server.ListenAndServe()

		if err != nil &&
			err != http.ErrServerClosed {

			log.Println(
				"[WS] Server stopped:",
				err,
			)
		}

	}()

}
