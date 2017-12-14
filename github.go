// Copyright 2017 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
)

type githubRequest struct {
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

func parseGitHubRequest(req *http.Request) (project, secret string) {
	var jr githubRequest
	jd := json.NewDecoder(req.Body)
	err := jd.Decode(&jr)
	if err != nil {
		return
	}

	var push bool
	for _, event := range jr.Hook.Events {
		if event == "push" {
			push = true
			break
		}
	}
	if !push {
		return
	}

	project = jr.Repository.Name
	secret = jr.Hook.Config.Secret
	return
}
