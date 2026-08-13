package main

import (
	"encoding/json"
	"log"

	"desktop/backend/protocol"

	"github.com/gorilla/websocket"
)

func logRealtimeJSON(label string, data []byte) {
	var payload any

	if err := json.Unmarshal(data, &payload); err != nil {
		log.Println(label, string(data))
		return
	}

	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Println(label, string(data))
		return
	}

	log.Printf("%s\n%s\n", label, string(pretty))
}

func (s *RealtimeServer) broadcast(data []byte) {
	s.captureChannelState(data)
	logRealtimeJSON("[WS] Broadcasting:", data)

	s.mu.Lock()

	clients := make([]*WSClient, 0, len(s.clients))
	for client := range s.clients {
		clients = append(clients, client)
	}

	s.mu.Unlock()

	for _, client := range clients {
		client.mu.Lock()
		err := client.conn.WriteMessage(
			websocket.TextMessage,
			data,
		)
		client.mu.Unlock()

		if err != nil {
			log.Println("[WS] Write error:", err)
			s.removeClient(client)
		}
	}
}

func (s *RealtimeServer) BroadcastJSON(
	message protocol.RealtimeMessage,
) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("[WS] JSON marshal error:", err)
		return
	}

	s.broadcast(data)
}

func (s *RealtimeServer) BroadcastChannelUpdate(
	id int,
	volume int,
	muted bool,
) {
	s.BroadcastJSON(protocol.RealtimeMessage{
		Type: protocol.MessageChannelUpdate,

		Channel: &protocol.ChannelUpdate{
			ID:     id,
			Volume: &volume,
			Muted:  &muted,
		},
	})
}

func (s *RealtimeServer) BroadcastDeviceStatus(
	connected bool,
) {
	s.BroadcastJSON(protocol.RealtimeMessage{
		Type:      protocol.MessageDeviceStatus,
		Connected: &connected,
	})
}
