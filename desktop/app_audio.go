package main

import (
	"fmt"

	"desktop/backend/models"
	"desktop/backend/protocol"
)

// ============================================================
// AUDIO -> WEBSOCKET BRIDGE
// ============================================================

type audioBridge struct {
	server *RealtimeServer
}

func (b *audioBridge) OnChannelUpdate(
	channel models.Channel,
) {
	if b.server == nil {
		return
	}

	b.server.BroadcastChannelUpdate(
		channel.ID,
		channel.Volume,
		channel.Muted,
	)
}

// ============================================================
// SETUP REALTIME
// ============================================================

func (a *App) setupRealtime() {
	if a.audio == nil {
		return
	}

	// Create realtime server.
	a.realtime =
		NewRealtimeServer()

	// Initial mixer state.
	a.realtime.SetInitialChannels(
		a.audio.GetChannels(),
	)

	// Audio -> WebSocket.
	a.audio.SetListener(
		&audioBridge{
			server: a.realtime,
		},
	)

	// WebSocket -> Audio.
	a.realtime.SetChannelUpdateHandler(
		func(
			update protocol.ChannelUpdate,
		) error {
			return a.audio.ApplyChannelUpdate(
				update.ID,
				update.Volume,
				update.Muted,
			)
		},
	)

	// Start server :8081.
	StartRealtimeServer(
		a.realtime,
	)

	fmt.Println(
		"[APP] REALTIME SERVER STARTED",
	)

	fmt.Println(
		"[APP] WebSocket: ws://0.0.0.0:8081/ws",
	)

	fmt.Println(
		"[AUDIO] EVENT BRIDGE CONNECTED",
	)
}

// ============================================================
// GET CHANNELS
// ============================================================

func (a *App) GetChannels() []models.Channel {
	if a.audio == nil {
		return []models.Channel{}
	}

	return a.audio.GetChannels()
}

// ============================================================
// SET VOLUME
// ============================================================

func (a *App) SetVolume(
	id int,
	volume int,
) error {
	if a.audio == nil {
		return fmt.Errorf(
			"audio manager not initialized",
		)
	}

	return a.audio.SetVolume(
		id,
		volume,
	)
}

// ============================================================
// SET MUTED
// ============================================================

func (a *App) SetMuted(
	id int,
	muted bool,
) error {
	if a.audio == nil {
		return fmt.Errorf(
			"audio manager not initialized",
		)
	}

	return a.audio.SetMuted(
		id,
		muted,
	)
}
