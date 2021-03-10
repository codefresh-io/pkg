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

import "testing"

func Test_manager_AddManifest(t *testing.T) {
	type args struct {
		repo     string
		envName  string
		appName  string
		manifest []byte
	}

	tests := map[string]struct {
		m       *manager
		args    args
		wantErr bool
	}{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := &manager{}
			if err := m.AddManifest(tt.args.repo, tt.args.envName, tt.args.appName, tt.args.manifest); (err != nil) != tt.wantErr {
				t.Errorf("manager.AddManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
