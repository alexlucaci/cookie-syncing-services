package handlers

import (
	"github.com/alexlucaci/cookie-syncing-services/business/mid"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

// Service constructs an http.Handler with all application routes defined.
func Service(shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, customerServiceUrl, service2Url string) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	ch := cookieHandlers{
		db:                 db,
		customerServiceUrl: customerServiceUrl,
		service2Url:        service2Url,
	}
	app.Handle(http.MethodGet, "/whoami", ch.whoami)
	app.Handle(http.MethodGet, "/sync/partner2/servepixel", ch.servePixel)

	return app
}
