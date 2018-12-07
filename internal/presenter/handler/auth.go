package handler

import (
	"net/http"
	"time"

	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/presenter/renderer"
	"github.com/pagient/pagient-server/internal/presenter/websocket"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

// CreateToken authenticates a user and creates a jwt token
func CreateToken(userService service.UserService, tokenService service.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userReq := &renderer.UserRequest{}
		if err := render.Bind(req, userReq); err != nil {
			render.Render(w, req, renderer.ErrBadRequest(err))
			return
		}

		user, valid, err := userService.Login(userReq.Username, userReq.Password)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		if !valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tokenAuth := jwtauth.New("HS256", []byte(config.General.Secret), nil)
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
		err = tokenService.CreateToken(token)
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
}

// DeleteToken deletes a valid jwt token
func DeleteToken(tokenService service.TokenService, wsHub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		jwtToken, _, err := jwtauth.FromContext(req.Context())
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		token, err := tokenService.ShowToken(jwtToken.Raw)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		err = tokenService.DeleteToken(token)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		wsHub.DisconnectClient(token.ID)

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			Expires:  time.Now(),
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetSessions returns all jwt tokens from a user
func GetSessions(userService service.UserService, tokenService service.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		jwtToken, _, err := jwtauth.FromContext(req.Context())
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		user, err := userService.ShowUserByToken(jwtToken.Raw)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		tokens, err := tokenService.ListTokensByUser(user.Username)
		if err != nil {
			render.Render(w, req, renderer.ErrInternalServer(err))
			return
		}

		render.RenderList(w, req, renderer.NewTokenListResponse(tokens))
	}
}
