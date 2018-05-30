package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
)

// PagerResponse is the response payload for the Pager data model.
type PagerResponse struct {
	*model.Pager
}

func NewPagerResponse(Pager *model.Pager) *PagerResponse {
	resp := &PagerResponse{Pager: Pager}

	return resp
}

func (pr *PagerResponse) Render(w http.ResponseWriter, req *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

type PagerListResponse []*PagerResponse

func NewPagerListResponse(Pagers []*model.Pager) []render.Renderer {
	list := []render.Renderer{}
	for _, Pager := range Pagers {
		list = append(list, NewPagerResponse(Pager))
	}
	return list
}
