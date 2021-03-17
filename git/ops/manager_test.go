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
	"testing"
)

func Test_Manager_AddManifest(t *testing.T) {
	type args struct {
		cloneURL string
		envName  string
		appName  string
		manifest []byte
	}

	tests := []struct {
		name    string
		m       Manager
		args    args
		wantErr bool
	}{
		{
			name: "asd",
			m:    NewManager(),
			args: args{
				cloneURL: "https://github.com/noam-codefresh/demo",
				envName:  "prod",
				appName:  "argo-workflows",
				manifest: []byte{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.AddManifest(tt.args.cloneURL, tt.args.envName, tt.args.appName, tt.args.manifest); (err != nil) != tt.wantErr {
				t.Errorf("manager.AddManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
