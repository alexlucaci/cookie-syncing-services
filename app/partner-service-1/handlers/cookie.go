package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alexlucaci/cookie-syncing-services/business/data/cookie"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

const base64GifPixel = "R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs="

// Use a struct to handle dependencies like db inside handlers
type cookieHandlers struct {
	db                 *sqlx.DB
	customerServiceUrl string
	service2Url        string
}

func (h *cookieHandlers) whoami(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rvs, ok := ctx.Value(web.KeyRequestValues).(*web.RequestValues)
	if !ok {
		return web.NewRedirectError("/home", http.StatusMovedPermanently)
	}

	var p1 string

	// If cookie does not exist, create it inside db and set it to the writer
	p1C, err := r.Cookie("p1")
	if err != nil {
		c, err := cookie.Create(ctx, h.db, rvs.Now)
		if err != nil {
			return web.NewRedirectError("/home", http.StatusMovedPermanently)
		}
		p1 = c.ID
		http.SetCookie(w, &http.Cookie{
			Name:  "p1",
			Value: p1,
			Path:  "/",
		})
	} else {
		p1 = p1C.Value
	}

	url := fmt.Sprintf("%s/sync/partner1?p1=%s", h.service2Url, p1)

	return web.RespondWithRedirect(ctx, w, r, url)
}

func (h *cookieHandlers) servePixel(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rvs, ok := ctx.Value(web.KeyRequestValues).(*web.RequestValues)
	if !ok {
		return web.RespondWithError(ctx, w, r, errors.New("request values missing from context"))
	}

	qParams := r.URL.Query()
	p1Q, ok := qParams["p1"]
	if !ok {
		return web.RespondWithError(ctx, w, r, errors.New("p1 was not found in query params"))
	}

	p2Q, ok := qParams["p2"]
	if !ok {
		return web.RespondWithError(ctx, w, r, errors.New("p2 was not found in query params"))
	}

	c, err := cookie.One(ctx, h.db, p1Q[0])
	if err != nil {
		if err == cookie.ErrNotFound {
			return web.RespondWithError(ctx, w, r, errors.New("cookie with given p1 id was not found"))
		}
		return web.RespondWithError(ctx, w, r, errors.New("cannot get cookie from db"))
	}

	// If cookie with p1 id does not have a partner cookie id then set it
	// the one got from query params
	if c.PartnerCookieID == nil {
		if err := cookie.SetPartnerCookie(ctx, h.db, c.ID, p2Q[0], rvs.Now); err != nil {
			return web.RespondWithError(ctx, w, r, errors.New("cannot set partner cookie"))
		}
	}

	// Set the header to image/gif
	w.Header().Set("Content-Type", "image/gif")

	// Decode it from raw base64 string
	data, err := base64.StdEncoding.DecodeString(base64GifPixel)
	if err != nil {
		return web.RespondWithError(ctx, w, r, errors.New("cannot decode gif"))
	}

	// Write the response to the client.
	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}
