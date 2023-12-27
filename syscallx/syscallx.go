// Package syscallx used to provide wrappers for syscalls.
//
// It is no longer needed.
package syscallx

import "golang.org/x/sys/unix"

// Wrappers for xattr syscalls.
var (
	GetXAttr    = unix.Getxattr
	ListXAttr   = unix.Listxattr
	SetXAttr    = unix.Setxattr
	RemoveXAttr = unix.Removexattr
	Msync       = unix.Msync
)
