// Copyright 2021 The Codefresh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ops

import (
	"context"
	"io/ioutil"

	"github.com/codefresh-io/pkg/git"
)

type (
	Manager interface {
		AddManifest(repo, envName, appName string, manifest []byte) error
	}

	manager struct {
	}
)

var (
	clone = cloneRepository
)

func NewManager() Manager {
	return &manager{}
}

func cloneRepository(cloneURL string) (git.Repository, error) {
	clonePath, err := ioutil.TempDir("", "repo-")
	if err != nil {
		return nil, err
	}

	return git.Clone(context.TODO(), &git.CloneOptions{
		URL:  cloneURL,
		Path: clonePath,
		Auth: nil, //get token from filesystem?
	})
}

func (m *manager) AddManifest(repo, envName, appName string, manifest []byte) error {
	r, err := clone(repo)
	if err != nil {
		return err
	}

	rootPath, err := r.Root()
	if err != nil {
		return err
	}

	c, err := LoadConfig(rootPath)
	if err != nil {
		return err
	}

	err = c.AddManifest(envName, appName, manifest)
	if err != nil {
		return err
	}

	return nil
}
