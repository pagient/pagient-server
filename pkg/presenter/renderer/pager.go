package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/model"
)

// PagerResponse is the response payload for the pager data model
type PagerResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// NewPagerResponse creates a new pager response from pager model
func NewPagerResponse(pager *model.Pager) *PagerResponse {
	resp := &PagerResponse{ID: pager.ID, Name: pager.Name}

	return resp
}

// Render preprocesses the response before marshalling
func (pr *PagerResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

// PagerListResponse is the list response payload for the pager data model
type PagerListResponse []*PagerResponse

// NewPagerListResponse creates a new pager list response from multiple pager models
func NewPagerListResponse(pagers []*model.Pager) []render.Renderer {
	list := make([]render.Renderer, len(pagers))
	for i, pager := range pagers {
		list[i] = NewPagerResponse(pager)
	}
	return list
}
