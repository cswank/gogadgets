//
// Copyright 2014-2016 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package serial // import "go.bug.st/serial.v1"

import "syscall"

const devFolder = "/dev"
const regexFilter = "^(cu|tty)\\..*"

// termios manipulation functions

var baudrateMap = map[int]int{
	0:      syscall.B9600, // Default to 9600
	50:     syscall.B50,
	75:     syscall.B75,
	110:    syscall.B110,
	134:    syscall.B134,
	150:    syscall.B150,
	200:    syscall.B200,
	300:    syscall.B300,
	600:    syscall.B600,
	1200:   syscall.B1200,
	1800:   syscall.B1800,
	2400:   syscall.B2400,
	4800:   syscall.B4800,
	9600:   syscall.B9600,
	19200:  syscall.B19200,
	38400:  syscall.B38400,
	57600:  syscall.B57600,
	115200: syscall.B115200,
	230400: syscall.B230400,
}

var databitsMap = map[int]int{
	0: syscall.CS8, // Default to 8 bits
	5: syscall.CS5,
	6: syscall.CS6,
	7: syscall.CS7,
	8: syscall.CS8,
}

const tcCMSPAR int = 0 // may be CMSPAR or PAREXT
const tcIUCLC int = 0

// syscall wrappers

//sys ioctl(fd int, req uint64, data uintptr) (err error)

const ioctlTcgetattr = syscall.TIOCGETA
const ioctlTcsetattr = syscall.TIOCSETA
