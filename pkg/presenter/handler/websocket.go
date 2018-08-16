package handler

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	ws "github.com/gorilla/websocket"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/presenter/websocket"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/rs/zerolog/log"
	"net/url"
)

// ServeWebsocket establishes the websocket connection per client
func ServeWebsocket(cfg *config.Config, tokenService service.TokenService, wsHub *websocket.Hub) http.HandlerFunc {
	wsUpgrader := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(req *http.Request) bool {
			hostURL, err := url.Parse(cfg.Server.Host)
			if err != nil {
				return false
			}

			if hostURL.String() == req.Header.Get("Origin") && hostURL.Host == req.Host {
				return true
			}
			return false
		},
	}

	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Error().
				Err(err).
				Msg("websocket connection could not be established")

			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		jwtToken, _, err := jwtauth.FromContext(req.Context())
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		token, err := tokenService.Get(jwtToken.Raw)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		client := websocket.NewClient(token.ID, wsHub, conn)
		wsHub.Register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
	}
}
