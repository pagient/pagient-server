package renderer

import (
	"net/http"

	"github.com/pagient/pagient-api/pkg/model"
)

// TokenResponse is the response payload for the token data model
type TokenResponse struct {
	*model.Token
}

// NewTokenResponse creates a new token response from token model
func NewTokenResponse(token *model.Token) *TokenResponse {
	resp := &TokenResponse{Token: token}

	return resp
}

// Render preprocesses the response before marshalling
func (cr *TokenResponse) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
