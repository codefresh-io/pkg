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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnvironment_bootstrapUrl(t *testing.T) {
	tests := map[string]struct {
		env  *environment
		want string
	}{
		"Simple": {
			&environment{
				TemplateRef: "https://github.com/foo/bar",
			},
			"https://github.com/foo/bar/" + bootstrapDir,
		},
		"With Tag": {
			&environment{
				TemplateRef: "https://github.com/foo/bar@v0.0.1",
			},
			"https://github.com/foo/bar/" + bootstrapDir + "?ref=v0.0.1",
		},
		"With Branch Name": {
			&environment{
				TemplateRef: "https://github.com/foo/bar#fizz",
			},
			"https://github.com/foo/bar/" + bootstrapDir + "?ref=fizz",
		},
		"With Branch SHA": {
			&environment{
				TemplateRef: "https://github.com/foo/bar#f24fcad",
			},
			"https://github.com/foo/bar/" + bootstrapDir + "?ref=f24fcad",
		},
	}
	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			got := tt.env.bootstrapUrl()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getAppFromFile(t *testing.T) {
	basicApp := &application{
		&v1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: v1alpha1.ApplicationSpec{
				Source: v1alpha1.ApplicationSource{
					RepoURL: "https://github.com/foo/bar",
					Path:    "kustomize/entities/overlays",
				},
				Destination: v1alpha1.ApplicationDestination{
					Server:    "https://kubernetes.default.svc",
					Namespace: "foo",
				},
			},
		},
		"",
		nil,
	}

	tests := map[string]struct {
		data []byte
		want *application
		err  string
	}{
		"Simple": {
			data: []byte(`
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: "foo"
spec:
  source:
    repoURL: https://github.com/foo/bar
    targetRevision: HEAD
    path: kustomize/entities/overlays
  destination:
    server: https://kubernetes.default.svc
    namespace: "foo"
`),
			want: basicApp,
			err:  "",
		},

		"Should return only first app": {
			data: []byte(`
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: "foo"
spec:
  source:
    repoURL: https://github.com/foo/bar
    targetRevision: HEAD
    path: kustomize/entities/overlays
  destination:
    server: https://kubernetes.default.svc
    namespace: "foo"

---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: "bar"
spec:
  source:
    repoURL: https://github.com/bar/bar
    targetRevision: HEAD
    path: kustomize/entities/overlays
  destination:
    server: https://kubernetes.default.svc
    namespace: "bar"
`),
			want: basicApp,
			err:  "",
		},
		"Should ignore other kinds": {
			data: []byte(`
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: "foo-proj"
spec:
  description: "foo project"
  sourceRepos:
  - "*"
  destinations:
  - namespace: "*"
    server: https://kubernetes.default.svc

---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: "foo"
spec:
  source:
    repoURL: https://github.com/foo/bar
    targetRevision: HEAD
    path: kustomize/entities/overlays
  destination:
    server: https://kubernetes.default.svc
    namespace: "foo"
`),
			want: basicApp,
			err:  "",
		},
		"Should fail to unmarshal": {
			data: []byte("foo"),
			want: nil,
			err:  "failed to unmarshal object in",
		},
	}
	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			tf, err := ioutil.TempFile("", "")
			assert.NoError(t, err)
			defer func() { _ = os.Remove(tf.Name()) }()

			_, err = tf.Write(tt.data)
			assert.NoError(t, err)
			env := &environment{
				c: &config{
					path: filepath.Dir(tf.Name()),
				},
			}
			got, err := env.getAppFromFile(tf.Name())
			if tt.err != "" {
				assert.Contains(t, err.Error(), tt.err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tf.Name(), got.Path)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Spec.Source.Path, got.SrcPath())
			assert.Equal(t, tt.want.Spec.Source.RepoURL, got.Spec.Source.RepoURL)
		})
	}
}

func TestApplication_childApps(t *testing.T) {
	must := func(path string, err error) string {
		assert.NoError(t, err)
		return path
	}

	tests := map[string]struct {
		env  *environment
		want []*application
		err  string
	}{
		"Simple": {
			&environment{
				c: &config{
					path: must(filepath.Abs("../../test/e2e/structures/uc1")),
				},
				RootApplicationPath: "root.yaml",
			},
			[]*application{
				{
					&v1alpha1.Application{
						ObjectMeta: metav1.ObjectMeta{
							Name: "leaf",
						},
					},
					must(filepath.Abs("../../test/e2e/structures/uc1/apps/app1.yaml")),
					nil,
				},
			},
			"",
		},
		"Two levels": {
			&environment{
				c: &config{
					path: must(filepath.Abs("../../test/e2e/structures/uc2")),
				},
				RootApplicationPath: "root.yaml",
			},
			[]*application{
				{
					&v1alpha1.Application{
						ObjectMeta: metav1.ObjectMeta{
							Name: "child1",
						},
					},
					must(filepath.Abs("../../test/e2e/structures/uc2/apps/app1.yaml")),
					nil,
				},
				{
					&v1alpha1.Application{
						ObjectMeta: metav1.ObjectMeta{
							Name: "leaf2",
						},
					},
					must(filepath.Abs("../../test/e2e/structures/uc2/apps/app2.yaml")),
					nil,
				},
			},
			"",
		},
	}

	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			app, err := tt.env.getRootApp()
			if tt.err != "" {
				assert.Contains(t, err.Error(), tt.err)
				return
			}
			assert.NoError(t, err)

			got, err := app.childApps()
			if tt.err != "" {
				assert.Contains(t, err.Error(), tt.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			for i, ca := range got {
				assert.Equal(t, tt.want[i].Path, ca.Path)
				assert.Equal(t, tt.want[i].Name, ca.Name)
			}
		})
	}
}

func TestApplication_leafApps(t *testing.T) {
	must := func(path string, err error) string {
		assert.NoError(t, err)
		return path
	}

	tests := map[string]struct {
		env  *environment
		want []*application
		err  string
	}{
		"Simple": {
			&environment{
				c: &config{
					path: must(filepath.Abs("../../test/e2e/structures/uc1")),
				},
				RootApplicationPath: "root.yaml",
			},
			[]*application{
				{
					&v1alpha1.Application{
						ObjectMeta: metav1.ObjectMeta{
							Name: "leaf",
						},
					},
					must(filepath.Abs("../../test/e2e/structures/uc1/apps/app1.yaml")),
					nil,
				},
			},
			"",
		},
		"Two levels": {
			&environment{
				c: &config{
					path: must(filepath.Abs("../../test/e2e/structures/uc2")),
				},
				RootApplicationPath: "root.yaml",
			},
			[]*application{
				{
					&v1alpha1.Application{
						ObjectMeta: metav1.ObjectMeta{
							Name: "leaf1",
						},
					},
					must(filepath.Abs("../../test/e2e/structures/uc2/apps/third/app3.yaml")),
					nil,
				},
				{
					&v1alpha1.Application{
						ObjectMeta: metav1.ObjectMeta{
							Name: "leaf2",
						},
					},
					must(filepath.Abs("../../test/e2e/structures/uc2/apps/app2.yaml")),
					nil,
				},
			},
			"",
		},
	}

	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			app, err := tt.env.getRootApp()
			if tt.err != "" {
				assert.Contains(t, err.Error(), tt.err)
				return
			}
			assert.NoError(t, err)

			got, err := app.LeafApps()
			if tt.err != "" {
				assert.Contains(t, err.Error(), tt.err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			// for i, ca := range got {
			// 	assert.Equal(t, tt.want[i].Path, ca.Path)
			// 	assert.Equal(t, tt.want[i].Name, ca.Name)
			// }
		})
	}
}
