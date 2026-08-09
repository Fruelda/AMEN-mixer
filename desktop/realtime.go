package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"desktop/backend/protocol"

	"github.com/gorilla/websocket"
)

// ============================================================
// CLIENT
// ============================================================

type WSClient struct {
	conn *websocket.Conn

	mu sync.Mutex
}

// ============================================================
// SERVER
// ============================================================

type RealtimeServer struct {
	clients map[*WSClient]bool

	mu sync.Mutex
}

// ============================================================
// WEBSOCKET CONFIG
// ============================================================

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {

		// Development mode.
		// Production sebaiknya dibatasi.

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

	log.Println("[WS] Incoming connection")

	conn, err :=
		upgrader.Upgrade(
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

	client :=
		&WSClient{

			conn: conn,
		}

	// Register

	s.mu.Lock()

	s.clients[client] = true

	total :=
		len(s.clients)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client connected: %d\n",
		total,
	)

	defer func() {

		s.removeClient(
			client,
		)

		_ =
			conn.Close()

	}()

	for {

		messageType, data, err :=
			conn.ReadMessage()

		if err != nil {

			log.Println(
				"[WS] Disconnect:",
				err,
			)

			break

		}

		if messageType != websocket.TextMessage {

			continue

		}

		log.Println(
			"[WS] Received:",
			string(data),
		)

		/*

			Sementara relay.

			Nanti COMMAND akan masuk
			ke handler backend.

		*/

		s.broadcast(
			data,
		)

	}

}

// ============================================================
// REMOVE CLIENT
// ============================================================

func (s *RealtimeServer) removeClient(
	client *WSClient,
) {

	s.mu.Lock()

	delete(
		s.clients,
		client,
	)

	total :=
		len(s.clients)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client removed: %d\n",
		total,
	)

}

// ============================================================
// BROADCAST RAW
// ============================================================

func (s *RealtimeServer) broadcast(
	data []byte,
) {

	s.mu.Lock()

	clients :=
		make(
			[]*WSClient,
			0,
			len(s.clients),
		)

	for client := range s.clients {

		clients =
			append(
				clients,
				client,
			)

	}

	s.mu.Unlock()

	for _, client := range clients {

		client.mu.Lock()

		err :=
			client.conn.WriteMessage(
				websocket.TextMessage,
				data,
			)

		client.mu.Unlock()

		if err != nil {

			log.Println(
				"[WS] Write error:",
				err,
			)

			s.removeClient(
				client,
			)

		}

	}

}

// ============================================================
// BROADCAST JSON
// ============================================================

func (s *RealtimeServer) BroadcastJSON(
	message protocol.RealtimeMessage,
) {

	data, err :=
		json.Marshal(
			message,
		)

	if err != nil {

		log.Println(
			"[WS] Marshal error:",
			err,
		)

		return

	}

	log.Println(
		"[WS] Broadcast:",
		string(data),
	)

	s.broadcast(
		data,
	)

}

// ============================================================
// CHANNEL UPDATE HELPER
// ============================================================

func (s *RealtimeServer) BroadcastChannelUpdate(
	id int,
	volume int,
	muted bool,
) {

	message :=
		protocol.RealtimeMessage{

			Type: protocol.MessageChannelUpdate,

			Channel: &protocol.ChannelUpdate{

				ID: id,

				Volume: &volume,

				Muted: &muted,
			},
		}

	s.BroadcastJSON(
		message,
	)

}

// ============================================================
// DEVICE STATUS
// ============================================================

func (s *RealtimeServer) BroadcastDeviceStatus(
	connected bool,
) {

	message :=
		protocol.RealtimeMessage{

			Type: protocol.MessageDeviceStatus,

			Connected: &connected,
		}

	s.BroadcastJSON(
		message,
	)

}

// ============================================================
// CLIENT COUNT
// ============================================================

func (s *RealtimeServer) ClientCount() int {

	s.mu.Lock()

	defer s.mu.Unlock()

	return len(
		s.clients,
	)

}

// ============================================================
// START SERVER
// ============================================================

func StartRealtimeServer(
	server *RealtimeServer,
) {

	mux :=
		http.NewServeMux()

	mux.HandleFunc(
		"/ws",
		server.HandleWebSocket,
	)

	mux.HandleFunc(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			w.Header().
				Set(
					"Content-Type",
					"text/plain",
				)

			_, _ =
				w.Write(
					[]byte(
						"AMEN MIXER REALTIME SERVER OK\n",
					),
				)

		},
	)

	go func() {

		log.Println(
			"[WS] ===================================",
		)

		log.Println(
			"[WS] AMEN REALTIME SERVER",
		)

		log.Println(
			"[WS] Listening :8081",
		)

		log.Println(
			"[WS] WS endpoint ws://0.0.0.0:8081/ws",
		)

		serverHTTP :=
			&http.Server{

				Addr: "0.0.0.0:8081",

				Handler: mux,

				ReadHeaderTimeout: 5 * time.Second,
			}

		err :=
			serverHTTP.ListenAndServe()

		if err != nil &&
			err != http.ErrServerClosed {

			log.Println(
				"[WS] Server error:",
				err,
			)

		}

	}()

}
