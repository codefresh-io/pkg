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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codefresh-io/cf-argo/pkg/kube"
	"github.com/codefresh-io/cf-argo/pkg/store"
	"github.com/codefresh-io/pkg/helpers"
	"github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// errors
var (
	ErrEnvironmentAlreadyExists = errors.New("environment already exists")
	ErrEnvironmentNotExist      = errors.New("environment does not exist")
	ErrAppNotFound              = errors.New("app not found")

	ConfigFileName = fmt.Sprintf("%s.yaml", store.AppName)
)

const (
	configVersion = "1.0"
)

type (
	Config interface {
		Persist() error
		GetEnvironment(name string) (Environment, error)
		AddEnvironmentP(ctx context.Context, env Environment, values interface{}, dryRun bool) error
		DeleteEnvironmentP(ctx context.Context, name string, values interface{}, dryRun bool) error
	}

	config struct {
		path         string                  // the path from which the config was loaded
		Version      string                  `json:"version"`
		Environments map[string]*environment `json:"environments"`
	}
)

func NewConfig(path string) Config {
	return &config{
		path:         path,
		Version:      configVersion,
		Environments: make(map[string]*environment),
	}
}

// LoadConfig loads the config from the specified path
func LoadConfig(path string) (Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, ConfigFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", path)
		}

		return nil, err
	}

	c := new(config)
	c.path = path

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	for name, e := range c.Environments {
		e.c = c
		e.name = name
	}

	return c, nil
}

// Persist saves the config to file
func (c *config) Persist() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(c.path, ConfigFileName), data, 0644)
}

func (c *config) GetEnvironment(name string) (e Environment, err error) {
	e, exists := c.Environments[name]
	if !exists {
		err = fmt.Errorf("%w: %s", ErrEnvironmentNotExist, name)
	}

	return
}

// AddEnvironmentP adds a new environment, copies all of the argocd apps to the relative
// location in the repository that c is managing, and persists the config object
func (c *config) AddEnvironmentP(ctx context.Context, env Environment, values interface{}, dryRun bool) error {
	if _, exists := c.Environments[env.Name()]; exists {
		return fmt.Errorf("%w: %s", ErrEnvironmentAlreadyExists, env.Name())
	}

	// copy all of the argocd apps to the correct location in the destination repo
	newEnv, err := c.installEnv(env)
	if err != nil {
		return err
	}

	c.Environments[env.Name()] = newEnv
	if err = c.Persist(); err != nil {
		return err
	}

	cs, err := store.Get().NewKubeClient(ctx).KubernetesClientSet()
	if err != nil {
		return err
	}

	_, err = cs.CoreV1().Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s-argocd", env.Name())},
	}, metav1.CreateOptions{})
	if err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return err
		}
	}

	manifests, err := kube.KustBuild(newEnv.bootstrapUrl(), values)
	if err != nil {
		return err
	}

	return store.Get().NewKubeClient(ctx).Apply(ctx, &kube.ApplyOptions{
		Manifests: manifests,
		DryRun:    dryRun,
	})
}

// DeleteEnvironmentP deletes an environment and persists the config object
func (c *config) DeleteEnvironmentP(ctx context.Context, name string, values interface{}, dryRun bool) error {
	env, exists := c.Environments[name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotExist, name)
	}

	err := env.cleanup()
	if err != nil {
		return err
	}

	delete(c.Environments, name)

	err = c.Persist()
	if err != nil {
		return err
	}

	manifests, err := kube.KustBuild(env.bootstrapUrl(), values)
	if err != nil {
		return err
	}

	return store.Get().NewKubeClient(ctx).Delete(ctx, &kube.DeleteOptions{
		Manifests: manifests,
		DryRun:    dryRun,
	})
}

func (c *config) firstEnv() *environment {
	for _, env := range c.Environments {
		return env
	}

	return nil
}

func (c *config) installEnv(env Environment) (*environment, error) {
	castEnv := env.(*environment)

	lapps, err := castEnv.leafApps()
	if err != nil {
		return nil, err
	}

	newEnv := &environment{
		name:                castEnv.name,
		c:                   c,
		TemplateRef:         castEnv.TemplateRef,
		RootApplicationPath: castEnv.RootApplicationPath,
	}

	for _, la := range lapps {
		if la.IsManaged() {
			if err = newEnv.installApp(castEnv.c.path, la.(*application)); err != nil {
				return nil, err
			}
		}
	}

	// copy the tpl "argocd-apps" to the matching dir in the dst repo
	src := filepath.Join(castEnv.c.path, filepath.Dir(castEnv.RootApplicationPath))

	var dstApplicationPath string

	if len(c.Environments) == 0 {
		dstApplicationPath = castEnv.RootApplicationPath
	} else {
		dstApplicationPath = c.firstEnv().RootApplicationPath
	}

	dst := filepath.Join(c.path, filepath.Dir(dstApplicationPath))

	err = helpers.CopyDir(src, dst)
	if err != nil {
		return nil, err
	}

	return newEnv, nil
}

func (c *config) getApp(appName string) (*application, error) {
	err := ErrAppNotFound

	var app *application

	for _, e := range c.Environments {
		app, err = e.getApp(appName)
		if err != nil && !errors.Is(err, ErrAppNotFound) {
			return nil, err
		}

		if app != nil {
			return app, nil
		}
	}

	return app, err
}
