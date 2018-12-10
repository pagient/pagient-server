package handler

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/service"
	"github.com/pagient/pagient-server/internal/ui/renderer"

	"github.com/go-chi/render"
)

// GetClients lists all configured clients
func GetClients(clientService service.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		clients, err := clientService.ListClients()
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewClientListResponse(clients))
	}
}
