package context

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/renderer"
)

// AuthCtx middleware is used to load a Client object from
// the basic auth headers passed through as the request. In case
// the Client could not be found, we stop here and return a 422.
func AuthCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var client *model.Client
		var err error

		username, _, ok := r.BasicAuth()
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
		}

		client, err = model.GetClient(username)

		if err != nil {
			render.Render(w, r, renderer.ErrRender(err))
			return
		}

		ctx := context.WithValue(r.Context(), clientKey, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
