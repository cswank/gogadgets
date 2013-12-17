package devices

import (
	"fmt"
	"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
	"os"
	"errors"
	"io/ioutil"
)

type GPIO struct {
	OutputDevice
	InputDevice
	units string
	export string
	exportPath string
	directionPath string
	valuePath string
	direction string
	edge string
}

func NewGPIO(pin *models.Pin) (*GPIO, error) {
	portMap, ok := Pins["gpio"][pin.Port]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such port: %s", pin.Port))
	}
	export, ok := portMap[pin.Pin]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such pin: %s", pin.Pin))
	}
	g := &GPIO{
		export: export,
		exportPath: "/sys/class/gpio/export",
		directionPath: fmt.Sprintf("/sys/class/gpio/gpio%s/direction", export),
		valuePath: fmt.Sprintf("/sys/class/gpio/gpio%s/value", export),
		direction: pin.Direction,
		edge: pin.Edge,
	}
	err := g.Init()
	return g, err
}

func (g *GPIO) Init() error {
	var err error
	if !utils.FileExists(g.directionPath) {
		err = g.writeValue(g.exportPath, g.export)
	}
	if err == nil {
		err = g.writeValue(g.directionPath, "out")
		if err == nil {
			err = g.writeValue(g.valuePath, "0")
		}
	}
	return err
}

func (g *GPIO) Update(msg *models.Message) {
	
}

func (g *GPIO) On(val *models.Value) error {
	return g.writeValue(g.valuePath, "1")
}

func (g *GPIO) Status() bool {
	data, err := ioutil.ReadFile(g.valuePath)
	return err == nil && string(data) == "1\n"
}

func (g *GPIO) Off() error {
	return g.writeValue(g.valuePath, "0")
}

func (g *GPIO) writeValue(path, value string) error {
	return ioutil.WriteFile(path, []byte(value), os.ModeDevice)
}

/*
fd_set exceptfds;
int    res;    

FD_ZERO(&exceptfds);
FD_SET(gpioFileDesc, &exceptfds);

res = select(gpioFileDesc+1, 
             NULL,               // readfds - not needed
             NULL,               // writefds - not needed
             &exceptfds,
             NULL);              // timeout (never)

if (res > 0 && FD_ISSET(gpioFileDesc, &exceptfds))
{
     // GPIO line changed
}
*/
func (g *GPIO) Wait(value interface{}) error {
	// func Open(path string, mode int, perm uint32) (fd int, err error)
	// func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error)
	// fd, err := syscall.Open(g.valuePath, syscall.O_RDONLY, 0777)
	// if err != nil {
	// 	return err
	// }
	// fdSet = &syscall.FdSet{fd}
	// n, err := syscall.Select(fd + 1, nil, nil, fdSet, nil)
	return nil
}
