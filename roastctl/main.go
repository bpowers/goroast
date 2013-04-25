// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

const (
	usage = `Usage: %s [OPTION...]
control program for coffee roasting

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

	// TODO: control stuff
}
