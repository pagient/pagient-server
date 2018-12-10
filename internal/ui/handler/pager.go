package handler

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/service"
	"github.com/pagient/pagient-server/internal/ui/renderer"

	"github.com/go-chi/render"
)

// GetPagers lists all configured pagers
func GetPagers(pagerService service.PagerService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pagers, err := pagerService.ListPagers()
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewPagerListResponse(pagers))
	}
}
