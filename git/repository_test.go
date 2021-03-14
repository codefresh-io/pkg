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

	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/storage"
	gg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
				URL:  "https://github.com/foo/bar",
				Auth: nil,
			},
			expectedPath:    "/foo/bar",
			expectedURL:     "https://github.com/foo/bar",
			expectedRefName: plumbing.HEAD,
		},
		"With tag": {
			opts: &CloneOptions{
				URL:  "https://github.com/foo/bar@tag",
				Auth: nil,
			},
			expectedPath:    "/foo/bar",
			expectedURL:     "https://github.com/foo/bar",
			expectedRefName: plumbing.NewTagReferenceName("tag"),
		},
		"With branch": {
			opts: &CloneOptions{
				URL:  "https://github.com/foo/bar#branch",
				Auth: nil,
			},
			expectedPath:    "/foo/bar",
			expectedURL:     "https://github.com/foo/bar",
			expectedRefName: plumbing.NewBranchReferenceName("branch"),
		},
		"With token": {
			opts: &CloneOptions{
				URL: "https://github.com/foo/bar",
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

	orig := clone

	defer func() { clone = orig }()

	for name, test := range tests {
		clone2 = func(ctx context.Context, s storage.Storer, worktree billy.Filesystem, o *gg.CloneOptions) (*gg.Repository, error) {
			return nil, nil
		}
		clone = func(ctx context.Context, s storage.Storer, worktree billy.Filesystem, o *gg.CloneOptions) (*gg.Repository, error) {
			return nil, nil
		}
		// func(ctx context.Context, s storage.Storer, worktree billy.Filesystem, o *gg.CloneOptions) (*gg.Repository, error) {
		// 	assert.Equal(t, test.expectedURL, o.URL)
		// 	assert.Equal(t, test.expectedRefName, o.ReferenceName)
		// 	assert.Equal(t, 1, o.Depth)

		// 	if o.Auth != nil {
		// 		bauth, _ := o.Auth.(*http.BasicAuth)
		// 		assert.Equal(t, test.expectedPassword, bauth.Password)
		// 	}

		// 	return nil, nil
		// }

		t.Run(name, func(t *testing.T) {
			_, _ = Clone(context.Background(), test.opts)
		})
	}
}
