package audio

import "desktop/backend/models"

// ============================================================
// GET ALL CHANNELS
// ============================================================

func (m *Manager) GetChannels() []models.Channel {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	channels := make(
		[]models.Channel,
		len(m.channels),
	)

	copy(
		channels,
		m.channels,
	)

	return channels
}

// ============================================================
// GET ONE CHANNEL
// ============================================================

func (m *Manager) GetChannel(
	id int,
) *models.Channel {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, channel := range m.channels {
		if channel.ID != id {
			continue
		}

		result := channel
		return &result
	}

	return nil
}

// ============================================================
// NORMALIZE VOLUME
// ============================================================

func normalizeVolume(
	volume int,
) int {
	if volume < 0 {
		return 0
	}

	if volume > 100 {
		return 100
	}

	return volume
}
