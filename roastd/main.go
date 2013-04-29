// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/bpowers/goroast/devices"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const (
	usage = `Usage: %s [OPTION...]
IO daemon for SR500 coffee roaster controller.

Options:
`
)

var (
	memProfile string
	cpuProfile string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&memProfile, "memprofile", "",
		"write memory profile to this file")
	flag.StringVar(&cpuProfile, "cpuprofile", "",
		"write cpu profile to this file")

	flag.Parse()
}

// startProfiling enables memory and/or CPU profiling if the
// appropriate command line flags have been set.
func startProfiling() {
	var err error
	// if we've passed in filenames to dump profiling data too,
	// start collecting profiling data.
	if memProfile != "" {
		runtime.MemProfileRate = 1
	}
	if cpuProfile != "" {
		var f *os.File
		if f, err = os.Create(cpuProfile); err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func stopProfiling() {
	if memProfile != "" {
		runtime.GC()
		f, err := os.Create(memProfile)
		if err != nil {
			log.Println(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
	if cpuProfile != "" {
		pprof.StopCPUProfile()
		cpuProfile = ""
	}
}

func main() {
	// if -memprof or -cpuprof haven't been set on the command
	// line, these are nops
	startProfiling()
	defer stopProfiling()

	tc1, err := devices.NewMax31855("/dev/spidev0.0")
	if err != nil {
		fmt.Printf("error: devices.NewMax31855('/dev/spidev0.0'): %s\n", err)
		return
	}
	defer tc1.Close()

	pin, err := os.Create("/sys/class/gpio/gpio22/value")
	if err != nil {
		log.Fatalf("open failed: %s", err)
	}
	defer func() {
		fmt.Printf("shutting heater off\n")
		pin.Write([]byte{'0'})
	}()

	timer := time.Tick(500 * time.Millisecond)

	pin.Write([]byte{'1'})

	for {
		<-timer

		temp, err := tc1.Read()
		if err != nil {
			fmt.Printf("error: tc1.Read(): %s\n", err)
			break
		}
		fmt.Printf("temp: %.2f°C (%.2f°F)\n", temp, temp*1.8+32)
		if temp > 66 {
			fmt.Printf("temp: threshold hit, cooling down\n")
			break
		}
	}

	pin.Write([]byte{'0'})

	for {
		<-timer

		temp, err := tc1.Read()
		if err != nil {
			fmt.Printf("error: tc1.Read(): %s\n", err)
			break
		}
		fmt.Printf("temp: %.2f°C (%.2f°F)\n", temp, temp*1.8+32)
	}
}
