package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/service"
)

// GetPagers lists all configured pagers
func GetPagers(pagerService service.PagerService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pagers, err := pagerService.GetAll()
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewPagerListResponse(pagers))
	}
}
