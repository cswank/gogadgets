// The main package handles connecting to the door lock and receiving various
// command line flags, including the JSON users database file, serial port location
// and the baud rate.
package main

import (
	"flag"
	"github.com/chimera/go-inside/door"
)

// Available command line flags with sane-ish defaults.
var users_file = flag.String("db", "/home/pi/users.json", "The users JSON file to use.")
var baud = flag.Int("baud", 19200, "The baudrate to connect to the serial port with.")

func main() {
	// Parse any command line flags.
	flag.Parse()

	// Create a new connection to the door lock
	door := &door.DoorLock{
		Baud:      *baud,
		UsersFile: *users_file,
		PossiblePorts: []string{
			"/dev/ttyACM0",
			"/dev/ttyACM1",
			"/dev/ttyACM2",
			"/dev/tty.usbmodem411",
			"/dev/tty.usbmodem621",
		},
	}

	// Handle inputting of user RFID codes
	door.Listen()

	// Make sure to disconnect from the door when we're done.
	defer door.Disconnect()
}
