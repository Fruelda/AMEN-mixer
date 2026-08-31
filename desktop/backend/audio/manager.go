package audio

import (
	"fmt"
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
	manager := &Manager{
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
				Name:      "Music",
				App:       "music",
				Volume:    55,
				Connected: true,
			},
			{
				ID:        4,
				Name:      "Chat",
				App:       "chat",
				Volume:    85,
				Connected: true,
			},
			{
				ID:        5,
				Name:      "Game",
				App:       "game",
				Volume:    60,
				Connected: false,
			},
			{
				ID:        6,
				Name:      "Valeton",
				App:       "valeton",
				Volume:    60,
				Connected: false,
			},
		},
	}

	manager.syncInitialMasterState()

	return manager
}

func (m *Manager) syncInitialMasterState() {
	if m.windows == nil {
		return
	}

	volume, muted, err := m.windows.GetMasterState()
	if err != nil {
		fmt.Printf(
			"[WINDOWS AUDIO] initial master state unavailable: %v\n",
			err,
		)
		return
	}

	for i := range m.channels {
		if m.channels[i].App != "master" {
			continue
		}

		m.channels[i].Volume = volume
		m.channels[i].Muted = muted

		fmt.Printf(
			"[WINDOWS AUDIO] initial master => volume=%d%% muted=%t\n",
			volume,
			muted,
		)

		return
	}
}

// =====================================================
// SET LISTENER
// =====================================================

func (m *Manager) SetListener(
	listener UpdateListener,
) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.listener = listener
}

// =====================================================
// GET CHANNELS
// =====================================================

func (m *Manager) GetChannels() []models.Channel {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(
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

	volume = clampVolume(volume)

	for i := range m.channels {
		channel := &m.channels[i]

		if channel.ID != id {
			continue
		}

		actualVolume := volume

		if m.windows != nil {
			var err error

			actualVolume, err = m.windows.SetVolume(
				channel.App,
				volume,
			)

			if err != nil {
				return err
			}
		}

		channel.Volume = actualVolume

		updated := *channel

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
// SET MUTE
// =====================================================

func (m *Manager) SetMute(
	id int,
	muted bool,
) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.channels {
		channel := &m.channels[i]

		if channel.ID != id {
			continue
		}

		actualMuted := muted

		if m.windows != nil {
			var err error

			actualMuted, err = m.windows.SetMute(
				channel.App,
				muted,
			)

			if err != nil {
				return err
			}
		}

		channel.Muted = actualMuted

		updated := *channel

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
			copy := channel

			return &copy
		}
	}

	return nil
}

// =====================================================
// CLOSE
// =====================================================

func (m *Manager) Close() {
	m.mutex.Lock()

	windows := m.windows
	m.windows = nil

	m.mutex.Unlock()

	if windows != nil {
		windows.Close()
	}
}

func (m *Manager) DebugSessions() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.windows == nil {
		return fmt.Errorf(
			"windows audio not initialized",
		)
	}

	return m.windows.DebugSessions()
}