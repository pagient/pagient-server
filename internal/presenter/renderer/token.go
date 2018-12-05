package renderer

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/model"

	"github.com/go-chi/render"
)

// TokenResponse is the response payload for the token data model
type TokenResponse struct {
	Token string `json:"token"`
}

// NewTokenResponse creates a new token response from token model
func NewTokenResponse(token *model.Token) *TokenResponse {
	resp := &TokenResponse{Token: token.Raw}

	return resp
}

// Render preprocesses the response before marshalling
func (cr *TokenResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

// TokenListResponse is the list response payload for the token data model
type TokenListResponse []*TokenResponse

// NewTokenListResponse creates a new token list response from multiple token models
func NewTokenListResponse(tokens []*model.Token) []render.Renderer {
	list := make([]render.Renderer, len(tokens))
	for i, token := range tokens {
		list[i] = NewTokenResponse(token)
	}
	return list
}
