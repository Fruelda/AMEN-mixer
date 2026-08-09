package audio

import (
	"fmt"
	"strings"
)

type WindowsAudio struct {
}

func NewWindowsAudio() *WindowsAudio {

	return &WindowsAudio{}

}

// =====================================================
// REFRESH AUDIO SESSION
// =====================================================

func (w *WindowsAudio) Refresh() error {

	fmt.Println(
		"[WINDOWS AUDIO] Refresh sessions",
	)

	/*
		Nanti:

		Core Audio API:
		- enumerate sessions
		- detect apps
		- update channel list

	*/

	return nil

}

// =====================================================
// SET APPLICATION VOLUME
// =====================================================

func (w *WindowsAudio) SetVolume(
	app string,
	volume int,
) error {

	app =
		strings.ToLower(
			app,
		)

	fmt.Printf(
		"[WINDOWS AUDIO] %s => %d%%\n",
		app,
		volume,
	)

	/*
		Nanti:

		1. Cari audio session:
		   app.exe


		2. Ambil:
		   ISimpleAudioVolume


		3. Set:

		   volume / 100.0


	*/

	return nil

}
