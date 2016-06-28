package gax

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// A user defined call stub.
type APICall func(context.Context) error

// scaleDuration returns the product of a and mult.
func scaleDuration(a time.Duration, mult float64) time.Duration {
	ns := float64(a) * mult
	return time.Duration(ns)
}

// ensureTimeout returns a context with the given timeout applied if there
// is no deadline on the context.
func ensureTimeout(ctx context.Context, timeout time.Duration) context.Context {
	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}
	return ctx
}

// invokeWithRetry calls stub using an exponential backoff retry mechanism
// based on the values provided in retrySettings.
func invokeWithRetry(ctx context.Context, stub APICall, callSettings CallSettings) error {
	retrySettings := callSettings.RetrySettings
	backoffSettings := callSettings.RetrySettings.BackoffSettings
	delay := backoffSettings.DelayTimeoutSettings.Initial
	timeout := backoffSettings.RPCTimeoutSettings.Initial
	for {
		// If the deadline is exceeded...
		if ctx.Err() != nil {
			return ctx.Err()
		}
		timeoutCtx, _ := context.WithTimeout(ctx, backoffSettings.RPCTimeoutSettings.Max)
		timeoutCtx, _ = context.WithTimeout(timeoutCtx, timeout)
		err := stub(timeoutCtx)
		code := grpc.Code(err)
		if code == codes.OK {
			return nil
		}
		if !retrySettings.RetryCodes[code] {
			return err
		}
		delayCtx, _ := context.WithTimeout(ctx, backoffSettings.DelayTimeoutSettings.Max)
		delayCtx, _ = context.WithTimeout(delayCtx, delay)
		<-delayCtx.Done()

		delay = scaleDuration(delay, backoffSettings.DelayTimeoutSettings.Multiplier)
		timeout = scaleDuration(timeout, backoffSettings.RPCTimeoutSettings.Multiplier)
	}
}

// Invoke calls stub with a child of context modified by the specified options.
func Invoke(ctx context.Context, stub APICall, opts ...CallOption) error {
	settings := &CallSettings{}
	callOptions(opts).Resolve(settings)
	ctx = ensureTimeout(ctx, settings.Timeout)
	if len(settings.RetrySettings.RetryCodes) > 0 {
		return invokeWithRetry(ctx, stub, *settings)
	}
	return stub(ctx)
}
