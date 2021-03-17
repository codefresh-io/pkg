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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/codefresh-io/pkg/git"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

type (
	Manager interface {
		AddManifest(cloneURL, rootSrcPath, appName string, manifest []byte) error
	}

	manager struct {
		fs billy.Filesystem
	}
)

func NewManager() Manager {
	return &manager{
		fs: memfs.New(),
	}
}

func (m *manager) AddManifest(cloneURL, rootSrcPath, appName string, manifest []byte) error {
	ctx := context.TODO()

	r, err := m.cloneRepository(ctx, cloneURL)
	if err != nil {
		return err
	}

	e := m.loadEnvironment(rootSrcPath)

	filename, err := e.AddManifest(appName, manifest)
	if err != nil {
		return err
	}

	return pushChanges(ctx, r, fmt.Sprintf("Added manifest for '%s'", filename))
}

func (m *manager) cloneRepository(ctx context.Context, cloneURL string) (git.Repository, error) {
	return git.Clone(ctx, m.fs, &git.CloneOptions{
		URL:  cloneURL,
		Auth: nil, //get token from filesystem?
	})
}

func (m *manager) loadEnvironment(rootSrcPath string) Environment {
	return NewEnvironment(m.fs, rootSrcPath)
}

func pushChanges(ctx context.Context, repo git.Repository, msg string) error {
	err := repo.Add(ctx, ".")
	if err != nil {
		return err
	}

	_, err = repo.Commit(ctx, msg)
	if err != nil {
		return err
	}

	return repo.Push(ctx, &git.PushOptions{
		Auth: &git.Auth{
			Password: "b3b09fbc4e3e6fa3a6f1082bf9957839a8f6c652",
		},
	})
}

func readFile(fs billy.Basic, filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}

func writeFile(fs billy.Basic, filename string, data []byte) error {
	f, err := fs.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(data)

	return err
}

func getYamls(fs billy.Filesystem, path string) ([]string, error) {
	fi, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)

	for _, f := range fi {
		name := f.Name()

		ext := filepath.Ext(name)
		if ext == ".yaml" || ext == ".yml" {
			res = append(res, fs.Join(path, name))
		}
	}

	return res, nil
}
