// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// FIXME: these definitions are from kernel headers.  I either need to
// revise them, or relicense under the GPL2.

package spi

const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8

	// XXX: sometimes overridden, but not that I can tell for SPI
	// on ARM.
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRMASK   = ((1 << _IOC_NRBITS) - 1)
	_IOC_TYPEMASK = ((1 << _IOC_TYPEBITS) - 1)
	_IOC_SIZEMASK = ((1 << _IOC_SIZEBITS) - 1)
	_IOC_DIRMASK  = ((1 << _IOC_DIRBITS) - 1)

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = (_IOC_NRSHIFT + _IOC_NRBITS)
	_IOC_SIZESHIFT = (_IOC_TYPESHIFT + _IOC_TYPEBITS)
	_IOC_DIRSHIFT  = (_IOC_SIZESHIFT + _IOC_SIZEBITS)

	_IOC_NONE  = 0
	_IOC_WRITE = 1
	_IOC_READ  = 2
)

func IOC(dir, ty, nr, size int) int32 {
	return (int32)((dir << _IOC_DIRSHIFT) |
		(ty << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) |
		(size << _IOC_SIZESHIFT))
}

func IOR(ty, nr, size int) int32 {
	return IOC(_IOC_READ, ty, nr, size)
}

func IOW(ty, nr, size int) int32 {
	return IOC(_IOC_WRITE, ty, nr, size)
}
