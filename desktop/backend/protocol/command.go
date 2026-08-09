package protocol

type Command struct {
	Type    string `json:"type"`
	Channel int    `json:"channel"`
	Value   int    `json:"value"`
}
