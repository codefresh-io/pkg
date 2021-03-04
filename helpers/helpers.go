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
