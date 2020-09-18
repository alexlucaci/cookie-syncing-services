package web

// Middleware is a function designed to run code before and/or after another Handler.
type Middleware func(Handler) Handler

// wrapMiddleware wraps middleware(s) around a Handler.
func wrapMiddleware(mws []Middleware, handler Handler) Handler {

	// Looping backwards means that the first middleware
	// of the slice will be executed first when a request comes.
	for i := len(mws) - 1; i >= 0; i-- {
		h := mws[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
