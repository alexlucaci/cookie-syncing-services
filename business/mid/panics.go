package mid

import (
	"context"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"runtime/debug"
)

// Panics returns a middleware that has the ability to recover from panics and converts
// the panic to an error so that is handled in Errors middleware.
func Panics(log *log.Logger) web.Middleware {

	// Middleware function the be executed.
	m := func(after web.Handler) web.Handler {

		// Handler to be attached in the chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// If the context object is missing the request values
			// do a graceful shutdown.
			rvs, ok := ctx.Value(web.KeyRequestValues).(*web.RequestValues)
			if !ok {
				return web.NewShutdownError("request values missing from context")
			}

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s :\n%s", rvs.TraceID, debug.Stack())
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return after(ctx, w, r)
		}

		return h
	}

	return m
}
