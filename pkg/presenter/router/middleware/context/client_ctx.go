package context

import (
	"context"
	"net/http"

	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/service"
	"github.com/rs/zerolog/log"
)

func ClientCtx(clientService service.ClientService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctxUser := req.Context().Value("user").(*model.User)
			client, err := clientService.GetByUser(ctxUser)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("get client failed")

				http.Error(w, http.StatusText(500), 500)
				return
			}

			ctx := context.WithValue(req.Context(), "client", client)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
