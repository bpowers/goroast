// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// FIXME: these definitions are from kernel headers.  I either need to
// revise them, or relicense under the GPL2.

package spi

import (
	"fmt"
	"os"
	"unsafe"
)

const (
	SPI_CPHA = 0x01
	SPI_CPOL = 0x02

	SPI_MODE_0 = (0 | 0)
	SPI_MODE_1 = (0 | SPI_CPHA)
	SPI_MODE_2 = (SPI_CPOL | 0)
	SPI_MODE_3 = (SPI_CPOL | SPI_CPHA)

	SPI_CS_HIGH   = 0x04
	SPI_LSB_FIRST = 0x08
	SPI_3WIRE     = 0x10
	SPI_LOOP      = 0x20
	SPI_NO_CS     = 0x40
	SPI_READY     = 0x80

	SPI_IOC_MAGIC = 'k'
)

type SPIIOTransaction struct {
	TXBuf       uint64
	RXBuf       uint64
	Len         uint32
	SpeedHz     uint32
	DelayUsecs  uint16
	BitsPerWord uint8
	CSChange    uint8
	Pad         uint32
}
const sizeof_SPIIOTransaction = 32

var (
	SPI_IOC_RD_MODE = IOR(SPI_IOC_MAGIC, 1, 1)
	SPI_IOC_WR_MODE = IOW(SPI_IOC_MAGIC, 1, 1)

	// Read / Write SPI bit justification
	SPI_IOC_RD_LSB_FIRST = IOR(SPI_IOC_MAGIC, 2, 1)
	SPI_IOC_WR_LSB_FIRST = IOW(SPI_IOC_MAGIC, 2, 1)

	// Read / Write SPI device word length (1..N)
	SPI_IOC_RD_BITS_PER_WORD = IOR(SPI_IOC_MAGIC, 3, 1)
	SPI_IOC_WR_BITS_PER_WORD = IOW(SPI_IOC_MAGIC, 3, 1)

	// Read / Write SPI device default max speed hz
	SPI_IOC_RD_MAX_SPEED_HZ = IOR(SPI_IOC_MAGIC, 4, 4)
	SPI_IOC_WR_MAX_SPEED_HZ = IOW(SPI_IOC_MAGIC, 4, 4)
)

func SPI_IOC_MESSAGE(count int) int32 {
	return IOW(SPI_IOC_MAGIC, 0, count*sizeof_SPIIOTransaction)
}

func Transaction(f *os.File, write, read []byte) error {
	if write != nil && read != nil {
		if len(write) != len(read) {
			return fmt.Errorf("write and read size mismatch (%d vs %d)",
				len(write), len(read))
		}
	}
	var length uint32
	if write != nil {
		length = uint32(len(write))
	} else {
		length = uint32(len(read))
	}
	trx := SPIIOTransaction{
		TXBuf: uint64(uintptr(unsafe.Pointer(&write[0]))),
		RXBuf: uint64(uintptr(unsafe.Pointer(&write[0]))),
		Len:   length,
	}
	return ioctl(f.Fd(), SPI_IOC_MESSAGE(1), unsafe.Pointer(&trx))
}
