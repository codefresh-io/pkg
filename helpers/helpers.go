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
package helpers

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/codefresh-io/pkg/log"
)

// ContextWithCancelOnSignals returns a context that is canceled when one of the specified signals
// are received
func ContextWithCancelOnSignals(ctx context.Context, sigs ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, sigs...)

	go func() {
		s := <-sig
		log.G(ctx).Debugf("got signal: %s", s)
		cancel()
	}()

	return ctx
}

// Die panics if err is not nil
func Die(err error) {
	if err != nil {
		panic(err)
	}
}

func CopyDir(source, destination string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}

		absDst := filepath.Join(destination, relPath)
		if err = ensureDir(absDst); err != nil {
			return err
		}

		if info.IsDir() {
			err = os.Mkdir(absDst, info.Mode())
			if err != nil {
				if os.IsExist(err.(*os.PathError).Unwrap()) {
					return nil
				}
			}

			return err
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(absDst, data, info.Mode())
		}
	})
}

func ensureDir(path string) error {
	dstDir := filepath.Dir(path)
	if _, err := os.Stat(dstDir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		return os.MkdirAll(dstDir, 0755)
	}

	return nil
}
