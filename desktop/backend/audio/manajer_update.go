package audio

import "fmt"

// ============================================================
// APPLY CHANNEL UPDATE
// ============================================================
//
// Bisa menerima:
//
// volume saja
// mute saja
// atau keduanya.
//
// nil = field tidak diubah.
//
// ============================================================

func (m *Manager) ApplyChannelUpdate(
	id int,
	volume *int,
	muted *bool,
) error {
	m.mutex.Lock()

	// ========================================================
	// FIND CHANNEL
	// ========================================================

	index :=
		-1

	for i := range m.channels {

		if m.channels[i].ID == id {
			index = i

			break
		}
	}

	if index == -1 {
		m.mutex.Unlock()

		return fmt.Errorf(
			"audio channel %d not found",
			id,
		)
	}

	channel :=
		&m.channels[index]

	changed :=
		false

	// ========================================================
	// VOLUME
	// ========================================================

	if volume != nil {
		newVolume :=
			normalizeVolume(
				*volume,
			)

		if channel.Volume != newVolume {

			if m.windows != nil {
				err :=
					m.windows.SetVolume(
						channel.App,
						newVolume,
					)

				if err != nil {
					m.mutex.Unlock()

					return fmt.Errorf(
						"set volume channel %d: %w",
						id,
						err,
					)
				}
			}

			channel.Volume =
				newVolume

			changed =
				true
		}
	}

	// ========================================================
	// MUTE
	// ========================================================

	if muted != nil {
		newMuted :=
			*muted

		if channel.Muted != newMuted {

			if m.windows != nil {
				err :=
					m.windows.SetMuted(
						channel.App,
						newMuted,
					)

				if err != nil {
					m.mutex.Unlock()

					return fmt.Errorf(
						"set mute channel %d: %w",
						id,
						err,
					)
				}
			}

			channel.Muted =
				newMuted

			changed =
				true
		}
	}

	// Copy final state.
	updated :=
		*channel

	listener :=
		m.listener

	// Lepas mutex sebelum broadcast.
	m.mutex.Unlock()

	// Kalau state sama, tidak usah broadcast.
	if !changed {
		return nil
	}

	// ========================================================
	// AUDIO -> REALTIME
	// ========================================================

	if listener != nil {
		listener.OnChannelUpdate(
			updated,
		)
	}

	return nil
}

// ============================================================
// SET VOLUME
// ============================================================

func (m *Manager) SetVolume(
	id int,
	volume int,
) error {
	return m.ApplyChannelUpdate(
		id,
		&volume,
		nil,
	)
}

// ============================================================
// SET MUTED
// ============================================================

func (m *Manager) SetMuted(
	id int,
	muted bool,
) error {
	return m.ApplyChannelUpdate(
		id,
		nil,
		&muted,
	)
}
