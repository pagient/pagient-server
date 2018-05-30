package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/renderer"
)

// GetClients lists all configured clients
func GetClients(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients, err := model.GetClients()
		if err != nil {
			render.Render(w, r, renderer.ErrRender(err))
			return
		}

		render.RenderList(w, r, renderer.NewClientListResponse(clients))
	}
}
