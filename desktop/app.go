package main

import (
	"context"
	"fmt"

	"desktop/backend/audio"
	"desktop/backend/models"
	"desktop/backend/protocol"
	"desktop/backend/serial"
)

// =============================================================
// APP
// =============================================================

type App struct {
	ctx context.Context

	audio *audio.Manager

	realtime *RealtimeServer

	serial *serial.Manager
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

	volume := channel.Volume
	muted := channel.Muted

	b.server.BroadcastJSON(
		protocol.RealtimeMessage{
			Type: protocol.MessageChannelUpdate,

			Channel: &protocol.ChannelUpdate{
				ID:     channel.ID,
				Volume: &volume,
				Muted:  &muted,
			},
		},
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
	go func() {
	err := a.audio.DebugSessions()

	if err != nil {
		fmt.Println(
			"[AUDIO DEBUG ERROR]",
			err,
		)
	}
}()

	// =========================================================
	// SERIAL / ESP32
	// =========================================================

	manager, err :=
		serial.NewAuto()

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

	fmt.Printf(
		"[SERIAL] CONNECTED: %s\n",
		manager.PortName(),
	)

	a.serial = manager

	// =========================================================
	// REGISTER HARDWARE DEVICE
	// =========================================================

	if a.realtime != nil {
		a.realtime.RegisterLocalDevice(
			"amen-mixer-01",
			"AMEN Hardware Mixer",
			"hardware",
		)
	}

	// =========================================================
	// DISCONNECT HANDLER
	// =========================================================

	manager.OnDisconnect =
		func(err error) {
			fmt.Println(
				"[SERIAL] ESP32 DISCONNECTED:",
				err,
			)

			if a.realtime != nil {
				a.realtime.UnregisterLocalDevice(
					"amen-mixer-01",
				)
			}
		}

	// =========================================================
	// ESP32 COMMAND HANDLER
	// =========================================================
	//
	// Hardware command diproses langsung di Go.
	//
	// ESP32
	//   ↓
	// serial.Manager
	//   ↓
	// audio.Manager
	//   ↓
	// Windows Core Audio
	//   ↓
	// audioBridge
	//   ↓
	// CHANNEL_UPDATE
	//   ↓
	// Frontend
	//
	// Tidak menggunakan serial-command Wails event untuk
	// mengeksekusi audio.
	//
	// Tidak mengirim COMMAND WebSocket untuk dieksekusi ulang.
	//
	// Dengan demikian 1 physical encoder detent = 1 audio update.
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

			if a.audio == nil {
				fmt.Println(
					"[SERIAL] AUDIO MANAGER NOT AVAILABLE",
				)

				return
			}

			// =================================================
			// ENCODER
			// =================================================

			if cmd.Type == "ENC" {
				channel :=
					a.audio.GetChannel(
						cmd.Channel,
					)

				if channel == nil {
					fmt.Printf(
						"[SERIAL] Unknown channel: %d\n",
						cmd.Channel,
					)

					return
				}

				newVolume :=
					channel.Volume +
						cmd.Value

				err :=
					a.audio.SetVolume(
						cmd.Channel,
						newVolume,
					)

				if err != nil {
					fmt.Printf(
						"[SERIAL] SetVolume failed: %v\n",
						err,
					)
				}

				return
			}

			// =================================================
			// BUTTON
			// =================================================

			if cmd.Type == "BTN" &&
				cmd.Value == 1 {

				channel :=
					a.audio.GetChannel(
						cmd.Channel,
					)

				if channel == nil {
					fmt.Printf(
						"[SERIAL] Unknown channel: %d\n",
						cmd.Channel,
					)

					return
				}

				err :=
					a.audio.SetMute(
						cmd.Channel,
						!channel.Muted,
					)

				if err != nil {
					fmt.Printf(
						"[SERIAL] SetMute failed: %v\n",
						err,
					)
				}

				return
			}
		}

	// =========================================================
	// START SERIAL
	// =========================================================

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
//
// Digunakan oleh UI.
//
// ESP32 tidak perlu melewati frontend.
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

// =============================================================
// SET MUTE
// =============================================================
//
// Digunakan oleh UI.
//
// ESP32 button diproses langsung pada OnCommand.
// =============================================================

func (a *App) SetMute(
	id int,
	muted bool,
) error {
	if a.audio == nil {
		return fmt.Errorf(
			"audio manager not initialized",
		)
	}

	return a.audio.SetMute(
		id,
		muted,
	)
}

// =============================================================
// SHUTDOWN
// =============================================================

func (a *App) shutdown(
	ctx context.Context,
) {
	_ = ctx

	fmt.Println(
		"[APP] SHUTDOWN",
	)

	if a.serial != nil {
		fmt.Println(
			"[SERIAL] CLOSING",
		)

		a.serial.Close()
		a.serial = nil

		fmt.Println(
			"[SERIAL] CLOSED",
		)
	}

	if a.audio != nil {
		fmt.Println(
			"[AUDIO] CLOSING",
		)

		a.audio.Close()
		a.audio = nil

		fmt.Println(
			"[AUDIO] CLOSED",
		)
	}
}
