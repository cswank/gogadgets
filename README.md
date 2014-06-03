# GoGadgets


## Gadgets

A gogadgets app consists of a collection of gadgets.  Most gadgets
control a physical device (led, thermometer, motion sensor, electric
motor, etc).  The gogadgets includes 5 built in gadgets::

### GPIO
This gadget is used to turn things on and off.

### Switch
Switch is really a GPIO that has been configured as input.  It is used
to wait for some input device (push button, motion sensor, etc) to change
state.  When the state of the input device changes then the rest of the
system receives a message for that state change.  All gadgets, whether
input gadgets or output gadgets, received these messages and can therefore
react to changes in other parts of the system.

### Thermometer
The 'extras' directory includes a compiled dts file (BB-W1-00A0.dtbo) for dallas 
one-wire devices.  If you load it into the beaglebone's device tree overlay then
you can connect a dallas 1-wire temperature sensor to it and include the thermometer 
in your gadgets system.

### Heater
Use this gadget to heat things up.  When you heat something up with Heater it listens
for temperature updates and turns itself off when the desired temperature has been
reached.  You must have a thermometer in your system to use this.

### Cooler

### Motor

### Recorder
This gadget doesn't actually control hardware.  It receives all the update messages
from the rest of the system and pushes them to a MongoDB (which should really be
hosted on a server somewhere on your sub-net).

## Installation

    # apt-get update
    # apt-get install build-essential
    # apt-get install libzmq3-dev
    # apt-get install git
    # apt-get install bzr
    # apt-get install golang
    # mkdir /opt/gadgets
    # export GOPATH=/opt/gadgets/
    # go get bitbucket.org/cswank/gogadgets
    # cd ./src/bitbucket.org/cswank/gogadgets/gogadgets
    # go install

## Try out one of the examples

The config file for the led example looks like::

    {
        "gadgets": [
            {
                "location": "lab",
                "name": "led",
                "pin": {
                    "type": "gpio",
                    "port": "8",
                    "pin": "9",
                }
            }
        ]
    }

So to try this example out you would connect an led to port 8, pin 9.  Next, start the app::

    # $GOPATH/bin/gogadgets -g $GOPATH/src/bitbucket.org/cswank/gogadgets/examples/led/config.json

Next you can turn on the led from the command line.  You can either ssh into the same beaglebone
that is running gogagdgets, or you can install gadgets again on another machine.  If you were
to install gogadgets on a remote machine then you can turn on the led like this::

    $ $GOPATH/bin/gogadgets -c "turn on lab led"

## Robot Command Language


All Gadgets in the system respond to Robot Command Language (RCL) messages.
A RCL command is a string consisting of::

      <action> <location> <device name> <for|to> <command argument> <units>

So, for example::

    turn on living room light
    turn on living room light for 30 minutes
    heat boiler to 22 C

## Methods

A Method is a sequence of RCL messages.  Let's say you had a gadgets system like this::

    {
        "gadgets": [
            {
                "location": "lab",
                "name": "led",
                "pin": {
                    "type": "gpio",
                    "port": "8",
                    "pin": "9",
                }
            },
            {
                "location": "lab",
                "name": "motion sensor",
                "pin": {
                    "type": "switch",
                    "port": "8",
                    "pin": "11",
                }
            }
        ]
    }
