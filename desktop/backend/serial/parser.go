package serial

import (
	"errors"
	"strconv"
	"strings"

	"desktop/backend/protocol"
)

// ============================================================
// PARSE COMMAND
// ============================================================

func Parse(
	line string,
) (*protocol.Command, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, errors.New("empty packet")
	}

	parts := strings.Split(line, ",")

	if len(parts) != 3 {
		return nil, errors.New("invalid packet")
	}

	// ========================================================
	// COMMAND TYPE
	// ========================================================

	var commandType protocol.CommandType

	switch parts[0] {
	case string(protocol.CommandEncoder):
		commandType = protocol.CommandEncoder

	case string(protocol.CommandButton):
		commandType = protocol.CommandButton

	default:
		return nil, errors.New("unknown packet")
	}

	// ========================================================
	// CHANNEL
	// ========================================================

	channel, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	// ========================================================
	// VALUE
	// ========================================================

	value, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	// ========================================================
	// COMMAND
	// ========================================================

	return &protocol.Command{
		Type:    commandType,
		Channel: channel,
		Value:   value,
	}, nil
}
