# Go-Inside

A Golang RFID access control system for [Chimera Arts & Maker Space](http://chimeraarts.org).

Intended to run on a Raspberry Pi, connected to an Arduino that controls an electronic door latch.

The Arduino code just listens for the character `"1"` over the serial port, and if it gets it, will trigger the electronic door latch using a MOSFET transistor. 

The `go-inside` binary in this repository is built to work on a Linux ARM operating system, like the RPi. If you need to run this application on other systems, it will need to be rebuilt using `go build`. 

The Arduino code can be found in the `servant` directory.

## Setup

First, you must upload the `servant/servant.pde` sketch onto your Arduino.

Download the `go-inside` binary (for Linux ARM distros) and place in desired directory (home director is just fine).

After that, run the `go-inside` binary:

    $ go-inside -port="/dev/ttyACM0"

The `port` will be the location of your Arduino serial port connection. On Raspberry Pi, it usually starts with `/dev/ttyACM`. On OSX, it is usually at `/dev/tty.usbmodem`.

<http://www.stuffaboutcode.com/2012/06/raspberry-pi-run-program-at-start-up.html>

Create a file at `/etc/init.d/go-inside`:

    #! /bin/sh
    case "$1" in
      start)
        echo "Starting go-inside"
        /usr/local/bin/go-inside
        ;;
      stop)
        echo "Stopping go-inside"
        killall go-inside
        ;;
      *)
        echo "Usage: /etc/init.d/go-inside {start|stop}"
        exit 1
        ;;
    esac

    exit 0

Make the script executable:

    sudo chmod 755 /etc/init.d/go-inside

The register the script to run when the system starts up/shuts down:

    sudo update-rc.d go-inside defaults

## Development

To build the package for ARM based systems (such as the Raspberry Pi) while on another platform (OSX, Windows), run the following:

    GOARCH=arm GOOS=linux go build

## Todo

- Negate the use of the Arduino and use the RPi's GPIO pins instead. Need to use a 3.3v transistor to trigger the door latch in order to accomplish this.

## Credit

Written by [Dana Woodman](http://danawoodman.com) for [Chimera Arts & Maker Space](http://chimeraarts.org). Based on the rs232 library generously shared by jpad.

## License

Licensed under the MIT license.
