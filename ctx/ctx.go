package ctx

import (
	"context"
	"os"
	"os/signal"

	"github.com/rkrmr33/pkg/log"
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
				log.G().Warnw("got signal", "sig", s)
				cancel()
			} else {
				log.G().Warn("forcing exit")
				os.Exit(1)
			}
		}
	}()

	return ctx
}

// ContextWithCancelOnDone returns a context that depends on the done channel
// of another context.
func ContextWithCancelOnDone(done <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-done
		cancel()
	}()

	return ctx
}
