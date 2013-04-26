// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package devices

import (
	"fmt"
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

func (m *Max31855) Read() (Celsius, error) {
	return 0, fmt.Errorf("not implemented")
}

func (m *Max31855) Close() {
	m.f.Close()
}
