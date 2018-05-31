package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
)

// ClientResponse is the response payload for the Client data model.
type ClientResponse struct {
	*model.Client
}

func NewClientResponse(client *model.Client) *ClientResponse {
	resp := &ClientResponse{Client: client}

	return resp
}

func (cr *ClientResponse) Render(w http.ResponseWriter, req *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

type ClientListResponse []*ClientResponse

func NewClientListResponse(clients []*model.Client) []render.Renderer {
	list := []render.Renderer{}
	for _, client := range clients {
		list = append(list, NewClientResponse(client))
	}
	return list
}
