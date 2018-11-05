package messages

// Update is produced by any device in the system to notify
// everything about a state change.
type Update struct {
	Location string
	Name     string
}

// Command is broadcast all the output devices in the system.  If a device
// decides it applys to itself then it does something (turns on a
// gpio, for example).
type Command struct {
}
