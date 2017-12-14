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

func parseGitLabRequest(req *http.Request) (project, secret string) {
	var jr gitlabRequest
	jd := json.NewDecoder(req.Body)
	err := jd.Decode(&jr)
	if err != nil {
		return
	}

	if jr.ObjectKind != "push" {
		return
	}

	project = jr.Project.Name
	secret = req.Header.Get("X-Gitlab-Token")
	return
}
