// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/bpowers/seshcookie"
	"log"
	"os"
	"net/http"
	"runtime"
	"runtime/pprof"
)

const (
	usage = `Usage: %s [OPTION...]
web dashboard for coffee roasting

Options:
`
)

var (
	memProfile string
	cpuProfile string
	devMode    bool
	cookieKey  string
	cookieName = "roasts"
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
	flag.BoolVar(&devMode, "dev", false, "run on port 8080, rather than 443")
	flag.StringVar(&cookieKey, "cookie-key", "",
		"key for http sessions")

	flag.Parse()

	if cookieKey == "" {

	}
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
	var err error

	// if -memprof or -cpuprof haven't been set on the command
	// line, these are nops
	startProfiling()
	defer stopProfiling()

	rootHandler := seshcookie.NewSessionHandler(
		&AuthHandler{
			http.FileServer(http.Dir("./static")),
			&authorizer{"."},
			&decider{},
		},
		cookieKey,
		nil)
	rootHandler.CookieName = cookieName

	// TODO: loop and do stuff
	http.Handle("/", rootHandler)
	http.Handle("/err/", http.FileServer(http.Dir("./err")))

	if devMode {
		err = http.ListenAndServe(
			":8080",
			nil)
	} else {
		go func() {
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
				http.Redirect(rw, r, "https://boosd.org/", 302)
			})
			http.ListenAndServe(":80", mux)
		}()
		// if we're serving over https, set the secure flag
		// for cookies
		seshcookie.Session.Secure = true
		err = http.ListenAndServeTLS(
			":443",
			"/home/bpowers/.tls/certchain.pem",
			"/home/bpowers/.tls/boosd.org_key.pem",
			nil)
	}
	if err != nil {
		log.Printf("ListenAndServe:", err)
	}

}
