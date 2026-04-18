package main

import (
	"log"

	// "go.bug.st/serial"
	"github.com/fatih/color"
)

// For the constants, mapping protocol to a uint8.
type ProtocolType uint8

// Each pair of three bits maps to a channel, from LSB to MSB.
// There are 4 groups of 3 bits, which corresponds
// to the pins of that protocol. For each protocol:
//   - UART: tx=1,   rx=2
//   - SPI:  miso=1, mosi=2, clk=3, cs=4
//   - I2C:  sda=1,  scl=2
//   - CAN:  canh=1, canl=2
//
// Unused values are set to 0.
//
// For example: SPI with respective pins on channels 4, 3, 5, 1:
//   - 0b0100,0011,0101,0001
//     4    3    5    1
//   - 0x4351
//
// Order matters! So if cli argument is --port scl7sda1, the
// correct order is sda then scl, so the ProtocolPins would be
//   - 0b0001,0111,0000,0000
//     1    7    0    0
//   - 0x1700
//
// It is used internally only. Not exported to python script.
type ProtocolPins uint16

var (
	errLog  = color.New(color.FgRed).SprintFunc()
	statLog = color.New(color.FgGreen).SprintFunc()
)

// Maps a protocol to an int
const (
	NONE ProtocolType = iota
	UART
	SPI
	I2C
	CAN
)

// Config struct, specifies which port, duration, which protocol
// to decode (if any), and if so, which channels to use.
type config struct {
	port     string
	duration float32
	protocol ProtocolType
	pins     ProtocolPins
}

func main() {
	config, err := getConfig()

	if err != nil {
		log.Fatal(err)
	}

	config.print()
}

// Print a configuration. For debugging purposes.
func (a *config) print() {
	preamble := "[print] "

	log.Printf(statLog(preamble) + "Current configuration:\n")
	log.Printf(statLog(preamble) + "  Port:\t%s\n", a.port)
	log.Printf(statLog(preamble) + "  Duration:\t%f\n", a.duration)
	log.Printf(statLog(preamble) + "  Protocol:\t%d\n", a.protocol)
	log.Printf(statLog(preamble) + "  Pins:\t0x%X\n", a.pins)
}
