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

	err :=
		json.Unmarshal(
			data,
			&header,
		)

	if err != nil {

		log.Println(
			"[WS] Invalid JSON:",
			err,
		)

		return
	}

	// ========================================================
	// ROUTING
	// ========================================================

	switch header.Type {

	case "device.register":

		s.handleDeviceRegister(
			client,
			data,
		)

	case "client.register":

		s.handleClientRegister(
			client,
			data,
		)

	case "mixer.command":

		s.handleMixerCommand(
			data,
		)

	case "CHANNEL_UPDATE":

		s.handleChannelUpdate(
			client,
			data,
		)

	default:

		log.Println(
			"[WS] Unknown:",
			header.Type,
		)
	}
}

// ============================================================
// DEVICE REGISTER
// ESP32 / HARDWARE
// ============================================================

func (s *RealtimeServer) handleDeviceRegister(
	client *WSClient,
	data []byte,
) {

	var device protocol.DeviceRegister

	err :=
		json.Unmarshal(
			data,
			&device,
		)

	if err != nil {

		log.Println(
			"[DEVICE] Invalid register:",
			err,
		)

		return
	}

	// ========================================================
	// SAVE IDENTITY
	// ========================================================

	client.ID =
		device.ID

	client.Name =
		device.Name

	client.Type =
		"hardware"

	log.Printf(
		"[DEVICE] %s connected\n",
		device.ID,
	)

	// ========================================================
	// WELCOME MESSAGE
	// ========================================================

	reply :=
		map[string]string{
			"type": "welcome",

			"message": "AMEN backend connected",
		}

	response, err :=
		json.Marshal(
			reply,
		)

	if err == nil {

		client.mu.Lock()

		err =
			client.conn.WriteMessage(
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

	// ========================================================
	// DEVICE CONNECTED
	// ========================================================

	s.BroadcastDeviceStatus(
		true,
	)

	// ========================================================
	// SEND CURRENT MIXER STATE
	// ========================================================

	s.SendState(
		client,
	)

	// ========================================================
	// UPDATE DEVICE LIST
	// ========================================================

	s.BroadcastDevices()
}

// ============================================================
// CLIENT REGISTER
// WAILS / IPHONE / ANDROID / BROWSER
// ============================================================

func (s *RealtimeServer) handleClientRegister(
	client *WSClient,
	data []byte,
) {

	var clientData struct {
		ID string `json:"id"`

		Name string `json:"name"`

		Type string `json:"clientType"`
	}

	err :=
		json.Unmarshal(
			data,
			&clientData,
		)

	if err != nil {

		log.Println(
			"[CLIENT] Invalid register:",
			err,
		)

		return
	}

	// ========================================================
	// SAVE IDENTITY
	// ========================================================

	client.ID =
		clientData.ID

	client.Name =
		clientData.Name

	client.Type =
		clientData.Type

	log.Printf(
		"[CLIENT] %s %s\n",
		client.Type,
		client.Name,
	)

	// ========================================================
	// SEND CURRENT STATE
	// ========================================================

	s.SendState(
		client,
	)

	// ========================================================
	// BROADCAST DEVICE LIST
	// ========================================================

	s.BroadcastDevices()
}

// ============================================================
// MIXER COMMAND
// ESP32 COMMAND
// ============================================================

func (s *RealtimeServer) handleMixerCommand(
	data []byte,
) {

	var cmd protocol.MixerCommand

	err :=
		json.Unmarshal(
			data,
			&cmd,
		)

	if err != nil {

		log.Println(
			"[MIXER] Invalid command:",
			err,
		)

		return
	}

	log.Printf(
		"[MIXER] CH=%d VALUE=%d\n",
		cmd.Channel,
		cmd.Value,
	)

	// Broadcast command ke client lain.
	s.broadcast(
		data,
	)
}

// ============================================================
// CHANNEL UPDATE
// HP / BROWSER / DESKTOP
// ============================================================

func (s *RealtimeServer) handleChannelUpdate(
	client *WSClient,
	data []byte,
) {

	var message protocol.RealtimeMessage

	err :=
		json.Unmarshal(
			data,
			&message,
		)

	if err != nil {

		log.Println(
			"[CHANNEL UPDATE] Invalid JSON:",
			err,
		)

		return
	}

	// ========================================================
	// VALIDATE
	// ========================================================

	if message.Channel == nil {

		log.Println(
			"[CHANNEL UPDATE] Missing channel payload",
		)

		return
	}

	update :=
		*message.Channel

	// ========================================================
	// LOG
	// ========================================================

	log.Printf(
		"[CHANNEL UPDATE] client=%s channel=%d volume=%v muted=%v\n",
		client.Name,
		update.ID,
		update.Volume,
		update.Muted,
	)

	// ========================================================
	// FALLBACK
	// ========================================================
	//
	// Kalau audio handler belum dipasang,
	// server masih bisa broadcast seperti sistem lama.
	//
	// ========================================================

	if s.onChannelUpdate == nil {

		log.Println(
			"[CHANNEL UPDATE] No audio handler - broadcast fallback",
		)

		s.broadcast(
			data,
		)

		return
	}

	// ========================================================
	// APPLY TO AUDIO MANAGER
	// ========================================================
	//
	// Jangan broadcast data mentah di sini.
	//
	// audio.Manager akan:
	//
	// 1. apply volume/mute
	// 2. update authoritative state
	// 3. panggil audioBridge
	// 4. BroadcastChannelUpdate()
	//
	// Dengan begitu update tidak broadcast dua kali.
	//
	// ========================================================

	err =
		s.onChannelUpdate(
			update,
		)

	if err != nil {

		log.Println(
			"[CHANNEL UPDATE] Apply failed:",
			err,
		)

		// Kalau gagal, kirim state server yang benar
		// kembali ke client pengirim.
		s.SendState(
			client,
		)

		return
	}
}
