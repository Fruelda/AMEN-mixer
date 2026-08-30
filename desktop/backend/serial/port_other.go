//go:build !windows

package serial

import (
	"fmt"

	bugserial "go.bug.st/serial"
)

type bugSerialPort struct {
	port bugserial.Port
}

func openPlatformPort(portName string) (serialPort, error) {
	mode := &bugserial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   bugserial.NoParity,
		StopBits: bugserial.OneStopBit,

		InitialStatusBits: &bugserial.ModemOutputBits{
			DTR: false,
			RTS: false,
		},
	}

	port, err := bugserial.Open(
		portName,
		mode,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"open serial port %s: %w",
			portName,
			err,
		)
	}

	return &bugSerialPort{
		port: port,
	}, nil
}

func (p *bugSerialPort) Read(buf []byte) (int, error) {
	if p == nil || p.port == nil {
		return 0, fmt.Errorf(
			"serial port is closed",
		)
	}

	return p.port.Read(buf)
}

func (p *bugSerialPort) Write(buf []byte) (int, error) {
	if p == nil || p.port == nil {
		return 0, fmt.Errorf(
			"serial port is closed",
		)
	}

	return p.port.Write(buf)
}

func (p *bugSerialPort) Close() error {
	if p == nil || p.port == nil {
		return nil
	}

	err := p.port.Close()
	p.port = nil

	return err
}
