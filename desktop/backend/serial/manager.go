package serial

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"desktop/backend/protocol"

	bugserial "go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const preferredUSBSerial = "A5069RR4"

const (
	serialChunkSize = 512
	maxPendingBytes = 8192
)

type Manager struct {
	port serialPort

	portName string

	OnCommand func(*protocol.Command)

	OnDisconnect func(error)
}

// New opens an explicitly selected serial port.
// Keep this function so the existing architecture can still use a fixed port
// when needed.
func New(portName string) (*Manager, error) {
	port, err := openPlatformPort(portName)
	if err != nil {
		return nil, err
	}

	return &Manager{
		port:     port,
		portName: portName,
	}, nil
}

// NewAuto detects the ESP32 serial port first, then opens it.
//
// Selection order:
//  1. AMEN_SERIAL_PORT environment variable, if explicitly set.
//  2. USB serial number already present in this project's PlatformIO config.
//  3. The only USB serial device, when exactly one exists.
//  4. The only serial port, when exactly one exists.
//
// If Windows exposes several ambiguous ports, we deliberately return an error
// instead of silently opening a random COM port.
func NewAuto() (*Manager, error) {
	portName, err := DetectPort()
	if err != nil {
		return nil, err
	}

	return New(portName)
}

func DetectPort() (string, error) {
	if configured := strings.TrimSpace(os.Getenv("AMEN_SERIAL_PORT")); configured != "" {
		fmt.Printf("[SERIAL] Using AMEN_SERIAL_PORT=%s\n", configured)
		return configured, nil
	}

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return detectPortFallback(err)
	}

	if len(ports) == 0 {
		return "", fmt.Errorf("no serial ports found")
	}

	usbPorts := make([]string, 0, len(ports))
	allPorts := make([]string, 0, len(ports))

	for _, port := range ports {
		if port == nil || strings.TrimSpace(port.Name) == "" {
			continue
		}

		allPorts = append(allPorts, port.Name)

		if port.IsUSB {
			usbPorts = append(usbPorts, port.Name)
		}

		fmt.Printf(
			"[SERIAL] Found port=%s usb=%t vid=%s pid=%s serial=%s product=%s manufacturer=%s\n",
			port.Name,
			port.IsUSB,
			port.VID,
			port.PID,
			port.SerialNumber,
			port.Product,
			port.Manufacturer,
		)

		if port.IsUSB && strings.EqualFold(strings.TrimSpace(port.SerialNumber), preferredUSBSerial) {
			fmt.Printf(
				"[SERIAL] Matched project USB serial %s on %s\n",
				preferredUSBSerial,
				port.Name,
			)

			return port.Name, nil
		}
	}

	if len(usbPorts) == 1 {
		fmt.Printf("[SERIAL] Selected only USB serial port: %s\n", usbPorts[0])
		return usbPorts[0], nil
	}

	if len(allPorts) == 1 {
		fmt.Printf("[SERIAL] Selected only serial port: %s\n", allPorts[0])
		return allPorts[0], nil
	}

	if len(usbPorts) > 1 {
		return "", fmt.Errorf(
			"multiple USB serial ports found (%s); set AMEN_SERIAL_PORT to the ESP32 COM port, for example COM5",
			strings.Join(usbPorts, ", "),
		)
	}

	return "", fmt.Errorf(
		"multiple serial ports found (%s); set AMEN_SERIAL_PORT to the ESP32 COM port, for example COM5",
		strings.Join(allPorts, ", "),
	)
}

func detectPortFallback(enumerationErr error) (string, error) {
	ports, err := bugserial.GetPortsList()
	if err != nil {
		return "", fmt.Errorf(
			"serial enumeration failed: %v; fallback failed: %w",
			enumerationErr,
			err,
		)
	}

	if len(ports) == 0 {
		return "", fmt.Errorf(
			"no serial ports found: %v",
			enumerationErr,
		)
	}

	for _, port := range ports {
		fmt.Printf(
			"[SERIAL] Found port=%s (basic enumeration)\n",
			port,
		)
	}

	if len(ports) == 1 {
		fmt.Printf(
			"[SERIAL] Selected only serial port: %s\n",
			ports[0],
		)

		return ports[0], nil
	}

	return "", fmt.Errorf(
		"USB metadata enumeration failed (%v) and multiple serial ports exist (%s); set AMEN_SERIAL_PORT to the ESP32 COM port",
		enumerationErr,
		strings.Join(ports, ", "),
	)
}

func (m *Manager) PortName() string {
	if m == nil {
		return ""
	}

	return m.portName
}

// Start reads raw serial chunks and performs line framing locally.
//
// The platform transport may validly return n=0 on an idle timeout,
// so the manager simply waits for the next read without treating
// that as EOF.
func (m *Manager) Start() {
	fmt.Printf(
		"[SERIAL] Reader started on %s\n",
		m.portName,
	)

	chunk := make([]byte, serialChunkSize)
	pending := make(
		[]byte,
		0,
		serialChunkSize*2,
	)

	for {
		n, err := m.port.Read(chunk)
		if err != nil {
			fmt.Println(
				"[SERIAL] READ ERROR:",
				err,
			)

			if m.OnDisconnect != nil {
				m.OnDisconnect(err)
			}

			return
		}

		if n == 0 {
			continue
		}

		pending = append(
			pending,
			chunk[:n]...,
		)

		for {
			newline := bytes.IndexByte(
				pending,
				'\n',
			)

			if newline < 0 {
				break
			}

			lineBytes := pending[:newline+1]
			line := string(lineBytes)

			remainder := len(pending) - (newline + 1)

			copy(
				pending,
				pending[newline+1:],
			)

			pending = pending[:remainder]

			m.handleLine(line)
		}

		if len(pending) > maxPendingBytes {
			fmt.Printf(
				"[SERIAL] Dropping oversized unterminated frame (%d bytes)\n",
				len(pending),
			)

			pending = pending[:0]
		}
	}
}

func (m *Manager) handleLine(line string) {
	fmt.Printf(
		"[SERIAL] RAW: %q\n",
		line,
	)

	cmd, err := Parse(line)
	if err != nil {
		// Firmware boot/status/debug lines share
		// the serial port with commands.
		return
	}

	fmt.Printf(
		"[SERIAL] CMD: %+v\n",
		cmd,
	)

	if m.OnCommand != nil {
		m.OnCommand(cmd)
	}
}

func (m *Manager) Close() {
	if m != nil && m.port != nil {
		_ = m.port.Close()
	}
}
