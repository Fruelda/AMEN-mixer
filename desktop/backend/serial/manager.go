package serial

import (
	"bufio"
	"fmt"

	"desktop/backend/protocol"

	"go.bug.st/serial"
)

// ============================================================
// MANAGER
// ============================================================

type Manager struct {
	port serial.Port

	OnCommand func(
		*protocol.Command,
	)
}

// ============================================================
// CREATE MANAGER
// ============================================================

func New(
	portName string,
) (*Manager, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(
		portName,
		mode,
	)

	if err != nil {
		return nil, err
	}

	return &Manager{
		port: port,
	}, nil
}

// ============================================================
// START SERIAL
// ============================================================

func (m *Manager) Start() {
	fmt.Println("[SERIAL] Started")

	reader := bufio.NewReader(
		m.port,
	)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(
				"[SERIAL] READ ERROR:",
				err,
			)

			continue
		}

		// ====================================================
		// PARSE COMMAND
		// ====================================================

		cmd, err := Parse(line)

		if err != nil {
			fmt.Println(
				"[SERIAL] PARSE ERROR:",
				err,
			)

			continue
		}

		fmt.Printf(
			"[SERIAL] CMD: %+v\n",
			cmd,
		)

		// ====================================================
		// DISPATCH COMMAND
		// ====================================================

		if m.OnCommand != nil {
			m.OnCommand(cmd)
		}
	}
}

// ============================================================
// CLOSE SERIAL
// ============================================================

func (m *Manager) Close() {
	if m.port == nil {
		return
	}

	_ = m.port.Close()
}
