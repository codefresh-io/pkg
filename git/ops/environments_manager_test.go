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
	"os"
	"testing"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
			fs := memfs.New()
			tf, err := fs.TempFile("", "")
			assert.NoError(t, err)
			defer func() { _ = os.Remove(tf.Name()) }()

			_, err = tf.Write(tt.data)
			assert.NoError(t, err)
			env := &environment{
				fs,
				"/",
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
