package audio

import "desktop/backend/models"

// ============================================================
// UPDATE LISTENER
// ============================================================
//
// Menerima perubahan channel dari Audio Manager.
//
// Implementasinya digunakan oleh realtime bridge
// untuk mengirim update audio ke WebSocket.
//
// ============================================================

type UpdateListener interface {
	OnChannelUpdate(
		channel models.Channel,
	)
}
