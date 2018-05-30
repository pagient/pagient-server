package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/renderer"
)

// GetPagers lists all available pagers
func GetPagers(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pagers, err := model.GetPagers()
		if err != nil {
			render.Render(w, req, renderer.ErrRender(err))
			return
		}

		render.RenderList(w, req, renderer.NewPagerListResponse(pagers))
	}
}
