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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/codefresh-io/cf-argo/pkg/store"
	"github.com/ghodss/yaml"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

const (
	labelsManagedBy = "app.kubernetes.io/managed-by"
	labelsName      = "app.kubernetes.io/name"
)

type (
	Application interface {
		IsManaged() bool
		LabelName() string
		LeafApps() ([]Application, error)
		SrcPath() string
		Uninstall() (bool, error)
	}

	application struct {
		*v1alpha1.Application
		// Path the path from where the application manifest was read from
		Path string
		// env the environment that contains this application
		env *environment
	}
)

func (a *application) IsManaged() bool {
	return a.labelValue(labelsManagedBy) == store.AppName
}

func (a *application) LabelName() string {
	return a.labelValue(labelsName)
}

func (a *application) LeafApps() ([]Application, error) {
	childApps, err := a.childApps()
	if err != nil {
		return nil, err
	}

	isLeaf := true
	res := []Application{}

	for _, childApp := range childApps {
		isLeaf = false

		childRes, err := childApp.LeafApps()
		if err != nil {
			return nil, err
		}

		res = append(res, childRes...)
	}

	if isLeaf {
		res = append(res, a)
	}

	return res, nil
}

func (a *application) SrcPath() string {
	return a.Spec.Source.Path
}

func (a *application) Uninstall() (bool, error) {
	uninstalled := false

	childApps, err := a.childApps()
	if err != nil {
		return uninstalled, err
	}

	totalUninstalled := 0

	for _, childApp := range childApps {
		if childApp.IsManaged() {
			childUninstalled, err := childApp.Uninstall()
			if err != nil {
				return uninstalled, err
			}

			if childUninstalled {
				err = os.Remove(childApp.Path)
				if err != nil {
					return uninstalled, err
				}

				totalUninstalled++
			}
		}
	}

	return len(childApps) == totalUninstalled, nil
}

// func (a *application) deleteFromFilesystem() error {
// 	srcDir := filepath.Join(a.env.c.path, a.SrcPath())
// 	err := os.RemoveAll(srcDir)
// 	if err != nil {
// 		return err
// 	}

// 	projectPath := filepath.Join(filepath.Dir(a.Path), fmt.Sprintf("%s-project.yaml", a.Name))
// 	err = os.Remove(projectPath)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.Remove(a.Path)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (a *application) setSrcPath(newPath string) {
	a.Spec.Source.Path = newPath
}

func (a *application) labelValue(label string) string {
	if a.Labels == nil {
		return ""
	}

	return a.Labels[label]
}

func (a *application) getBaseLocation() (string, error) {
	refKust := filepath.Join(a.env.c.path, a.SrcPath(), "kustomization.yaml")

	bytes, err := ioutil.ReadFile(refKust)
	if err != nil {
		return "", err
	}

	k := &kustomize.Kustomization{}

	err = yaml.Unmarshal(bytes, k)
	if err != nil {
		return "", err
	}

	return filepath.Clean(filepath.Join(a.SrcPath(), k.Resources[0])), nil
}

func (a *application) save() error {
	data, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(a.Path, data, 0644)
}

func (a *application) childApps() ([]*application, error) {
	filenames, err := filepath.Glob(filepath.Join(a.env.c.path, a.SrcPath(), "*.yaml"))
	if err != nil {
		return nil, err
	}

	res := []*application{}

	for _, f := range filenames {
		childApp, err := a.env.getAppFromFile(f)
		if err != nil {
			fmt.Printf("file is not an argo-cd application manifest %s\n", f)
			continue
		}

		if childApp != nil {
			res = append(res, childApp)
		}
	}

	return res, nil
}
