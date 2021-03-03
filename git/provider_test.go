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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewProvider(t *testing.T) {
	tests := map[string]struct {
		opts             *Options
		expectedProvider Provider
		expectedError    string
	}{
		"Github": {
			&Options{
				Type: "github",
			},
			&github{},
			"",
		},
		"No Type": {
			&Options{},
			nil,
			ErrProviderNotSupported.Error(),
		},
		"Bad Type": {
			&Options{Type: "foo"},
			nil,
			ErrProviderNotSupported.Error(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := NewProvider(test.opts)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
				return
			}
			assert.IsType(t, test.expectedProvider, p)
		})
	}
}
