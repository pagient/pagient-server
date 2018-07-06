package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/presenter/renderer"
	"github.com/pagient/pagient-api/pkg/service"
)

// PagerHandler struct
type PagerHandler struct {
	pagerService service.PagerService
}

// NewPagerHandler initializes a PagerHandler
func NewPagerHandler(pagerService service.PagerService) *PagerHandler {
	return &PagerHandler{
		pagerService: pagerService,
	}
}

// GetPagers lists all configured pagers
func (handler *PagerHandler) GetPagers(w http.ResponseWriter, req *http.Request) {
	pagers, err := handler.pagerService.GetAll()
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	render.RenderList(w, req, renderer.NewPagerListResponse(pagers))
}
