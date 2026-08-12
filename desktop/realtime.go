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

	ID   string
	Name string
	Type string
}

// ============================================================
// SERVER
// ============================================================

type RealtimeServer struct {
	clients map[*WSClient]bool

	channels map[int]protocol.Channel

	mu sync.Mutex
}

// ============================================================
// CONFIG
// ============================================================

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ============================================================
// CREATE SERVER
// ============================================================

func NewRealtimeServer() *RealtimeServer {

	return &RealtimeServer{
		clients: make(
			map[*WSClient]bool,
		),

		channels: make(
			map[int]protocol.Channel,
		),
	}
}

// ============================================================
// PRETTY JSON LOG
// ============================================================
//
// Hanya untuk tampilan terminal.
// Tidak mengubah payload WebSocket.
//
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

	// Jika bukan JSON valid,
	// tampilkan seperti biasa.
	if err != nil {

		log.Println(
			label,
			string(data),
		)

		return
	}

	// Format JSON multiline.
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
// WEBSOCKET HANDLER
// ============================================================

func (s *RealtimeServer) HandleWebSocket(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println(
		"[WS] Incoming connection",
	)

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
			Type: "unknown",
		}

	// ========================================================
	// ADD CLIENT
	// ========================================================

	s.mu.Lock()

	s.clients[client] = true

	total :=
		len(s.clients)

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

		if messageType !=
			websocket.TextMessage {

			continue
		}

		// ====================================================
		// RECEIVED LOG
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

	switch header.Type {

	// ========================================================
	// DEVICE REGISTER
	// ========================================================

	case "device.register":

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

		// ====================================================
		// WELCOME
		// ====================================================

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

		// ====================================================
		// DEVICE STATUS
		// ====================================================

		s.BroadcastDeviceStatus(
			true,
		)

		// ====================================================
		// CONNECTED DEVICES
		// ====================================================
		// Initial state untuk ESP32/device baru
		s.SendState(
			client,
		)
		s.BroadcastDevices()

	// ========================================================
	// CLIENT REGISTER
	// ========================================================

	case "client.register":

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
		// Initial state untuk Wails/HP/iPad
		s.SendState(
			client,
		)
		// ====================================================
		// CONNECTED DEVICES
		// ====================================================

		s.BroadcastDevices()

	// ========================================================
	// ESP32 COMMAND
	// ========================================================

	case "mixer.command":

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

		s.broadcast(
			data,
		)

	// ========================================================
	// CHANNEL UPDATE
	// ========================================================

	case "CHANNEL_UPDATE":

		log.Println(
			"[CHANNEL UPDATE]",
		)

		s.broadcast(
			data,
		)

	// ========================================================
	// UNKNOWN
	// ========================================================

	default:

		log.Println(
			"[WS] Unknown:",
			header.Type,
		)
	}
}

// ============================================================
// CONNECTED DEVICES
// ============================================================

func (s *RealtimeServer) getConnectedDevices() []protocol.DeviceInfo {

	s.mu.Lock()

	defer s.mu.Unlock()

	devices :=
		make(
			[]protocol.DeviceInfo,
			0,
			len(s.clients),
		)

	for client := range s.clients {

		// socket yang belum register
		// tidak ditampilkan

		if client.ID == "" {
			continue
		}

		devices =
			append(
				devices,
				protocol.DeviceInfo{
					ID: client.ID,

					Name: client.Name,

					ClientType: client.Type,

					Connected: true,
				},
			)
	}

	return devices
}

// ============================================================
// BROADCAST DEVICES
// ============================================================

func (s *RealtimeServer) BroadcastDevices() {

	devices :=
		s.getConnectedDevices()

	message :=
		protocol.RealtimeMessage{
			Type: protocol.MessageDevices,

			Devices: devices,
		}

	log.Printf(
		"[DEVICES] Connected: %d\n",
		len(devices),
	)

	s.BroadcastJSON(
		message,
	)
}

// ============================================================
// REMOVE CLIENT
// ============================================================

func (s *RealtimeServer) removeClient(
	client *WSClient,
) {

	s.mu.Lock()

	_, exists :=
		s.clients[client]

	if !exists {

		s.mu.Unlock()

		return
	}

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

	// ========================================================
	// UPDATE CONNECTED DEVICES
	// ========================================================

	s.BroadcastDevices()
}

// ============================================================
// BROADCAST RAW
// ============================================================

func (s *RealtimeServer) broadcast(
	data []byte,
) {

	// ========================================================
	// LOG
	// ========================================================
	// Simpan state terbaru jika message
	// adalah CHANNEL_UPDATE
	s.captureChannelState(
		data,
	)
	logRealtimeJSON(
		"[WS] Broadcasting:",
		data,
	)

	// ========================================================
	// COPY CLIENT LIST
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
	// SEND
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
// JSON BROADCAST
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
// CHANNEL UPDATE
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

	// ========================================================
	// WEBSOCKET
	// ========================================================

	mux.HandleFunc(
		"/ws",
		server.HandleWebSocket,
	)

	// ========================================================
	// HEALTH CHECK
	// ========================================================

	mux.HandleFunc(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			_, _ =
				w.Write(
					[]byte(
						"AMEN MIXER REALTIME SERVER OK",
					),
				)
		},
	)

	// ========================================================
	// START SERVER
	// ========================================================

	go func() {

		log.Println(
			"[WS] =================================",
		)

		log.Println(
			"[WS] AMEN REALTIME SERVER",
		)

		log.Println(
			"[WS] Listening :8081",
		)

		httpServer :=
			&http.Server{
				Addr: "0.0.0.0:8081",

				Handler: mux,

				ReadHeaderTimeout: 5 * time.Second,
			}

		err :=
			httpServer.ListenAndServe()

		if err != nil {

			log.Println(
				"[WS]",
				err,
			)
		}
	}()
}
