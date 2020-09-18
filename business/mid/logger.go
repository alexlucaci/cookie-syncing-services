package mid

import (
	"context"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"log"
	"net/http"
	"time"
)

// Logger returns a middleware that writes information about the request to the logs.
// Format used: TraceID: (404) POST /probe -> IP ADDR (latency).
func Logger(log *log.Logger) web.Middleware {

	// Middleware function the be executed.
	m := func(before web.Handler) web.Handler {

		// Handler to be attached in the chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context object is missing the request values
			// do a graceful shutdown.
			rvs, ok := ctx.Value(web.KeyRequestValues).(*web.RequestValues)
			if !ok {
				return web.NewShutdownError("request values missing from context")
			}

			err := before(ctx, w, r)

			log.Printf("%s : (%d) : %s %s -> %s (%s)",
				rvs.TraceID, rvs.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(rvs.Now),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
