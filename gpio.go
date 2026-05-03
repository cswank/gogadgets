//go:build !windows
// +build !windows

package gogadgets

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// Kept for test compatibility.
var (
	GPIO_DEV_PATH = "/sys/class/gpio"
	GPIO_DEV_MODE = os.ModeDevice
)

// GPIO interacts with the Linux GPIO character device interface
// (/dev/gpiochipN) to control pins. The Chip field on Pin selects
// the gpiochip device (default "0"). Line offsets come from the
// existing pin maps (Pins and PiPins).
type GPIO struct {
	line      uint32
	direction string
	activeLow bool
	edge      string
	fd        int // line handle fd (output) or event fd (input)
}

// chardev ioctl structures (match kernel UAPI gpio.h)
type gpiohandleRequest struct {
	LineOffsets   [64]uint32
	Flags        uint32
	DefaultValues [64]uint8
	ConsumerLabel [32]byte
	Lines        uint32
	Fd           int32
}

type gpiohandleData struct {
	Values [64]uint8
}

type gpioeventRequest struct {
	LineOffset    uint32
	HandleFlags   uint32
	EventFlags    uint32
	ConsumerLabel [32]byte
	Fd            int32
}

const (
	gpioMagic = 0xB4

	gpioHandleRequestInput     = 1 << 0
	gpioHandleRequestOutput    = 1 << 1
	gpioHandleRequestActiveLow = 1 << 2

	gpioEventRequestRisingEdge  = 1 << 0
	gpioEventRequestFallingEdge = 1 << 1
	gpioEventRequestBothEdges   = gpioEventRequestRisingEdge | gpioEventRequestFallingEdge
)

func ioc(dir, typ, nr, size uintptr) uintptr {
	return (dir << 30) | (size << 16) | (typ << 8) | nr
}

var (
	ioctlGetLineHandle = ioc(3, gpioMagic, 0x03, unsafe.Sizeof(gpiohandleRequest{}))
	ioctlGetLineEvent  = ioc(3, gpioMagic, 0x04, unsafe.Sizeof(gpioeventRequest{}))
	ioctlGetLineValues = ioc(3, gpioMagic, 0x08, unsafe.Sizeof(gpiohandleData{}))
	ioctlSetLineValues = ioc(3, gpioMagic, 0x09, unsafe.Sizeof(gpiohandleData{}))
)

func gpioIoctl(fd int, req uintptr, arg unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(arg))
	if errno != 0 {
		return errno
	}
	return nil
}

func NewGPIO(pin *Pin) (OutputDevice, error) {
	return newGPIO(pin)
}

func newGPIO(pin *Pin) (*GPIO, error) {
	var lineStr string
	var ok bool
	if pin.Platform == "rpi" {
		lineStr, ok = PiPins[pin.Pin]
		if !ok {
			return nil, fmt.Errorf("no such pin: %s", pin.Pin)
		}
	} else {
		var portMap map[string]string
		portMap, ok = Pins["gpio"][pin.Port]
		if !ok {
			return nil, fmt.Errorf("no such port: %s", pin.Port)
		}
		lineStr, ok = portMap[pin.Pin]
		if !ok {
			return nil, fmt.Errorf("no such pin: %s", pin.Pin)
		}
	}

	lineNum, err := strconv.Atoi(lineStr)
	if err != nil {
		return nil, fmt.Errorf("invalid line number %q: %w", lineStr, err)
	}

	if pin.Direction == "" {
		pin.Direction = "out"
	}

	chip := pin.Chip
	if chip == "" {
		chip = "0"
	}

	g := &GPIO{
		line:      uint32(lineNum),
		direction: pin.Direction,
		activeLow: pin.ActiveLow == "1",
		edge:      pin.Edge,
	}

	chipPath := fmt.Sprintf("/dev/gpiochip%s", chip)
	chipFd, err := syscall.Open(chipPath, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", chipPath, err)
	}
	defer syscall.Close(chipFd)

	if g.direction == "out" {
		err = g.requestOutput(chipFd)
	} else {
		err = g.requestEvent(chipFd)
	}
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GPIO) requestOutput(chipFd int) error {
	req := gpiohandleRequest{
		Lines: 1,
		Flags: gpioHandleRequestOutput,
	}
	req.LineOffsets[0] = g.line
	copy(req.ConsumerLabel[:], "gogadgets")
	if g.activeLow {
		req.Flags |= gpioHandleRequestActiveLow
	}
	if err := gpioIoctl(chipFd, ioctlGetLineHandle, unsafe.Pointer(&req)); err != nil {
		return fmt.Errorf("request line handle for line %d: %w", g.line, err)
	}
	g.fd = int(req.Fd)
	return nil
}

func (g *GPIO) requestEvent(chipFd int) error {
	req := gpioeventRequest{
		LineOffset:  g.line,
		HandleFlags: gpioHandleRequestInput,
	}
	copy(req.ConsumerLabel[:], "gogadgets")
	if g.activeLow {
		req.HandleFlags |= gpioHandleRequestActiveLow
	}
	switch g.edge {
	case "rising":
		req.EventFlags = gpioEventRequestRisingEdge
	case "falling":
		req.EventFlags = gpioEventRequestFallingEdge
	case "both", "":
		req.EventFlags = gpioEventRequestBothEdges
	default:
		return fmt.Errorf("unknown edge type: %s", g.edge)
	}
	if err := gpioIoctl(chipFd, ioctlGetLineEvent, unsafe.Pointer(&req)); err != nil {
		return fmt.Errorf("request line event for line %d: %w", g.line, err)
	}
	g.fd = int(req.Fd)
	return nil
}

func (g *GPIO) Commands(location, name string) *Commands {
	return nil
}

func (g *GPIO) Update(msg *Message) bool {
	return false
}

func (g *GPIO) On(val *Value) error {
	data := gpiohandleData{}
	data.Values[0] = 1
	return gpioIoctl(g.fd, ioctlSetLineValues, unsafe.Pointer(&data))
}

func (g *GPIO) Off() error {
	data := gpiohandleData{}
	return gpioIoctl(g.fd, ioctlSetLineValues, unsafe.Pointer(&data))
}

func (g *GPIO) Status() map[string]bool {
	data := gpiohandleData{}
	err := gpioIoctl(g.fd, ioctlGetLineValues, unsafe.Pointer(&data))
	return map[string]bool{"gpio": err == nil && data.Values[0] != 0}
}

// Wait blocks until a GPIO event (edge transition) occurs on the line.
func (g *GPIO) Wait() error {
	var buf [16]byte
	_, err := syscall.Read(g.fd, buf[:])
	return err
}

func (g *GPIO) Close() error {
	if g.fd > 0 {
		return syscall.Close(g.fd)
	}
	return nil
}
