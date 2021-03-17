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
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/go-git/go-billy/v5"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	ErrAppNotFound              = errors.New("app not found")
	yamlSeparator = regexp.MustCompile(`\n---`)
)

type (
	Environment interface {
		Name() string
		AddManifest(appName string, manifest []byte) (string, error)
	}

	environment struct {
		fs billy.Filesystem
		rootPath string
	}
)

func NewEnvironment(fs billy.Filesystem, rootPath string) Environment {
	return &environment{fs, rootPath}
}

func (e *environment) Name() string {
	return filepath.Dir(e.rootPath)
}

func (e *environment) AddManifest(appName string, manifest []byte) (string, error) {
	app, err := e.getApp(appName)
	if err != nil {
		return "", err
	}

	return app.AddManifest(manifest)
}

func (e *environment) getApp(appName string) (*application, error) {
	yamls, err := getYamls(e.fs, e.rootPath)
	if err != nil {
		return nil, err
	}

	for _, yaml := range(yamls) {
		app, _ := e.getAppFromFile(yaml)

		res, _ := e.getAppRecurse(app, appName)
		if res != nil {
			return res, nil
		}
	}

	return nil, ErrAppNotFound
}

func (e *environment) getAppRecurse(root *application, appName string) (*application, error) {
	if root == nil || root.IsManaged() {
		return nil, nil
	}

	if root.LabelName() == appName {
		return root, nil
	}

	appsDir := root.SrcPath() // check if it's not in this repo

	filenames, err := getYamls(e.fs, appsDir)
	if err != nil {
		return nil, err
	}

	for _, f := range filenames {
		app, err := e.getAppFromFile(f)
		if err != nil || app == nil {
			// not an argocd app - ignore
			continue
		}

		res, err := e.getAppRecurse(app, appName)
		if err != nil || res != nil {
			return res, err
		}
	}

	return nil, nil
}

func (e *environment) getAppFromFile(path string) (*application, error) {
	data, err := readFile(e.fs, path)
	if err != nil {
		return nil, err
	}

	for _, text := range yamlSeparator.Split(string(data), -1) {
		if strings.TrimSpace(text) == "" {
			continue
		}

		u := &unstructured.Unstructured{}

		err := yaml.Unmarshal([]byte(text), u)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal object in %s: %w", path, err)
		}

		if u.GetKind() == "Application" {
			app := &v1alpha1.Application{}
			if err := yaml.Unmarshal([]byte(text), app); err != nil {
				return nil, err
			}

			return &application{app, path, e}, nil
		}
	}

	return nil, nil
}

// func createDummy(fs billy.Basic, path string) error {
// 	file, err := fs.Create(fs.Join(path, "DUMMY"))
// 	if err != nil {
// 		return err
// 	}

// 	return file.Close()
// }
