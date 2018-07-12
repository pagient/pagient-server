package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/service"
)

// ClientHandler struct
type ClientHandler struct {
	clientService service.ClientService
}

// NewClientHandler initializes a ClientHandler
func NewClientHandler(clientService service.ClientService) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
	}
}

// GetClients lists all configured clients
func (handler *ClientHandler) GetClients(w http.ResponseWriter, req *http.Request) {
	clients, err := handler.clientService.GetAll()
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	render.RenderList(w, req, renderer.NewClientListResponse(clients))
}
