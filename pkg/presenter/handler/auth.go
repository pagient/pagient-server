package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/presenter/renderer"
	"github.com/pagient/pagient-server/pkg/presenter/websocket"
	"github.com/pagient/pagient-server/pkg/service"
)

// AuthHandler struct
type AuthHandler struct {
	cfg          *config.Config
	userService  service.UserService
	tokenService service.TokenService
	wsHub        *websocket.Hub
}

// NewAuthHandler initializes a AuthHandler
func NewAuthHandler(cfg *config.Config, userService service.UserService, tokenService service.TokenService, hub *websocket.Hub) *AuthHandler {
	return &AuthHandler{
		cfg:          cfg,
		userService:  userService,
		tokenService: tokenService,
		wsHub:        hub,
	}
}

// CreateToken authenticates a user and creates a jwt token
func (handler *AuthHandler) CreateToken(w http.ResponseWriter, req *http.Request) {
	userReq := &renderer.UserRequest{}
	if err := render.Bind(req, userReq); err != nil {
		render.Render(w, req, renderer.ErrBadRequest(err))
		return
	}

	user, valid, err := handler.userService.Login(userReq.Username, userReq.Password)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	if !valid {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	tokenAuth := jwtauth.New("HS256", []byte(handler.cfg.General.Secret), nil)
	jwtToken, _, err := tokenAuth.Encode(jwtauth.Claims{
		"exp": jwtauth.ExpireIn(12 * time.Hour),
	})
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	token := &model.Token{
		Raw:  jwtToken.Raw,
		User: *user,
	}
	err = handler.tokenService.Add(token)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken.Raw,
		Path:     "/",
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
	})

	render.Render(w, req, renderer.NewTokenResponse(token))
}

// DeleteToken deletes a valid jwt token
func (handler *AuthHandler) DeleteToken(w http.ResponseWriter, req *http.Request) {
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

	err = handler.tokenService.Remove(token)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	handler.wsHub.DisconnectClient(token.ID)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusNoContent)
}

// GetSessions returns all jwt tokens from a user
func (handler *AuthHandler) GetSessions(w http.ResponseWriter, req *http.Request) {
	jwtToken, _, err := jwtauth.FromContext(req.Context())
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	user, err := handler.userService.GetByToken(jwtToken.Raw)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	tokens, err := handler.tokenService.GetByUser(user.Username)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	render.RenderList(w, req, renderer.NewTokenListResponse(tokens))
}
