package audio

import (
	"sync"

	"desktop/backend/models"
)

type Manager struct {
	mutex sync.RWMutex

	channels []models.Channel

	windows *WindowsAudio
}

func New() *Manager {

	return &Manager{

		windows: NewWindowsAudio(),

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

// =====================================================
// GET CHANNELS
// =====================================================

func (m *Manager) GetChannels() []models.Channel {

	m.mutex.RLock()

	defer m.mutex.RUnlock()

	result :=
		make(
			[]models.Channel,
			len(m.channels),
		)

	copy(
		result,
		m.channels,
	)

	return result

}

// =====================================================
// SET VOLUME
// =====================================================

func (m *Manager) SetVolume(
	id int,
	volume int,
) error {

	m.mutex.Lock()

	defer m.mutex.Unlock()

	for i := range m.channels {

		channel :=
			&m.channels[i]

		if channel.ID == id {

			channel.Volume =
				volume

			// Forward ke Windows Audio

			if m.windows != nil {

				return m.windows.SetVolume(
					channel.App,
					volume,
				)

			}

			return nil

		}

	}

	return nil

}

// =====================================================
// FIND CHANNEL
// =====================================================

func (m *Manager) GetChannel(
	id int,
) *models.Channel {

	m.mutex.RLock()

	defer m.mutex.RUnlock()

	for i := range m.channels {

		if m.channels[i].ID == id {

			channel :=
				m.channels[i]

			return &channel

		}

	}

	return nil

}
