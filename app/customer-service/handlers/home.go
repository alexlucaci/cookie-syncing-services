package handlers

import (
	"context"
	"github.com/alexlucaci/cookie-syncing-services/foundation/web"
	"net/http"
)

type homeHandlers struct {
	service1Url string
}

type HomePage struct {
	WhoAmILink string
}

func (h *homeHandlers) home(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := HomePage{WhoAmILink: h.service1Url + "/whoami"}
	return web.RenderTemplate(w, "business/templates/home.html", page)
}
