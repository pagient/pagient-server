package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/pagient/pagient-api/pkg/config"
	"github.com/pagient/pagient-api/pkg/model"
	"github.com/pagient/pagient-api/pkg/presenter/renderer"
	"github.com/pagient/pagient-api/pkg/presenter/websocket"
	"github.com/pagient/pagient-api/pkg/service"
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
	data := &renderer.UserRequest{}
	if err := render.Bind(req, data); err != nil {
		render.Render(w, req, renderer.ErrBadRequest(err))
		return
	}

	user := data.User
	valid, err := handler.userService.Login(user.Username, user.Password)
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	if !valid {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	tokenAuth := jwtauth.New("HS256", []byte(handler.cfg.Server.SecretKey), nil)
	_, tokenString, err := tokenAuth.Encode(jwtauth.Claims{
		"user": user.Username,
		"exp":  jwtauth.ExpireIn(12 * time.Hour),
	})
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	err = handler.tokenService.Add(&model.Token{
		Token: tokenString,
		User:  user.Username,
	})
	if err != nil && !service.IsModelNotExistErr(err) {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
	})

	render.Render(w, req, renderer.NewTokenResponse(&model.Token{
		Token: tokenString,
	}))
}

// DeleteToken deletes a valid jwt token
func (handler *AuthHandler) DeleteToken(w http.ResponseWriter, req *http.Request) {
	token, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	username, ok := claims.Get("user")
	if !ok {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	err = handler.tokenService.Remove(&model.Token{
		Token: token.Raw,
		User:  username.(string),
	})
	if err != nil {
		render.Render(w, req, renderer.ErrInternalServer(err))
		return
	}

	handler.wsHub.Disconnect(token.Signature)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusNoContent)
}
