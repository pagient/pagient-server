package renderer

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/model"

	"github.com/go-chi/render"
)

// ClientResponse is the response payload for the client data model
type ClientResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// NewClientResponse creates a new client response from client model
func NewClientResponse(client *model.Client) *ClientResponse {
	resp := &ClientResponse{ID: client.ID, Name: client.Name}

	return resp
}

// Render preprocesses the response before marshalling
func (cr *ClientResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

// ClientListResponse is the list response payload for the client data model
type ClientListResponse []*ClientResponse

// NewClientListResponse creates a new client list response from multiple client models
func NewClientListResponse(clients []*model.Client) []render.Renderer {
	list := make([]render.Renderer, len(clients))
	for i, client := range clients {
		list[i] = NewClientResponse(client)
	}
	return list
}
