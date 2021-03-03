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
package git

import (
	"context"
	"testing"

	gg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
)

func Test_Clone(t *testing.T) {
	tests := map[string]struct {
		opts             *CloneOptions
		expectedPath     string
		expectedURL      string
		expectedPassword string
		expectedRefName  plumbing.ReferenceName
	}{
		"Simple": {
			opts: &CloneOptions{
				Path: "/foo/bar",
				URL:  "https://github.com/foo/bar",
				Auth: nil,
			},
			expectedPath:    "/foo/bar",
			expectedURL:     "https://github.com/foo/bar",
			expectedRefName: plumbing.HEAD,
		},
		"With tag": {
			opts: &CloneOptions{
				Path: "/foo/bar",
				URL:  "https://github.com/foo/bar@tag",
				Auth: nil,
			},
			expectedPath:    "/foo/bar",
			expectedURL:     "https://github.com/foo/bar",
			expectedRefName: plumbing.NewTagReferenceName("tag"),
		},
		"With branch": {
			opts: &CloneOptions{
				Path: "/foo/bar",
				URL:  "https://github.com/foo/bar#branch",
				Auth: nil,
			},
			expectedPath:    "/foo/bar",
			expectedURL:     "https://github.com/foo/bar",
			expectedRefName: plumbing.NewBranchReferenceName("branch"),
		},
		"With token": {
			opts: &CloneOptions{
				Path: "/foo/bar",
				URL:  "https://github.com/foo/bar",
				Auth: &Auth{
					Password: "password",
				},
			},
			expectedPath:     "/foo/bar",
			expectedURL:      "https://github.com/foo/bar",
			expectedPassword: "password",
			expectedRefName:  plumbing.HEAD,
		},
	}

	orig := plainClone

	defer func() { plainClone = orig }()

	for name, test := range tests {
		plainClone = func(ctx context.Context, path string, isBare bool, o *gg.CloneOptions) (*gg.Repository, error) {
			assert.Equal(t, test.expectedPath, path)
			assert.Equal(t, test.expectedURL, o.URL)
			assert.Equal(t, test.expectedRefName, o.ReferenceName)
			assert.Equal(t, 1, o.Depth)
			assert.False(t, isBare)

			if o.Auth != nil {
				bauth, _ := o.Auth.(*http.BasicAuth)
				assert.Equal(t, test.expectedPassword, bauth.Password)
			}

			return nil, nil
		}

		t.Run(name, func(t *testing.T) {
			_, _ = Clone(context.Background(), test.opts)
		})
	}
}
