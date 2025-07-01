package rateg

import (
	"context"

	"golang.org/x/time/rate"
)

type LimitFnOption struct {
	Limit  float64
	Bucket int
}

type Fn[T any, R any] func(T) (R, error)
type FnCtx[T any, R any] func(context.Context, T) (R, error)

// LimitFn limits the function to be called only once at a specified time interval
func LimitFn[T any, R any](f Fn[T, R], option LimitFnOption) Fn[T, R] {
	// Use closures to save a limiter
	limiter := rate.NewLimiter(rate.Limit(option.Limit), option.Bucket)
	// Returns a new function, which is used to limit the function to be called only once at a specified time interval
	return func(t T) (R, error) {
		err := limiter.Wait(context.Background())
		if err != nil {
			var empty R
			return empty, err
		}
		// Execute the function that needs to be limited
		return f(t)
	}
}

// LimitFnCtx limits the function to be called only once at a specified time interval
func LimitFnCtx[T any, R any](f FnCtx[T, R], option LimitFnOption) FnCtx[T, R] {
	// Use closures to save a limiter
	limiter := rate.NewLimiter(rate.Limit(option.Limit), option.Bucket)
	// Returns a new function, which is used to limit the function to be called only once at a specified time interval
	return func(ctx context.Context, t T) (R, error) {
		err := limiter.Wait(ctx)
		if err != nil {
			var empty R
			return empty, err
		}
		// Execute the function that needs to be limited
		return f(ctx, t)
	}
}
