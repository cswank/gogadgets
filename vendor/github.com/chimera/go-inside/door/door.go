package door

import (
	"fmt"
	"log"

	"code.google.com/p/gopass"

	"github.com/chimera/go-inside/rs232"
	"github.com/chimera/go-inside/users"
)

type DoorLock struct {
	Serial        rs232.Port
	Baud          int
	UsersFile     string
	PossiblePorts []string // Possible ports that the Arduino could be on...
}

func (d *DoorLock) Unlock() error {

	// Loop over the available ports and try to connect in order.
	for _, port := range d.PossiblePorts {

		// Configure the serial connection.
		options := rs232.Options{
			BitRate:  uint32(d.Baud),
			DataBits: 8,
			StopBits: 1,
			Parity:   rs232.PARITY_NONE,
			Timeout:  0,
		}

		// Open a connection to the serial port.
		p, err := rs232.Open(port, options)

		// Handle connection errors.
		if err != nil {
			// log.Fatalf("Failed to connect to serial port with error code: %d (%s)\n", e.Code, errType)
			// return fmt.Errorf("Failed to connect to serial port with error code: %d (%s)\n", e.Code, errType)
			log.Printf("Could not connect to port %s\n", port)
			continue
		}

		// Attach a reference of the serial port to the DoorLock struct.
		d.Serial = *p

		log.Printf("Opened serial port %s\n", d.Serial.String())

		// Unlock door
		_, err = d.Serial.Write([]byte("1"))
		if err != nil {
			return fmt.Errorf("Could not unlock door, with error: %s", err)
		}
		log.Println("Door unlocked!")

		// If we successfully connect and write to the port, finish the loop.
		return nil
	}

	// None of the expected ports could be connected to.
	return fmt.Errorf("Failed to connect to available ports!")
}

func (d *DoorLock) Disconnect() {
	log.Println("Disconnecting from door lock...")
	d.Serial.Close()
}

func (d *DoorLock) Listen() {


	// Listen for incoming RFID codes.
	for {

		// Prompt for an RFID code.
		var code string
		code, err := gopass.GetPass("Please input your RFID code for access: ")
		if err != nil {
			log.Fatal(err)
		}

		// If a code is received, send it to get authenticated.
		if code != "" {
			err = users.AuthenticateCode(code, d.UsersFile)
			if err != nil {
				log.Println(err.Error())
			} else {
				// Log them in if their code is valid.
				err := d.Unlock()
				if err != nil {
					log.Fatal(err)
				}
				defer d.Disconnect()
			}
		}
	}
}
