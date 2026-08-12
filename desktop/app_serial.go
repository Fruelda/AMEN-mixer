package main

import (
	"fmt"

	"desktop/backend/protocol"
	"desktop/backend/serial"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) startSerial() {
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

			// Wails frontend.
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

			// WebSocket clients.
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
