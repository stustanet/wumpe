// Copyright 2017 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type githubRequest struct {
	Ref  string `json:"ref"`
	Hook struct {
		Events []string `json:"events"`
		Config struct {
			Secret string `json:"secret"`
		} `json:"config"`
	} `json:"hook"`
	Repository struct {
		Name string `json:"name"`
	}
}

func parseGitHubRequest(req *http.Request) (h hook, status int) {
	event := req.Header.Get("X-GitHub-Event")
	if event != "push" {
		log.Println("wrong event kind:", event)
		status = http.StatusBadRequest
		return
	}

	signature := req.Header.Get("X-Hub-Signature")
	// signature must have 5 bytes prefix + 40 bytes SHA1 HMAC
	if len(signature) != 45 || !strings.HasPrefix(signature, "sha1=") {
		log.Println("signature missing")
		status = http.StatusBadRequest
		return
	}

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("reading payload error:", err)
		status = http.StatusBadRequest
		return
	}

	// try to parse request to JSON
	var jr githubRequest
	err = json.Unmarshal(payload, &jr)
	if err != nil {
		log.Println("decode error:", err)
		status = http.StatusBadRequest
		return
	}

	// check if a handler exists for this project
	h, ok := cfg.Hooks[jr.Repository.Name]
	if !ok {
		log.Println("hook not found:", jr.Repository.Name)
		status = http.StatusTeapot
		return
	}

	// check if the ref matches
	if jr.Ref != h.Ref {
		log.Println("ref mismatch:", h.Ref)
		status = http.StatusTeapot
		return
	}

	// decode signature from hex to binary format
	var buf [20]byte
	receivedMAC := buf[:]
	hex.Decode(receivedMAC, []byte(signature[5:]))

	// check if the secret matches by computing the HMAC for the given body with
	// the secret and compare it with the signature in the request.
	mac := hmac.New(sha1.New, []byte(h.Secret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(expectedMAC, receivedMAC) {
		status = http.StatusForbidden
		return
	}

	status = http.StatusOK
	return
}
