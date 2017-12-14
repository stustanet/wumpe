// Copyright 2017 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
)

type gitlabRequest struct {
	ObjectKind string `json:"object_kind"`
	Ref        string `json:"ref"`
	Project    struct {
		Name string `json:"name"`
	}
}

func parseGitLabRequest(req *http.Request) (h hook, status int) {
	var jr gitlabRequest
	jd := json.NewDecoder(req.Body)
	err := jd.Decode(&jr)
	if err != nil {
		status = http.StatusBadRequest
		return
	}

	if jr.ObjectKind != "push" {
		status = http.StatusBadRequest
		return
	}

	// check if a handler exists for this project
	h, ok := cfg.Hooks[jr.Project.Name]
	if !ok {
		status = http.StatusTeapot
		return
	}

	// check if the ref matches
	if jr.Ref != h.Ref {
		status = http.StatusTeapot
		return
	}

	// check if the secret matches
	if req.Header.Get("X-Gitlab-Token") != h.Secret {
		status = http.StatusForbidden
		return
	}

	status = http.StatusOK
	return
}
