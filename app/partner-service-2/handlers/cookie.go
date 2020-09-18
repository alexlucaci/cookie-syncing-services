package handlers

import (
	"context"
	"fmt"
	"github.com/alexlucaci/cookie-syncing-services/business/data/cookie"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
)

// Use a struct to handle dependencies like db inside handlers
type cookieHandlers struct {
	db          *sqlx.DB
	service1Url string
}

func (h *cookieHandlers) sync(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rvs, ok := ctx.Value(web.KeyRequestValues).(*web.RequestValues)
	if !ok {
		return web.RespondWithError(ctx, w, r, errors.New("request values missing from context"))
	}

	var p2 string

	// If cookie does not exist, create it inside db and set it to the writer
	p2C, err := r.Cookie("p2")
	if err != nil {
		c, err := cookie.Create(ctx, h.db, rvs.Now)
		if err != nil {
			return web.NewRedirectError("/home", http.StatusMovedPermanently)
		}
		p2 = c.ID
		http.SetCookie(w, &http.Cookie{
			Name:  "p2",
			Value: p2,
			Path:  "/",
		})
	} else {
		p2 = p2C.Value
	}

	qParams := r.URL.Query()
	p1Q, ok := qParams["p1"]
	if !ok {
		return web.RespondWithError(ctx, w, r, errors.New("p1 was not found in query params"))
	}

	c, err := cookie.One(ctx, h.db, p2)
	if err != nil {
		if err == cookie.ErrNotFound {
			return web.RespondWithError(ctx, w, r, errors.New("cookie with given p2 id was not found"))
		}
		return web.RespondWithError(ctx, w, r, errors.New("cannot get cookie from db"))
	}

	// If cookie with p2 id does not have a partner cookie id then set it
	// the one got from query params
	if c.PartnerCookieID == nil {
		if err := cookie.SetPartnerCookie(ctx, h.db, c.ID, p1Q[0], rvs.Now); err != nil {
			return web.RespondWithError(ctx, w, r, errors.New("cannot set partner cookie"))
		}
	}

	url := fmt.Sprintf("%s/sync/partner2/servepixel?p1=%s&p2=%s", h.service1Url, p1Q[0], p2)

	return web.RespondWithRedirect(ctx, w, r, url)
}
