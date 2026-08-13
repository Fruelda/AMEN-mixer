package audio

import (
	"sync"

	"desktop/backend/models"
)

// ============================================================
// MANAGER
// ============================================================

type Manager struct {
	mutex sync.RWMutex

	channels []models.Channel
	windows  *WindowsAudio
	listener UpdateListener
}

// ============================================================
// CREATE MANAGER
// ============================================================

func New() *Manager {
	return &Manager{
		windows: NewWindowsAudio(),

		channels: []models.Channel{
			{
				ID:        1,
				Name:      "Master",
				App:       "master",
				Volume:    100,
				Muted:     false,
				Connected: true,
				Selected:  true,
			},
			{
				ID:        2,
				Name:      "Browser",
				App:       "browser",
				Volume:    70,
				Muted:     false,
				Connected: true,
			},
			{
				ID:        3,
				Name:      "Spotify",
				App:       "spotify",
				Volume:    55,
				Muted:     false,
				Connected: true,
			},
			{
				ID:        4,
				Name:      "Discord",
				App:       "discord",
				Volume:    85,
				Muted:     false,
				Connected: true,
			},
			{
				ID:        5,
				Name:      "Valeton",
				App:       "valeton",
				Volume:    60,
				Muted:     false,
				Connected: false,
			},
		},
	}
}

// ============================================================
// SET UPDATE LISTENER
// ============================================================

func (m *Manager) SetListener(
	listener UpdateListener,
) {
	m.mutex.Lock()
	m.listener = listener
	m.mutex.Unlock()
}
