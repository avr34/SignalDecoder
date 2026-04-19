package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"

	"go.bug.st/serial"
)

// Performs a logical xor of two bools.
// Funny how the language has no native support
// for this.
func xor(A, B bool) bool {
	return (!A && B) || (A && !B)
}

// Initializes all the flags, parses the flags, checks for
// invalidities, and if nothing's invalid, returns a pointer
// to a config struct.
//
// If there IS an error, it returns an error. It will not
// halt the program.
func getConfig() (*config, error) {
	preamble := "[getConfig]: "

	var ret config
	var err error

	port := flag.String("port", "", "Port at which to read data from the Logic Analyzer.")
	pins := flag.String("pins", "", "(Optional) Channels to use to perform protocol decoding.")
	duration := flag.Float64("duration", float64(0x0), "Duration (in ms) to run signal capture.")
	protocol := flag.String("protocol", "", "(Optional) Protocol to decode.")

	flag.Parse()

	log.Print(statLog(preamble) + "Parsed arguments")

	// Handle invalid inputs.
	if *port == "" {
		return nil, errors.New(errLog(preamble) + "Must specify a port.")
	} else {
		ret.port = *port
	}

	if *duration <= 0.0 {
		return nil, errors.New(errLog(preamble) + "Must specify a (positive) duration.")
	} else {
		ret.duration = float32(*duration)
	}

	// why doesn't this language have logical xor!
	// basically, one cannot exist w/o the other.
	if xor(*protocol == "", *pins == "") {
		return nil, errors.New(errLog(preamble) + "Cannot use protocol without pin declaration and vice versa.")
	}

	// Check if port exists.
	ports, err := serial.GetPortsList()

	if err != nil {
		return nil, errors.New(errLog(preamble) + "Issues opening serial port")
	}

	if len(ports) == 0 {
		return nil, errors.New(errLog(preamble) + "No serial port found")
	} else {
		if exists := slices.Contains(ports, *port); !exists {
			return nil, fmt.Errorf(errLog(preamble)+"Port %s not found\n", *port)
		}
	}

	switch strings.ToUpper(*protocol) {
	case "":
		ret.protocol = NONE
	case "UART":
		ret.protocol = UART
	case "SPI":
		ret.protocol = SPI
	case "I2C":
		ret.protocol = I2C
	case "CAN":
		ret.protocol = CAN
	default:
		return nil, errors.New(errLog(preamble) + "Unrecognized comm protocol.")
	}

	ret.pins, err = parsePins(ret.protocol, *pins)
	if err != nil {
		return nil, err
	}

	log.Print(statLog(preamble) + "CLI argument error checking complete.")

	return &ret, nil
}

// Call individual parsing function based on the protocol
func parsePins(p ProtocolType, pins string) (ProtocolPins, error) {
	preamble := "[parsePins]: "

	log.Print(statLog(preamble) + "Parsing pins argument")
	switch p {
	case UART:
		ret, err := parseUART(pins)
		if err != nil {
			return 0, err
		} else {
			return ret, nil
		}
	case SPI:
		ret, err := parseSPI(pins)
		if err != nil {
			return 0, err
		} else {
			return ret, nil
		}
	case I2C:
		ret, err := parseI2C(pins)
		if err != nil {
			return 0, err
		} else {
			return ret, nil
		}
	case CAN:
		ret, err := parseCAN(pins)
		if err != nil {
			return 0, err
		} else {
			return ret, nil
		}
	case NONE:
		log.Print(statLog(preamble) + "No pin configuration provided. Proceeding.")
		return 0, nil
	default:
		return 0, errors.New(errLog(preamble) + "Unknown protocol referenced.")
	}

	log.Print(statLog(preamble) + "Pin parsing complete.")

	return 0, nil
}
