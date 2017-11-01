// Copyright 2017 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/naoina/toml"
)

type config struct {
	Listen string
	Hooks  map[string]struct {
		Ref    string
		Secret string
		Dir    string
		Cmd    string
	}
}

type request struct {
	ObjectKind string `json:"object_kind"`
	Ref        string `json:"ref"`
	Project    struct {
		Name string `json:"name"`
	}
}

var cfg config

func sendErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func Build(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" || req.URL.Path != "/build" {
		sendErr(w, http.StatusTeapot)
		return
	}

	var jr request
	jd := json.NewDecoder(req.Body)
	err := jd.Decode(&jr)
	if err != nil {
		sendErr(w, http.StatusBadRequest)
		return
	}

	if jr.ObjectKind != "push" {
		sendErr(w, http.StatusBadRequest)
		return
	}

	hook, ok := cfg.Hooks[jr.Project.Name]
	if !ok {
		sendErr(w, http.StatusTeapot)
		return
	}

	if jr.Ref != hook.Ref || req.Header.Get("X-Gitlab-Token") != hook.Secret {
		sendErr(w, http.StatusForbidden)
		return
	}

	cmd := exec.Command("/usr/bin/git", "pull")
	cmd.Dir = hook.Dir
	out, err := cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		sendErr(w, http.StatusInternalServerError)
		return
	}
	out = bytes.TrimSpace(out)
	if bytes.HasPrefix(out, []byte("Already up-to-date")) {
		// no new commits
		return
	}

	cmd = exec.Command(hook.Cmd)
	cmd.Dir = hook.Dir
	out, err = cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		sendErr(w, http.StatusInternalServerError)
		return
	}
}

func main() {
	f, err := os.Open("/etc/wumpe.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	http.HandleFunc("/", Build)
	log.Fatal(http.ListenAndServe(cfg.Listen, nil))
}
