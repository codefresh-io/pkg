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
	"fmt"

	"net/http"

	gh "github.com/google/go-github/v32/github"
)

type github struct {
	opts   *Options
	client *gh.Client
}

func newGithub(opts *Options) (Provider, error) {
	var (
		c   *gh.Client
		err error
	)

	hc := &http.Client{}

	if opts.Auth != nil {
		hc.Transport = &gh.BasicAuthTransport{
			Username: opts.Auth.Username,
			Password: opts.Auth.Password,
		}
	}

	if opts.Host != "" {
		c, err = gh.NewEnterpriseClient(opts.Host, opts.Host, hc)
		if err != nil {
			return nil, err
		}
	} else {
		c = gh.NewClient(hc)
	}

	g := &github{
		opts:   opts,
		client: c,
	}

	return g, nil
}

func (g *github) GetRepository(ctx context.Context, opts *GetRepoOptions) (string, error) {
	r, res, err := g.client.Repositories.Get(ctx, opts.Owner, opts.Name)

	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return "", ErrRepoNotFound
		}

		return "", err
	}

	return *r.CloneURL, nil
}

func (g *github) CreateRepository(ctx context.Context, opts *CreateRepoOptions) (string, error) {
	authUser, _, err := g.client.Users.Get(ctx, "") // get authenticated user details
	if err != nil {
		return "", err
	}

	org := ""
	if *authUser.Login != opts.Owner {
		org = opts.Owner
	}

	r, _, err := g.client.Repositories.Create(ctx, org, &gh.Repository{
		Name:    gh.String(opts.Name),
		Private: gh.Bool(opts.Private),
	})
	if err != nil {
		return "", err
	}

	if r.CloneURL == nil {
		return "", fmt.Errorf("repo clone url is nil")
	}

	return *r.CloneURL, err
}
