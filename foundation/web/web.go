// web package contains a mini web framework.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// ctxKey represents the value type for the keys stored in the context.
type ctxKey int

// KeyRequestValues is the key used to store and retrieve a RequestValues from context.Context.
const KeyRequestValues ctxKey = 1

// RequestValues represents the state of each request.
type RequestValues struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// Handler is a type which handles an http request in the framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App represents the entrypoint of our application and configures
// the context objects for each http handlers.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mws      []Middleware
}

func NewApp(shutdown chan os.Signal, mws ...Middleware) *App {
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mws:        mws,
	}

	return &app
}

// Handle is where we bind a given Handler to a given HTTP verb and path.
// Here we are also wrapping the middlewares.
func (a *App) Handle(method string, path string, handler Handler, mws ...Middleware) {

	// First wrap the handler's specific middlewares around the handler.
	handler = wrapMiddleware(mws, handler)

	// Then wrap the application's middlewares to the handler.
	handler = wrapMiddleware(a.mws, handler)

	// And finally the function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Set the context object with the request values.
		rvs := RequestValues{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyRequestValues, &rvs)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	// Mount the handler for the specified verb and path.
	a.ContextMux.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
