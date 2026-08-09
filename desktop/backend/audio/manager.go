package audio

import (
	"sync"

	"desktop/backend/models"
)

type Manager struct {
	mutex sync.RWMutex

	channels []models.Channel
}

func New() *Manager {

	return &Manager{

		channels: []models.Channel{

			{
				ID:        1,
				Name:      "Master",
				App:       "master",
				Volume:    100,
				Connected: true,
				Selected:  true,
			},

			{
				ID:        2,
				Name:      "Browser",
				App:       "browser",
				Volume:    70,
				Connected: true,
			},

			{
				ID:        3,
				Name:      "Spotify",
				App:       "spotify",
				Volume:    55,
				Connected: true,
			},

			{
				ID:        4,
				Name:      "Discord",
				App:       "discord",
				Volume:    85,
				Connected: true,
			},

			{
				ID:        5,
				Name:      "Valeton",
				App:       "valeton",
				Volume:    60,
				Connected: false,
			},
		},
	}
}

func (m *Manager) GetChannels() []models.Channel {

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.channels
}

func (m *Manager) SetVolume(id int, volume int) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.channels {

		if m.channels[i].ID == id {

			m.channels[i].Volume = volume
			return

		}

	}

}
