package audio

import (
	"fmt"
	"strings"
)

type WindowsAudio struct{}

func NewWindowsAudio() *WindowsAudio {
	return &WindowsAudio{}
}

// ============================================================
// REFRESH
// ============================================================

func (w *WindowsAudio) Refresh() error {
	fmt.Println(
		"[WINDOWS AUDIO] Refresh sessions",
	)

	// TODO:
	// Windows Core Audio API:
	//
	// - enumerate sessions
	// - detect apps
	// - map ke mixer channels

	return nil
}

// ============================================================
// SET VOLUME
// ============================================================

func (w *WindowsAudio) SetVolume(
	app string,
	volume int,
) error {
	app =
		strings.ToLower(
			strings.TrimSpace(
				app,
			),
		)

	volume =
		normalizeVolume(
			volume,
		)

	fmt.Printf(
		"[WINDOWS AUDIO] %s => %d%%\n",
		app,
		volume,
	)

	// TODO Windows:
	//
	// Cari audio session berdasarkan process.
	//
	// ISimpleAudioVolume
	// SetMasterVolume(
	//     float32(volume) / 100,
	// )

	return nil
}

// ============================================================
// SET MUTE
// ============================================================

func (w *WindowsAudio) SetMuted(
	app string,
	muted bool,
) error {
	app =
		strings.ToLower(
			strings.TrimSpace(
				app,
			),
		)

	fmt.Printf(
		"[WINDOWS AUDIO] %s muted => %t\n",
		app,
		muted,
	)

	// TODO Windows:
	//
	// ISimpleAudioVolume.SetMute(...)

	return nil
}
