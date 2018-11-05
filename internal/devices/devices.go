package devices

import "github.com/cswank/gogadgets/internal/messages"

// Output can be turned on or off.  When it is on it does something.
type Output interface {
	On() error
	Off() error
	Update(msg messages.Update) (*messages.Update, error)
}

// Input is for devices (sensors) that read input and report it back to the system.
type Input interface {
	// Start lets the device begin reading its physical device and reporting the
	// results by calling the update callback.
	Start(update func(messages.Update)) error
	Get() messages.Update
}
