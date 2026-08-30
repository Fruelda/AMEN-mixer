package serial

// serialPort is the tiny transport surface Manager needs.
// Platform-specific files implement opening/configuring
// the underlying port.
type serialPort interface {
	Read([]byte) (int, error)
	Write([]byte,) (int, error)
	Close() error
}