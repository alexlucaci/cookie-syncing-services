package mid

import (
	"context"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"log"
	"net/http"
)

// Errors returns a middleware that handles and log errors coming out of the call chain.
// Application defined errors are handled and specific response is returned based on the error type.
func Errors(log *log.Logger) web.Middleware {

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

			// Run the chain and catch any error.
			if err := before(ctx, w, r); err != nil {

				// Log the error.
				log.Printf("%s : ERROR : %v", rvs.TraceID, err)

				// Send an error response
				if err := web.RespondWithError(ctx, w, r, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}
			}

			// The error has been handled so we can stop propagating it.
			return nil
		}

		return h
	}

	return m
}
