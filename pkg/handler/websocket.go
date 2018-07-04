package handler

import (
	"net/http"

	"github.com/go-chi/render"
	ws "github.com/gorilla/websocket"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/renderer"
	"github.com/pagient/pagient-api/pkg/websocket"
	"github.com/rs/zerolog/log"
)

// GetPagers lists all available pagers
func ServeWebsocket(cfg *config.Config, hub *websocket.Hub) http.HandlerFunc {
	upgrader := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Error().
				Err(err).
				Msg("websocket connection could not be established")
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		client := websocket.NewClient(hub, conn)
		hub.Register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump()
	}
}
