package models

type Channel struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	App       string `json:"app"`
	Volume    int    `json:"volume"`
	Muted     bool   `json:"muted"`
	Connected bool   `json:"connected"`
	Selected  bool   `json:"selected"`
}
