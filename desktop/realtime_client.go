package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsPongWait      = 12 * time.Second
	wsPingInterval  = 5 * time.Second
	wsPingWriteWait = 3 * time.Second
)

// ============================================================
// HANDLE WEBSOCKET
// ============================================================

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

	_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(
			time.Now().Add(wsPongWait),
		)
	})

	s.mu.Lock()
	s.clients[client] = true
	total := len(s.clients)
	s.mu.Unlock()

	log.Printf(
		"[WS] Client connected: %d\n",
		total,
	)

	done := make(chan struct{})

	go s.runClientHeartbeat(
		client,
		done,
	)

	defer func() {
		close(done)

		s.removeClient(client)

		_ = conn.Close()
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

		logRealtimeJSON(
			"[WS] Received:",
			data,
		)

		s.handleMessage(
			client,
			data,
		)
	}
}

// ============================================================
// HEARTBEAT
// ============================================================

func (s *RealtimeServer) runClientHeartbeat(
	client *WSClient,
	done <-chan struct{},
) {
	ticker := time.NewTicker(
		wsPingInterval,
	)

	defer ticker.Stop()

	for {
		select {
		case <-done:
			return

		case <-ticker.C:
			client.mu.Lock()

			err := client.conn.WriteControl(
				websocket.PingMessage,
				nil,
				time.Now().Add(
					wsPingWriteWait,
				),
			)

			client.mu.Unlock()

			if err != nil {
				log.Println(
					"[WS] Heartbeat failed:",
					err,
				)

				_ = client.conn.Close()

				return
			}
		}
	}
}
