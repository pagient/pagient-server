package context

import (
	"context"
	"net/http"

	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/rs/zerolog/log"
)

// ClientCtx middleware is used to load a Client object from the authenticated user
func ClientCtx(clientService service.ClientService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctxUser := req.Context().Value(UserKey).(*model.User)
			client, err := clientService.ShowClientByUser(ctxUser.Username)
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("get client failed")

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(req.Context(), ClientKey, client)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
