package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/service"
)

// GetClients lists all configured clients
func GetClients(clientService service.ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		clients, err := clientService.GetAll()
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewClientListResponse(clients))
	}
}
