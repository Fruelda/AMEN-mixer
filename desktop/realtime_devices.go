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

	devices := make(
		[]protocol.DeviceInfo,
		0,
		len(s.clients),
	)

	for client := range s.clients {
		if client.ID == "" {
			continue
		}

		devices = append(
			devices,
			protocol.DeviceInfo{
				ID:         client.ID,
				Name:       client.Name,
				ClientType: client.Type,
				Connected:  true,
			},
		)
	}

	return devices
}

// ============================================================
// BROADCAST DEVICES
// ============================================================

func (s *RealtimeServer) BroadcastDevices() {
	devices := s.getConnectedDevices()

	log.Printf(
		"[DEVICES] Connected: %d\n",
		len(devices),
	)

	s.BroadcastJSON(protocol.RealtimeMessage{
		Type:    protocol.MessageDevices,
		Devices: devices,
	})
}

// ============================================================
// REMOVE CLIENT
// ============================================================

func (s *RealtimeServer) removeClient(
	client *WSClient,
) {
	s.mu.Lock()

	if _, exists := s.clients[client]; !exists {
		s.mu.Unlock()
		return
	}

	delete(s.clients, client)

	total := len(s.clients)

	s.mu.Unlock()

	log.Printf(
		"[WS] Client removed: %d\n",
		total,
	)

	// ========================================================
	// BROADCAST UPDATED DEVICE LIST
	// ========================================================

	s.BroadcastDevices()
}
