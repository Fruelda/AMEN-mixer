package serial

import "fmt"

func (m *Manager) Send(data string) error {

	if m.port == nil {
		return fmt.Errorf("serial port not connected")
	}

	_, err := m.port.Write([]byte(data + "\n"))

	return err
}
