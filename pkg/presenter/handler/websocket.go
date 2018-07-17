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

// WebsocketHandler struct
type WebsocketHandler struct {
	cfg          *config.Config
	tokenService service.TokenService
	wsHub        *websocket.Hub
	wsUpgrader   ws.Upgrader
}

// NewWebsocketHandler initializes a WebsocketHandler
func NewWebsocketHandler(cfg *config.Config, tokenService service.TokenService, hub *websocket.Hub) *WebsocketHandler {
	upgrader := ws.Upgrader{
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

	return &WebsocketHandler{
		cfg:          cfg,
		tokenService: tokenService,
		wsHub:        hub,
		wsUpgrader:   upgrader,
	}
}

// ServeWebsocket establishes the websocket connection per client
func (handler *WebsocketHandler) ServeWebsocket(w http.ResponseWriter, req *http.Request) {
	conn, err := handler.wsUpgrader.Upgrade(w, req, nil)
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

	token, err := handler.tokenService.Get(jwtToken.Raw)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	client := websocket.NewClient(token.ID, handler.wsHub, conn)
	handler.wsHub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
}
