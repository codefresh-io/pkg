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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/codefresh-io/pkg/helpers"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// errors
var (
	yamlSeparator = regexp.MustCompile(`\n---`)
)

const (
	bootstrapDir = "bootstrap"
)

type (
	Environment interface {
		Name() string
		Uninstall() (bool, error)
		AddManifest(appName string, manifest []byte) error
	}

	environment struct {
		c                   *config
		name                string
		RootApplicationPath string `json:"rootAppPath"`
		TemplateRef         string `json:"templateRef"`
	}
)

func (e *environment) Name() string {
	return e.name
}

// Uninstall removes all managed apps and returns true if there are no more
// apps left in the environment.
func (e *environment) Uninstall() (bool, error) {
	rootApp, err := e.getRootApp()
	if err != nil {
		return false, err
	}

	uninstalled, err := rootApp.Uninstall()
	if uninstalled {
		return true, createDummy(filepath.Join(e.c.path, rootApp.SrcPath()))
	}

	return false, err
}

func (e *environment) AddManifest(appName string, manifest []byte) error {
	app, err := e.getApp(appName)
	if err != nil {
		return err
	}

	return app.AddManifest(manifest)
}

func (e *environment) getRootApp() (*application, error) {
	return e.getAppFromFile(filepath.Join(e.c.path, e.RootApplicationPath))
}

func (e *environment) getApp(appName string) (*application, error) {
	rootApp, err := e.getRootApp()
	if err != nil {
		return nil, err
	}

	app, err := e.getAppRecurse(rootApp, appName)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, fmt.Errorf("%w: %s", ErrAppNotFound, appName)
	}

	return app, nil
}

func (e *environment) leafApps() ([]Application, error) {
	rootApp, err := e.getRootApp()
	if err != nil {
		return nil, err
	}

	return rootApp.LeafApps()
}

func (e *environment) bootstrapUrl() string {
	var parts []string

	switch {
	case strings.Contains(e.TemplateRef, "#"):
		parts = strings.Split(e.TemplateRef, "#")
	case strings.Contains(e.TemplateRef, "@"):
		parts = strings.Split(e.TemplateRef, "@")
	default:
		parts = []string{e.TemplateRef}
	}

	bootstrapUrl := fmt.Sprintf("%s/%s", parts[0], bootstrapDir)

	if len(parts) > 1 {
		return fmt.Sprintf("%s?ref=%s", bootstrapUrl, parts[1])
	}

	return bootstrapUrl
}

func (e *environment) cleanup() error {
	_, err := e.getRootApp()
	if err != nil {
		return err
	}

	return nil // rootApp.deleteFromFilesystem()
}

func (e *environment) installApp(srcRootPath string, app *application) error {
	appName := app.LabelName()

	refApp, err := e.c.getApp(appName)
	if err != nil {
		if !errors.Is(err, ErrAppNotFound) {
			return err
		}

		return e.installNewApp(srcRootPath, app)
	}

	baseLocation, err := refApp.getBaseLocation()
	if err != nil {
		return err
	}

	absSrc := filepath.Join(srcRootPath, app.SrcPath())

	dst := filepath.Clean(filepath.Join(baseLocation, "..", "overlays", e.name))
	absDst := filepath.Join(e.c.path, dst)

	err = helpers.CopyDir(absSrc, absDst)
	if err != nil {
		return err
	}

	app.setSrcPath(dst)

	return app.save()
}

func (e *environment) installNewApp(srcRootPath string, app Application) error {
	appFolder := filepath.Clean(filepath.Join(app.SrcPath(), "..", ".."))
	absSrc := filepath.Join(srcRootPath, appFolder)
	absDst := filepath.Join(e.c.path, appFolder)

	return helpers.CopyDir(absSrc, absDst)
}

func (e *environment) getAppRecurse(root *application, appName string) (*application, error) {
	if root.LabelName() == appName {
		return root, nil
	}

	appsDir := root.SrcPath() // check if it's not in this repo

	filenames, err := filepath.Glob(filepath.Join(e.c.path, appsDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	for _, f := range filenames {
		app, err := e.getAppFromFile(f)
		if err != nil || app == nil {
			// not an argocd app - ignore
			continue
		}

		if !app.IsManaged() {
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
	data, err := ioutil.ReadFile(path)
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

func createDummy(path string) error {
	file, err := os.Create(filepath.Join(path, "DUMMY"))
	if err != nil {
		return err
	}

	return file.Close()
}
