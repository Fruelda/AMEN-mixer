package protocol

/*
|--------------------------------------------------------------------------
| REALTIME MESSAGE
|--------------------------------------------------------------------------
|
| Format utama WebSocket
|
|--------------------------------------------------------------------------
*/

type MessageType string

const (
	MessageState MessageType = "STATE"

	MessageChannelUpdate MessageType = "CHANNEL_UPDATE"

	MessageCommand MessageType = "COMMAND"

	MessageDeviceStatus MessageType = "DEVICE_STATUS"
)

type RealtimeMessage struct {
	Type MessageType `json:"type"`

	Channels []Channel `json:"channels,omitempty"`

	Channel *ChannelUpdate `json:"channel,omitempty"`

	Command *MixerCommand `json:"command,omitempty"`

	Connected *bool `json:"connected,omitempty"`
}

/*
|--------------------------------------------------------------------------
| CHANNEL
|--------------------------------------------------------------------------
*/

type Channel struct {
	ID int `json:"id"`

	Name string `json:"name"`

	App string `json:"app"`

	Volume int `json:"volume"`

	Muted bool `json:"muted"`

	Connected bool `json:"connected"`

	Selected bool `json:"selected"`
}

/*
|--------------------------------------------------------------------------
| CHANNEL UPDATE
|--------------------------------------------------------------------------
*/

type ChannelUpdate struct {
	ID int `json:"id"`

	Volume *int `json:"volume,omitempty"`

	Muted *bool `json:"muted,omitempty"`
}
