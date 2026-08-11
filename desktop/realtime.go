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

	conn, err := upgrader.Upgrade(
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

	client := &WSClient{

		conn: conn,

		Type: "unknown",
	}

	s.mu.Lock()

	s.clients[client] = true

	total := len(s.clients)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client connected: %d\n",
		total,
	)

	defer func() {

		s.removeClient(client)

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

	err := json.Unmarshal(
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

	case "device.register":

		var device protocol.DeviceRegister

		err :=
			json.Unmarshal(
				data,
				&device,
			)

		if err != nil {
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

			"type": "welcome",

			"message": "AMEN backend connected",
		}

		response, _ :=
			json.Marshal(reply)

		client.mu.Lock()

		client.conn.WriteMessage(
			websocket.TextMessage,
			response,
		)

		client.mu.Unlock()

		s.BroadcastDeviceStatus(true)

	case "client.register":

		var clientData struct {
			ID string `json:"id"`

			Name string `json:"name"`

			Type string `json:"clientType"`
		}

		json.Unmarshal(
			data,
			&clientData,
		)

		client.ID = clientData.ID
		client.Name = clientData.Name
		client.Type = clientData.Type

		log.Printf(
			"[CLIENT] %s %s\n",
			client.Type,
			client.Name,
		)

	case "mixer.command":

		var cmd protocol.MixerCommand

		err :=
			json.Unmarshal(
				data,
				&cmd,
			)

		if err != nil {
			return
		}

		log.Printf(
			"[MIXER] CH=%d VALUE=%d\n",
			cmd.Channel,
			cmd.Value,
		)

		s.broadcast(data)

	case "CHANNEL_UPDATE":

		log.Println(
			"[CHANNEL UPDATE]",
		)

		s.broadcast(data)

	default:

		log.Println(
			"[WS] Unknown:",
			header.Type,
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
	log.Println(
		"[WS] Broadcasting:",
		string(data),
	)
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

			s.removeClient(client)

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
		json.Marshal(message)

	if err != nil {

		return
	}

	s.broadcast(data)

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

	s.BroadcastJSON(message)

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

	s.BroadcastJSON(message)

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

			_, _ =
				w.Write(
					[]byte(
						"AMEN MIXER REALTIME SERVER OK",
					),
				)

		},
	)

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
