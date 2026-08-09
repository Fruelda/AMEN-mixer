package audio

type WindowsAudio struct{}

func NewWindowsAudio() *WindowsAudio {
	return &WindowsAudio{}
}

func (w *WindowsAudio) Refresh() error {

	// nanti Core Audio API

	return nil
}

func (w *WindowsAudio) SetVolume(app string, volume int) error {

	// nanti

	return nil
}
