// Copyright 2017 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

var cfg config

func sendErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func Build(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" || req.URL.Path != "/build" {
		sendErr(w, http.StatusTeapot)
		return
	}

	var project string
	var secret string

	// try to detect which type of request (GitLab or GitHub) this is
	switch {
	case strings.HasPrefix(req.UserAgent(), "GitHub-Hookshot/"):
		project, secret = parseGitHubRequest(req)
	case req.Header.Get("X-Gitlab-Token") != "":
		project, secret = parseGitLabRequest(req)
	}
	if project == "" || secret == "" {
		sendErr(w, http.StatusBadRequest)
		return
	}

	// check if a handler exists for this project
	hook, ok := cfg.Hooks[project]
	if !ok {
		sendErr(w, http.StatusTeapot)
		return
	}

	// check if the secret matches
	if secret != hook.Secret {
		sendErr(w, http.StatusForbidden)
		return
	}

	// pull the updates
	cmd := exec.Command("/usr/bin/git", "pull")
	cmd.Dir = hook.Dir
	out, err := cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		sendErr(w, http.StatusInternalServerError)
		return
	}
	out = bytes.TrimSpace(out)
	var gitUnchanged = []byte("Already up-to-date.")
	if bytes.HasPrefix(out, gitUnchanged) {
		// no new commits
		w.WriteHeader(http.StatusConflict)
		w.Write(gitUnchanged)
		return
	}

	// if new commits were pulled, call the hook command
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
