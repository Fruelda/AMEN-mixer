package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *RealtimeServer) HandleWebSocket(
	w http.ResponseWriter,
	r *http.Request,
) {
	log.Println("[WS] Incoming connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WS] Upgrade error:", err)
		return
	}

	client := &WSClient{
		conn: conn,
		Type: "unknown",
	}

	s.mu.Lock()
	s.clients[client] = true
	total := len(s.clients)
	s.mu.Unlock()

	log.Printf("[WS] Client connected: %d\n", total)

	defer func() {
		s.removeClient(client)
		_ = conn.Close()
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WS] Disconnect:", err)
			break
		}

		if messageType != websocket.TextMessage {
			continue
		}

		logRealtimeJSON("[WS] Received:", data)
		s.handleMessage(client, data)
	}
}
