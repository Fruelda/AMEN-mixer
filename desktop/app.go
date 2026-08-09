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

type App struct {
	ctx      context.Context
	audio    *audio.Manager
	realtime *RealtimeServer
}

func NewApp() *App {
	return &App{
		audio: audio.New(),
	}
}

func (a *App) startup(ctx context.Context) {

	a.ctx = ctx

	fmt.Println("===================================")
	fmt.Println("AMEN MIXER STARTING...")
	fmt.Println("===================================")

	// =========================================================
	// REALTIME SERVER
	// =========================================================

	a.realtime = NewRealtimeServer()

	StartRealtimeServer(a.realtime)

	fmt.Println("[APP] REALTIME SERVER STARTED")
	fmt.Println("[APP] WebSocket: ws://0.0.0.0:8081/ws")

	// =========================================================
	// SERIAL / ESP32
	// =========================================================

	manager, err := serial.New(
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

		// PENTING:
		// Jangan return sebelum bagian lain diperlukan.
		// Tapi karena realtime sudah jalan, kita cukup keluar
		// dari fungsi startup.
		return
	}

	fmt.Println("[SERIAL] CONNECTED!")

	// =========================================================
	// ESP32 COMMAND
	// =========================================================

	manager.OnCommand = func(
		cmd *protocol.Command,
	) {

		fmt.Println("===================================")
		fmt.Println("[SERIAL] COMMAND FROM ESP32")
		fmt.Printf("[SERIAL] TYPE    : %s\n", cmd.Type)
		fmt.Printf("[SERIAL] CHANNEL : %d\n", cmd.Channel)
		fmt.Printf("[SERIAL] VALUE   : %d\n", cmd.Value)
		fmt.Println("===================================")

		// =====================================================
		// 1. WAILS DESKTOP
		// =====================================================

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

			fmt.Println(
				"[WAILS] serial-command emitted",
			)
		}

		// =====================================================
		// 2. BROWSER / HP / LAPTOP
		// =====================================================

		if a.realtime != nil {

			a.realtime.BroadcastJSON(
				RealtimeMessage{

					Type: "ESP32_COMMAND",

					Channel: map[string]any{

						"type": cmd.Type,

						"channel": cmd.Channel,

						"value": cmd.Value,
					},
				},
			)

			fmt.Println(
				"[WS] ESP32 command broadcasted",
			)
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
