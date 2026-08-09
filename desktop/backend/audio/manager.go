package audio

import (
	"sync"

	"desktop/backend/models"
)

type Manager struct {
	mutex sync.RWMutex

	channels []models.Channel

	windows *WindowsAudio

	listener UpdateListener
}

// =====================================================
// NEW
// =====================================================

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
// SET LISTENER
// =====================================================

func (m *Manager) SetListener(
	listener UpdateListener,
) {

	m.listener =
		listener

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

		if channel.ID != id {

			continue

		}

		channel.Volume =
			volume

		updated :=
			*channel

		// Windows Audio

		if m.windows != nil {

			err :=
				m.windows.SetVolume(
					channel.App,
					volume,
				)

			if err != nil {

				return err

			}

		}

		// Notify realtime bridge

		if m.listener != nil {

			m.listener.OnChannelUpdate(
				updated,
			)

		}

		return nil

	}

	return nil

}

// =====================================================
// GET SINGLE CHANNEL
// =====================================================

func (m *Manager) GetChannel(
	id int,
) *models.Channel {

	m.mutex.RLock()

	defer m.mutex.RUnlock()

	for _, channel := range m.channels {

		if channel.ID == id {

			copy :=
				channel

			return &copy

		}

	}

	return nil

}
