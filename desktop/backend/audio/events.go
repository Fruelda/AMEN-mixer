package audio

import "desktop/backend/models"

// UpdateListener menerima perubahan audio
// dari Audio Manager.

type UpdateListener interface {
	OnChannelUpdate(
		channel models.Channel,
	)
}
