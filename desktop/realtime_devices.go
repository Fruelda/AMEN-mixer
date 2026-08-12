package main

import (
	"log"

	"desktop/backend/protocol"
)

// ============================================================
// GET CONNECTED DEVICES
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

		// ====================================================
		// IGNORE UNREGISTERED SOCKET
		// ====================================================

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
// BROADCAST CONNECTED DEVICES
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
		len(
			s.clients,
		)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client removed: %d\n",
		total,
	)

	// ========================================================
	// SEND UPDATED DEVICE LIST
	// ========================================================

	s.BroadcastDevices()
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
