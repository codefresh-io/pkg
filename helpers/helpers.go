// Copyright 2022 The Codefresh Authors.
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
	"os"
	"os/signal"

	"github.com/codefresh-io/pkg/log"
)

// ContextWithCancelOnSignals returns a context that is canceled when one of the specified signals
// are received
func ContextWithCancelOnSignals(ctx context.Context, sigs ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, sigs...)

	go func() {
		cancels := 0

		for {
			s := <-sig
			cancels++

			if cancels == 1 {
				log.G(ctx).Printf("got signal: %s", s)
				cancel()
			} else {
				log.G(ctx).Printf("forcing exit")
				os.Exit(1)
			}
		}
	}()

	return ctx
}

// Die panics if err is not nil
func Die(err error) {
	if err != nil {
		panic(err)
	}
}
