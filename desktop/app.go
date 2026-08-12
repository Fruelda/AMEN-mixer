package main

import (
	"context"
	"fmt"

	"desktop/backend/audio"
	"desktop/backend/models"
	"desktop/backend/protocol"
	"desktop/backend/serial"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// =============================================================
// APP
// =============================================================

type App struct {
	ctx context.Context

	audio *audio.Manager

	realtime *RealtimeServer
}

// =============================================================
// AUDIO BRIDGE
// =============================================================

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

// =============================================================
// NEW APP
// =============================================================

func NewApp() *App {

	return &App{
		audio: audio.New(),
	}
}

// =============================================================
// STARTUP
// =============================================================

func (a *App) startup(
	ctx context.Context,
) {

	a.ctx =
		ctx

	fmt.Println(
		"===================================",
	)

	fmt.Println(
		"AMEN MIXER STARTING...",
	)

	fmt.Println(
		"===================================",
	)

	// =========================================================
	// REALTIME SERVER
	// =========================================================

	a.realtime =
		NewRealtimeServer()

	// =========================================================
	// INITIAL CHANNEL STATE
	// =========================================================
	//
	// Seed pertama diambil dari Audio Manager.
	//
	// Setelah server berjalan,
	// setiap CHANNEL_UPDATE akan memperbarui
	// realtime state.
	//
	// =========================================================

	a.realtime.SetInitialChannels(
		a.audio.GetChannels(),
	)

	StartRealtimeServer(
		a.realtime,
	)

	fmt.Println(
		"[APP] REALTIME SERVER STARTED",
	)

	fmt.Println(
		"[APP] WebSocket: ws://0.0.0.0:8081/ws",
	)

	// =========================================================
	// AUDIO EVENT BRIDGE
	// =========================================================

	a.audio.SetListener(
		&audioBridge{
			server: a.realtime,
		},
	)

	fmt.Println(
		"[AUDIO] EVENT BRIDGE CONNECTED",
	)

	// =========================================================
	// SERIAL / ESP32
	// =========================================================

	manager, err :=
		serial.New(
			"/dev/cu.usbserial-1420",
		)

	if err != nil {

		fmt.Println(
			"[SERIAL] OPEN ERROR:",
			err,
		)

		fmt.Println(
			"[SERIAL] ESP32 belum terhubung.",
		)

		fmt.Println(
			"[SERIAL] Mixer tetap berjalan tanpa ESP32.",
		)

		return
	}

	fmt.Println(
		"[SERIAL] CONNECTED!",
	)

	// =========================================================
	// ESP32 COMMAND HANDLER
	// =========================================================

	manager.OnCommand =
		func(
			cmd *protocol.Command,
		) {

			fmt.Println(
				"===================================",
			)

			fmt.Println(
				"[SERIAL] COMMAND FROM ESP32",
			)

			fmt.Printf(
				"[SERIAL] TYPE    : %s\n",
				cmd.Type,
			)

			fmt.Printf(
				"[SERIAL] CHANNEL : %d\n",
				cmd.Channel,
			)

			fmt.Printf(
				"[SERIAL] VALUE   : %d\n",
				cmd.Value,
			)

			fmt.Println(
				"===================================",
			)

			// =================================================
			// WAILS EVENT
			// =================================================

			if a.ctx != nil {

				runtime.EventsEmit(
					a.ctx,
					"serial-command",
					map[string]any{
						"type": cmd.Type,

						"channel": cmd.Channel,

						"value": cmd.Value,
					},
				)
			}

			// =================================================
			// WEBSOCKET COMMAND
			// =================================================

			if a.realtime != nil {

				a.realtime.BroadcastJSON(
					protocol.RealtimeMessage{
						Type: protocol.MessageCommand,

						Command: &protocol.MixerCommand{
							Type: cmd.Type,

							Channel: cmd.Channel,

							Value: cmd.Value,
						},
					},
				)
			}
		}

	go manager.Start()

	fmt.Println(
		"[SERIAL] BACKGROUND SERIAL STARTED",
	)
}

// =============================================================
// GREET
// =============================================================

func (a *App) Greet(
	name string,
) string {

	return fmt.Sprintf(
		"Hello %s, It's show time!",
		name,
	)
}

// =============================================================
// GET CHANNELS
// =============================================================

func (a *App) GetChannels() []models.Channel {

	return a.audio.GetChannels()
}

// =============================================================
// SET VOLUME
// =============================================================

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
