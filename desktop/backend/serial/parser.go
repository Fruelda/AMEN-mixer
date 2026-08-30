package serial

import (
	"errors"
	"strconv"
	"strings"

	"desktop/backend/protocol"
)

func Parse(line string) (*protocol.Command, error) {

	line = strings.TrimSpace(line)

	if line == "" {
		return nil, errors.New("empty packet")
	}

	part := strings.Split(line, ",")

	switch part[0] {

	case "ENC":

		if len(part) != 3 {
			return nil, errors.New("invalid encoder packet")
		}

		channel, err := strconv.Atoi(part[1])

		if err != nil {
			return nil, err
		}

		delta, err := strconv.Atoi(part[2])

		if err != nil {
			return nil, err
		}

		return &protocol.Command{

			Type: "ENC",

			Channel: channel,

			Value: delta,
		}, nil

	case "BTN":

		if len(part) != 3 {
			return nil, errors.New("invalid button packet")
		}

		channel, err := strconv.Atoi(part[1])

		if err != nil {
			return nil, err
		}

		value, err := strconv.Atoi(part[2])

		if err != nil {
			return nil, err
		}

		// Firmware in this repository emits BTN,<channel>,0 on a press,
		// while the desktop protocol defines 1 as "pressed". Normalize the
		// existing firmware packet here so the rest of the app keeps one
		// consistent command contract.
		if value == 0 {
			value = 1
		}

		return &protocol.Command{

			Type: "BTN",

			Channel: channel,

			Value: value,
		}, nil

	}

	return nil, errors.New("unknown packet")

}
