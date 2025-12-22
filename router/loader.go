package router

import (
	"context"
	"time"
)

// LoaderFunc loads data before rendering a page
type LoaderFunc func(ctx context.Context, params Params) (any, error)

// LoaderTimeout is the default timeout for data loaders
const LoaderTimeout = 30 * time.Second

// Loader sets the data loader for a route
func (r *Route) Loader(fn LoaderFunc) *Route {
	r.LoaderFn = fn
	return r
}

// executeLoader runs the loader function with timeout support
// This function assumes LoaderFn is not nil and should only be called when LoaderFn is set
func (r *Route) executeLoader(ctx context.Context, params Params) (any, error) {
	// Create timeout context
	loaderCtx, cancel := context.WithTimeout(ctx, LoaderTimeout)
	defer cancel()

	// Channel for result
	type result struct {
		data any
		err  error
	}

	done := make(chan result, 1)

	// Execute loader in goroutine
	go func() {
		data, err := r.LoaderFn(loaderCtx, params)
		done <- result{data, err}
	}()

	// Wait for result or timeout
	select {
	case res := <-done:
		return res.data, res.err
	case <-loaderCtx.Done():
		return nil, &LoaderError{
			Status:  408,
			Message: "Data loader timed out",
			Err:     loaderCtx.Err(),
		}
	}
}
