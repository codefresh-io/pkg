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
//go:generate mockery -name Provider
//go:generate mockery -name Repository

package git

import (
	"context"
	"errors"
)

type (
	// Provider represents a git provider
	Provider interface {
		// CreateRepository creates the repository in the remote provider and returns a
		// clone url
		CreateRepository(ctx context.Context, opts *CreateRepoOptions) (string, error)

		GetRepository(ctx context.Context, opts *GetRepoOptions) (string, error)

		// CloneRepository tries to clone the repository and return it if it exists or
		// ErrRepoNotFound if the repo does not exist
		CloneRepository(ctx context.Context, cloneURL string) (Repository, error)
	}

	// Options for a new git provider
	Options struct {
		Type string
		Auth *Auth
		Host string
	}

	// Auth for git provider
	Auth struct {
		Username string
		Password string
	}

	CreateRepoOptions struct {
		Owner   string
		Name    string
		Private bool
	}

	GetRepoOptions struct {
		Owner string
		Name  string
	}
)

// Errors
var (
	ErrProviderNotSupported = errors.New("git provider not supported")
)

// New creates a new git provider
func NewProvider(opts *Options) (Provider, error) {
	switch opts.Type {
	case "github":
		return newGithub(opts)
	default:
		return nil, ErrProviderNotSupported
	}
}
