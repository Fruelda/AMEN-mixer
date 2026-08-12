package main

import (
	"encoding/json"
	"log"
	"sort"

	"desktop/backend/models"
	"desktop/backend/protocol"

	"github.com/gorilla/websocket"
)

// ============================================================
// INITIAL CHANNEL STATE
// ============================================================

func (s *RealtimeServer) SetInitialChannels(
	channels []models.Channel,
) {

	s.mu.Lock()

	defer s.mu.Unlock()

	s.channels =
		make(
			map[int]protocol.Channel,
			len(channels),
		)

	for _, channel := range channels {

		s.channels[channel.ID] =
			protocol.Channel{
				ID: channel.ID,

				Name: channel.Name,

				App: channel.App,

				Volume: channel.Volume,

				Muted: channel.Muted,

				Connected: channel.Connected,

				Selected: channel.Selected,
			}
	}

	log.Printf(
		"[STATE] Initialized: %d channels\n",
		len(s.channels),
	)
}

// ============================================================
// CAPTURE CHANNEL UPDATE
// ============================================================
//
// Dipanggil setiap CHANNEL_UPDATE melewati server.
//
// HP -> server
// Wails -> server
// ESP32 -> Wails -> server
//
// Dengan ini state terakhir selalu disimpan.
//
// ============================================================

func (s *RealtimeServer) captureChannelState(
	data []byte,
) {

	var message protocol.RealtimeMessage

	err :=
		json.Unmarshal(
			data,
			&message,
		)

	if err != nil {
		return
	}

	if message.Type !=
		protocol.MessageChannelUpdate {

		return
	}

	if message.Channel == nil {
		return
	}

	update :=
		message.Channel

	s.mu.Lock()

	defer s.mu.Unlock()

	channel, exists :=
		s.channels[update.ID]

	if !exists {

		channel =
			protocol.Channel{
				ID: update.ID,
			}
	}

	if update.Volume != nil {

		channel.Volume =
			*update.Volume
	}

	if update.Muted != nil {

		channel.Muted =
			*update.Muted
	}

	s.channels[update.ID] =
		channel
}

// ============================================================
// GET STATE SNAPSHOT
// ============================================================

func (s *RealtimeServer) getChannelState() []protocol.Channel {

	s.mu.Lock()

	channels :=
		make(
			[]protocol.Channel,
			0,
			len(s.channels),
		)

	for _, channel := range s.channels {

		channels =
			append(
				channels,
				channel,
			)
	}

	s.mu.Unlock()

	// supaya urutan channel konsisten
	// 1, 2, 3, 4, 5...

	sort.Slice(
		channels,
		func(
			i int,
			j int,
		) bool {

			return channels[i].ID <
				channels[j].ID
		},
	)

	return channels
}

// ============================================================
// SEND STATE TO CLIENT
// ============================================================

func (s *RealtimeServer) SendState(
	client *WSClient,
) {

	channels :=
		s.getChannelState()

	message :=
		protocol.RealtimeMessage{
			Type: protocol.MessageState,

			Channels: channels,
		}

	data, err :=
		json.Marshal(
			message,
		)

	if err != nil {

		log.Println(
			"[STATE] Marshal error:",
			err,
		)

		return
	}

	logRealtimeJSON(
		"[STATE] Sending:",
		data,
	)

	client.mu.Lock()

	err =
		client.conn.WriteMessage(
			websocket.TextMessage,
			data,
		)

	client.mu.Unlock()

	if err != nil {

		log.Println(
			"[STATE] Send error:",
			err,
		)

		return
	}

	log.Printf(
		"[STATE] Sent to %s (%d channels)\n",
		client.Name,
		len(channels),
	)
}
