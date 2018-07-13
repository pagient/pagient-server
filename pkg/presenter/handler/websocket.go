package handler

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	ws "github.com/gorilla/websocket"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/presenter/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
)

// WebsocketHandler struct
type WebsocketHandler struct {
	cfg        *config.Config
	wsHub      *websocket.Hub
	wsUpgrader ws.Upgrader
}

// NewWebsocketHandler initializes a WebsocketHandler
func NewWebsocketHandler(cfg *config.Config, hub *websocket.Hub) *WebsocketHandler {
	upgrader := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(req *http.Request) bool {
			hostUrl, err := url.Parse(cfg.Server.Host)
			if err != nil {
				return false
			}

			if hostUrl.String() == req.Header.Get("Origin") && hostUrl.Host == req.Host {
				return true
			}
			return false
		},
	}

	return &WebsocketHandler{
		cfg:        cfg,
		wsHub:      hub,
		wsUpgrader: upgrader,
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

	token, _, err := jwtauth.FromContext(req.Context())
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	client := websocket.NewClient(token.Signature, handler.wsHub, conn)
	handler.wsHub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
}
