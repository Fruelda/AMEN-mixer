package protocol

/*
|--------------------------------------------------------------------------
| AMEN Mixer Command Protocol
|--------------------------------------------------------------------------
|
| Komunikasi:
|
| ESP32
|  |
|  | Serial
|  |
| Go Backend
|  |
|  | WebSocket
|  |
| Vue Frontend
|
|--------------------------------------------------------------------------
*/

type CommandType string

const (
	CommandEncoder CommandType = "ENC"

	CommandButton CommandType = "BTN"
)

type MixerCommand struct {
	Type CommandType `json:"type"`

	/*
		channel yang dikontrol

		1 = Master
		2 = Browser
		3 = Spotify

	*/

	Channel int `json:"channel"`

	/*
		ENC:
		+/- volume

		BTN:
		1 pressed

	*/

	Value int `json:"value"`
}

// Backward compatibility
type Command = MixerCommand
