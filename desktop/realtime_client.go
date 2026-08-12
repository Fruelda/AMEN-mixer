package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ============================================================
// WEBSOCKET CONNECTION
// ============================================================

func (s *RealtimeServer) HandleWebSocket(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println(
		"[WS] Incoming connection",
	)

	// ========================================================
	// UPGRADE HTTP -> WEBSOCKET
	// ========================================================

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

	// ========================================================
	// CREATE CLIENT
	// ========================================================

	client :=
		&WSClient{
			conn: conn,
			Type: "unknown",
		}

	// ========================================================
	// REGISTER SOCKET
	// ========================================================

	s.mu.Lock()

	s.clients[client] =
		true

	total :=
		len(
			s.clients,
		)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client connected: %d\n",
		total,
	)

	// ========================================================
	// CLEANUP
	// ========================================================

	defer func() {

		s.removeClient(
			client,
		)

		_ =
			conn.Close()

	}()

	// ========================================================
	// READ LOOP
	// ========================================================

	for {

		messageType,
			data,
			err :=
			conn.ReadMessage()

		if err != nil {

			log.Println(
				"[WS] Disconnect:",
				err,
			)

			break
		}

		// Hanya proses text WebSocket.
		if messageType !=
			websocket.TextMessage {

			continue
		}

		// ====================================================
		// LOG MESSAGE
		// ====================================================

		logRealtimeJSON(
			"[WS] Received:",
			data,
		)

		// ====================================================
		// ROUTE MESSAGE
		// ====================================================

		s.handleMessage(
			client,
			data,
		)
	}
}
