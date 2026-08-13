package main

import (
	"encoding/json"
	"log"

	"desktop/backend/protocol"

	"github.com/gorilla/websocket"
)

// ============================================================
// MESSAGE ROUTER
// ============================================================

func (s *RealtimeServer) handleMessage(
	client *WSClient,
	data []byte,
) {
	var header struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &header); err != nil {
		log.Println("[WS] Invalid JSON:", err)
		return
	}

	switch header.Type {
	case "device.register":
		s.handleDeviceRegister(client, data)

	case "client.register":
		s.handleClientRegister(client, data)

	case "mixer.command":
		s.handleMixerCommand(data)

	case "CHANNEL_UPDATE":
		s.handleChannelUpdate(client, data)

	default:
		log.Println("[WS] Unknown:", header.Type)
	}
}

// ============================================================
// DEVICE REGISTER
// ============================================================

func (s *RealtimeServer) handleDeviceRegister(
	client *WSClient,
	data []byte,
) {
	var device protocol.DeviceRegister

	if err := json.Unmarshal(data, &device); err != nil {
		log.Println("[DEVICE] Invalid register:", err)
		return
	}

	client.ID = device.ID
	client.Name = device.Name
	client.Type = "hardware"

	log.Printf(
		"[DEVICE] %s connected\n",
		device.ID,
	)

	reply := map[string]string{
		"type":    "welcome",
		"message": "AMEN backend connected",
	}

	if response, err := json.Marshal(reply); err == nil {
		client.mu.Lock()

		err = client.conn.WriteMessage(
			websocket.TextMessage,
			response,
		)

		client.mu.Unlock()

		if err != nil {
			log.Println(
				"[WS] Welcome write error:",
				err,
			)
		}
	}

	s.BroadcastDeviceStatus(true)
	s.SendState(client)
	s.BroadcastDevices()
}

// ============================================================
// CLIENT REGISTER
// ============================================================

func (s *RealtimeServer) handleClientRegister(
	client *WSClient,
	data []byte,
) {
	var payload struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"clientType"`
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		log.Println("[CLIENT] Invalid register:", err)
		return
	}

	client.ID = payload.ID
	client.Name = payload.Name
	client.Type = payload.Type

	log.Printf(
		"[CLIENT] %s %s\n",
		client.Type,
		client.Name,
	)

	s.SendState(client)
	s.BroadcastDevices()
}

// ============================================================
// MIXER COMMAND
// ============================================================

func (s *RealtimeServer) handleMixerCommand(
	data []byte,
) {
	var cmd protocol.MixerCommand

	if err := json.Unmarshal(data, &cmd); err != nil {
		log.Println("[MIXER] Invalid command:", err)
		return
	}

	log.Printf(
		"[MIXER] CH=%d VALUE=%d\n",
		cmd.Channel,
		cmd.Value,
	)

	s.broadcast(data)
}

// ============================================================
// CHANNEL UPDATE
// ============================================================

func (s *RealtimeServer) handleChannelUpdate(
	client *WSClient,
	data []byte,
) {
	var message protocol.RealtimeMessage

	if err := json.Unmarshal(data, &message); err != nil {
		log.Println(
			"[CHANNEL UPDATE] Invalid JSON:",
			err,
		)
		return
	}

	if message.Channel == nil {
		log.Println(
			"[CHANNEL UPDATE] Missing channel payload",
		)
		return
	}

	update := *message.Channel

	log.Printf(
		"[CHANNEL UPDATE] client=%s channel=%d volume=%v muted=%v\n",
		client.Name,
		update.ID,
		update.Volume,
		update.Muted,
	)

	// Fallback kalau audio handler belum tersedia.
	if s.onChannelUpdate == nil {
		log.Println(
			"[CHANNEL UPDATE] No audio handler - broadcast fallback",
		)

		s.broadcast(data)
		return
	}

	/*
		Audio manager menjadi authoritative state:
		1. apply volume/mute
		2. update state
		3. trigger bridge
		4. broadcast CHANNEL_UPDATE
	*/
	if err := s.onChannelUpdate(update); err != nil {
		log.Println(
			"[CHANNEL UPDATE] Apply failed:",
			err,
		)

		s.SendState(client)
	}
}
