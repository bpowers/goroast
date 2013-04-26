// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"text/template"
)

type statusHandler struct {
	http.Handler
}

func (h *statusHandler) serveStatus(w http.ResponseWriter, req *http.Request) {
	managerTemplate := template.Must(template.ParseFiles("./tmpl/status.html"))
	data := "&lt;no data&gt;"
	var out bytes.Buffer

	cmd := exec.Command("/usr/local/bin/roastd")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("exec.Command('/usr/local/bin/roastd'): %s\n", err)
		goto out
	}
	data = out.String()

out:
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	managerTemplate.Execute(w, data)
}

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		h.serveStatus(w, req)
		return
	}

	h.Handler.ServeHTTP(w, req)
}
