// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package devices

import (
	"fmt"
	"github.com/bpowers/goroast/spi"
	"os"
)

type Celsius float64

type Max31855 struct {
	f *os.File
}

func NewMax31855(path string) (*Max31855, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile('%s', os.O_RDWR, 0)", path)
	}
	return &Max31855{f}, nil
}

// Reads are 4-bytes
func (m *Max31855) Read() (Celsius, error) {
	buf := make([]byte, 4)

	if err := spi.Transaction(m.f, nil, buf); err != nil {
		return 0, fmt.Errorf("spi.Transaction(%v, nil, buf): %s", m.f, err)
	}

	return Celsius(.25 * float64(buf[0] << 6 | buf[1] >> 2)), nil
}

func (m *Max31855) Close() {
	m.f.Close()
}
