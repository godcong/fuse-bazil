package fuse

import (
	"syscall"
)

const (
	ENODATA = Errno(syscall.ENODATA)
)

const (
	errNoXAttr = ENODATA
)

func init() {
	errnoNames[errNoXAttr] = "ENODATA"
}
