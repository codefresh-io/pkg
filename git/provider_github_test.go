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

	g "github.com/codefresh-io/pkg/git/github"
	"github.com/stretchr/testify/assert"
)

func Test_github_GetRepository(t *testing.T) {
	type fields struct {
		Repositories g.Repositories
	}

	tests := map[string]struct {
		fields  fields
		opts    *GetRepoOptions
		want    string
		wantErr string
	}{
		// "simple": {
		// 	fields: fields{},
		// 	opts: &GetRepoOptions{
		// 		Owner: "owner",
		// 		Name:  "repo",
		// 	},
		// 	want:    "https://github.com/owner/repo",
		// 	wantErr: "",
		// },
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := &github{
				Repositories: tt.fields.Repositories,
			}
			got, err := g.GetRepository(context.Background(), tt.opts)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("github.GetRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
