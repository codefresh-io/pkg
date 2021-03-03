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
//go:generate mockery -name Repository

package git

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	gg "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type (
	// Repository represents a git repository
	Repository interface {
		Add(ctx context.Context, pattern string) error

		AddRemote(ctx context.Context, name, url string) error

		// Commits all files and returns the commit sha
		Commit(ctx context.Context, msg string) (string, error)

		Push(context.Context, *PushOptions) error

		IsNewRepo() (bool, error)

		Root() (string, error)
	}

	PushOptions struct {
		RemoteName string
		Auth       *Auth
	}

	repo struct {
		r *gg.Repository
	}
)

// Errors
var (
	ErrNilOpts      = errors.New("options cannot be nil")
	ErrRepoNotFound = errors.New("git repository not found")
)

// go-git functions (we mock those in tests)
var (
	plainClone = gg.PlainCloneContext
	plainInit  = gg.PlainInit
)

func Clone(ctx context.Context, opts *CloneOptions) (Repository, error) {
	if opts == nil {
		return nil, ErrNilOpts
	}

	auth := getAuth(opts.Auth)

	cloneOpts := &gg.CloneOptions{
		Depth:    1,
		URL:      opts.URL,
		Auth:     auth,
		Progress: os.Stderr,
	}

	if ref := getRef(opts.URL); ref != "" {
		cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(ref)
		cloneOpts.URL = opts.URL[:strings.LastIndex(opts.URL, ref)-1]
	} else if i := strings.LastIndex(opts.URL, "@"); i > -1 {
		cloneOpts.ReferenceName = plumbing.NewTagReferenceName(opts.URL[i+1:])
		cloneOpts.URL = opts.URL[:i]
	}

	err := cloneOpts.Validate()
	if err != nil {
		return nil, err
	}

	r, err := plainClone(ctx, opts.Path, false, cloneOpts)
	if err != nil {
		return nil, err
	}

	return &repo{r}, nil
}

func Init(ctx context.Context, path string) (Repository, error) {
	if path == "" {
		path = "."
	}

	r, err := plainInit(path, false)
	if err != nil {
		return nil, err
	}

	return &repo{r}, err
}

func (r *repo) Add(ctx context.Context, pattern string) error {
	w, err := r.r.Worktree()
	if err != nil {
		return err
	}

	return w.AddGlob(pattern)
}

func (r *repo) AddRemote(ctx context.Context, name, url string) error {
	cfg := &config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	}

	err := cfg.Validate()
	if err != nil {
		return err
	}

	_, err = r.r.CreateRemote(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Commit(ctx context.Context, msg string) (string, error) {
	wt, err := r.r.Worktree()
	if err != nil {
		return "", err
	}

	h, err := wt.Commit(msg, &gg.CommitOptions{
		All: true,
	})
	if err != nil {
		return "", err
	}

	return h.String(), err
}

func (r *repo) Push(ctx context.Context, opts *PushOptions) error {
	if opts == nil {
		return ErrNilOpts
	}

	auth := getAuth(opts.Auth)
	pushOpts := &gg.PushOptions{
		RemoteName: opts.RemoteName,
		Auth:       auth,
		Progress:   os.Stdout,
	}

	err := pushOpts.Validate()
	if err != nil {
		return err
	}

	err = r.r.PushContext(ctx, pushOpts)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) IsNewRepo() (bool, error) {
	remotes, err := r.r.Remotes()
	if err != nil {
		return false, err
	}

	return len(remotes) == 0, nil
}

func (r *repo) Root() (string, error) {
	wt, err := r.r.Worktree()
	if err != nil {
		return "", err
	}

	return wt.Filesystem.Root(), nil
}

func getAuth(auth *Auth) transport.AuthMethod {
	if auth != nil {
		username := auth.Username
		if username == "" {
			username = "codefresh"
		}

		return &http.BasicAuth{
			Username: username,
			Password: auth.Password,
		}
	}

	return nil
}

func getRef(cloneURL string) string {
	u, err := url.Parse(cloneURL)
	if err != nil {
		return ""
	}

	return u.Fragment
}
