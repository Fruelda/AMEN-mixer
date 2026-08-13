package protocol

// ============================================================
// COMMAND PROTOCOL
// ============================================================
//
// Alur:
//
// ESP32
//   ↓ Serial
// Go Backend
//   ↓ WebSocket
// Vue Frontend
//
// ============================================================

type CommandType string

const (
	CommandEncoder CommandType = "ENC"
	CommandButton  CommandType = "BTN"
)

// ============================================================
// MIXER COMMAND
// ============================================================

type MixerCommand struct {
	Type CommandType `json:"type"`

	// Channel yang dikontrol:
	// 1 = Master
	// 2 = Browser
	// 3 = Spotify
	// dst.
	Channel int `json:"channel"`

	// ENC:
	// +/- perubahan volume.
	//
	// BTN:
	// 1 = pressed.
	Value int `json:"value"`
}

// ============================================================
// BACKWARD COMPATIBILITY
// ============================================================
//
// Serial layer masih menggunakan protocol.Command.
// Alias ini membuat migrasi tidak perlu dilakukan sekaligus.
//
// ============================================================

type Command = MixerCommand
