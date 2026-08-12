package main

import (
	"encoding/json"
	"log"

	"desktop/backend/protocol"

	"github.com/gorilla/websocket"
)

// ============================================================
// PRETTY JSON LOG
// ============================================================

func logRealtimeJSON(
	label string,
	data []byte,
) {

	var payload interface{}

	err :=
		json.Unmarshal(
			data,
			&payload,
		)

	if err != nil {

		log.Println(
			label,
			string(data),
		)

		return
	}

	pretty, err :=
		json.MarshalIndent(
			payload,
			"",
			"  ",
		)

	if err != nil {

		log.Println(
			label,
			string(data),
		)

		return
	}

	log.Printf(
		"%s\n%s\n",
		label,
		string(pretty),
	)
}

// ============================================================
// BROADCAST RAW MESSAGE
// ============================================================

func (s *RealtimeServer) broadcast(
	data []byte,
) {

	// ========================================================
	// UPDATE SNAPSHOT
	// ========================================================
	//
	// captureChannelState ada di state_sync.go.
	//
	// Kalau message adalah CHANNEL_UPDATE,
	// state terbaru akan disimpan untuk client
	// yang connect setelahnya.
	//
	// ========================================================

	s.captureChannelState(
		data,
	)

	// ========================================================
	// LOG
	// ========================================================

	logRealtimeJSON(
		"[WS] Broadcasting:",
		data,
	)

	// ========================================================
	// COPY CLIENT LIST
	// ========================================================
	//
	// Jangan menahan server mutex ketika network write.
	//
	// ========================================================

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

	// ========================================================
	// SEND TO ALL CLIENTS
	// ========================================================

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
			"[WS] JSON marshal error:",
			err,
		)

		return
	}

	s.broadcast(
		data,
	)
}

// ============================================================
// BROADCAST CHANNEL UPDATE
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
// BROADCAST DEVICE STATUS
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
