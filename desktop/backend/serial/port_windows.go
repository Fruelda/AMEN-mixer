//go:build windows

package serial

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

// Bit mask DCB Win32.
//
// Port sengaja dikonfigurasi 115200 8N1 tanpa hardware/software
// flow control. DTR dan RTS juga dimatikan supaya membuka COM port
// tidak ikut memainkan control line ESP32-S3.
const (
	dcbBinary uint32 = 0x00000001

	dcbOutXCTSFlow uint32 = 0x00000004
	dcbOutXDSRFlow uint32 = 0x00000008

	dcbDTRControlMask uint32 = 0x00000030

	dcbDSRSensitivity uint32 = 0x00000040

	dcbTXContinueOnXOFF uint32 = 0x00000080

	dcbOutX uint32 = 0x00000100
	dcbInX  uint32 = 0x00000200

	dcbErrorChar uint32 = 0x00000400
	dcbNull      uint32 = 0x00000800

	dcbRTSControlMask uint32 = 0x00003000

	dcbAbortOnError uint32 = 0x00004000
)

type windowsSerialPort struct {
	handle windows.Handle
}

func openPlatformPort(portName string) (serialPort, error) {
	name := strings.TrimSpace(portName)

	if name == "" {
		return nil, fmt.Errorf("empty serial port name")
	}

	// Win32 device namespace.
	//
	// Bentuk ini aman juga untuk COM10, COM11, dst.
	if !strings.HasPrefix(name, `\\.\`) {
		name = `\\.\` + name
	}

	path, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, fmt.Errorf(
			"serial path %q: %w",
			portName,
			err,
		)
	}

	handle, err := windows.CreateFile(
		path,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"open %s: %w",
			portName,
			err,
		)
	}

	port := &windowsSerialPort{
		handle: handle,
	}

	if err := port.configure115200N81(); err != nil {
		_ = port.Close()
		return nil, err
	}

	// Buang data lama yang sudah berada di RX queue sebelum Wails
	// mengambil alih COM port.
	if err := windows.PurgeComm(
		handle,
		windows.PURGE_RXCLEAR|windows.PURGE_RXABORT,
	); err != nil {
		_ = port.Close()

		return nil, fmt.Errorf(
			"purge %s input buffer: %w",
			portName,
			err,
		)
	}

	fmt.Printf(
		"[SERIAL] Windows native COM transport active on %s\n",
		portName,
	)

	return port, nil
}

func (p *windowsSerialPort) configure115200N81() error {
	if p == nil || p.handle == 0 {
		return fmt.Errorf("serial port is closed")
	}

	var dcb windows.DCB

	if err := windows.GetCommState(
		p.handle,
		&dcb,
	); err != nil {
		return fmt.Errorf(
			"GetCommState: %w",
			err,
		)
	}

	dcb.BaudRate = 115200
	dcb.ByteSize = 8
	dcb.Parity = windows.NOPARITY
	dcb.StopBits = windows.ONESTOPBIT

	// Binary serial mode.
	dcb.Flags |= dcbBinary

	// DTR_CONTROL_DISABLE = 0.
	dcb.Flags &^= dcbDTRControlMask

	// RTS_CONTROL_DISABLE = 0.
	dcb.Flags &^= dcbRTSControlMask

	// Hardware flow control OFF.
	dcb.Flags &^= dcbOutXCTSFlow
	dcb.Flags &^= dcbOutXDSRFlow
	dcb.Flags &^= dcbDSRSensitivity

	// Software XON/XOFF OFF.
	dcb.Flags &^= dcbOutX
	dcb.Flags &^= dcbInX

	// Jangan transform karakter atau abort otomatis.
	dcb.Flags &^= dcbErrorChar
	dcb.Flags &^= dcbNull
	dcb.Flags &^= dcbAbortOnError

	// Tetap izinkan transmisi.
	dcb.Flags |= dcbTXContinueOnXOFF

	dcb.XonChar = 17
	dcb.XoffChar = 19

	if err := windows.SetCommState(
		p.handle,
		&dcb,
	); err != nil {
		return fmt.Errorf(
			"SetCommState: %w",
			err,
		)
	}

	// Read akan bangun periodik ketika idle.
	// Kalau ada data, ReadFile tetap dapat mengembalikannya langsung.
	timeouts := windows.CommTimeouts{
		ReadIntervalTimeout:         0xFFFFFFFF,
		ReadTotalTimeoutMultiplier:  0,
		ReadTotalTimeoutConstant:    250,
		WriteTotalTimeoutMultiplier: 0,
		WriteTotalTimeoutConstant:   1000,
	}

	if err := windows.SetCommTimeouts(
		p.handle,
		&timeouts,
	); err != nil {
		return fmt.Errorf(
			"SetCommTimeouts: %w",
			err,
		)
	}

	return nil
}

func (p *windowsSerialPort) Read(buf []byte) (int, error) {
	if p == nil || p.handle == 0 {
		return 0, fmt.Errorf(
			"serial port is closed",
		)
	}

	if len(buf) == 0 {
		return 0, nil
	}

	var bytesRead uint32

	err := windows.ReadFile(
		p.handle,
		buf,
		&bytesRead,
		nil,
	)
	if err != nil {
		return int(bytesRead), fmt.Errorf(
			"ReadFile: %w",
			err,
		)
	}

	return int(bytesRead), nil
}

func (p *windowsSerialPort) Write(buf []byte) (int, error) {
	if p == nil || p.handle == 0 {
		return 0, fmt.Errorf(
			"serial port is closed",
		)
	}

	if len(buf) == 0 {
		return 0, nil
	}

	var bytesWritten uint32

	err := windows.WriteFile(
		p.handle,
		buf,
		&bytesWritten,
		nil,
	)
	if err != nil {
		return int(bytesWritten), fmt.Errorf(
			"WriteFile: %w",
			err,
		)
	}

	return int(bytesWritten), nil
}

func (p *windowsSerialPort) Close() error {
	if p == nil || p.handle == 0 {
		return nil
	}

	handle := p.handle
	p.handle = 0

	return windows.CloseHandle(handle)
}
