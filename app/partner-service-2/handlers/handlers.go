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
func Service(shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, customerServiceUrl, service1Url string) http.Handler {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	ch := cookieHandlers{db: db, service1Url: service1Url}
	app.Handle(http.MethodGet, "/sync/partner1", ch.sync)

	return app
}
