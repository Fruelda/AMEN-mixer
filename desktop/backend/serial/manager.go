package serial

import (
	"bufio"
	"fmt"

	"desktop/backend/protocol"

	"go.bug.st/serial"
)

type Manager struct {
	port serial.Port

	OnCommand func(*protocol.Command)
}

func New(portName string) (*Manager, error) {

	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	return &Manager{
		port: port,
	}, nil
}

func (m *Manager) Start() {

	fmt.Println("Serial Started...")

	reader := bufio.NewReader(m.port)

	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("READ ERROR:", err)
			continue
		}

		fmt.Printf("RAW: %q\n", line)

		cmd, err := Parse(line)
		if err != nil {
			fmt.Println("PARSE ERROR:", err)
			continue
		}

		fmt.Printf("CMD: %+v\n", cmd)

		if m.OnCommand != nil {
			m.OnCommand(cmd)
		}
	}
}

func (m *Manager) Close() {

	if m.port != nil {
		m.port.Close()
	}
}
