package renderer

import (
	"net/http"

	"github.com/pagient/pagient-server/internal/model"
)

// UserRequest is the request payload for user data model
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Bind postprocesses the decoding of the request body
func (pr *UserRequest) Bind(r *http.Request) error {
	return nil
}

// GetModel returns a User model
func (pr *UserRequest) GetModel() *model.User {
	return &model.User{
		Username: pr.Username,
		Password: pr.Password,
	}
}
