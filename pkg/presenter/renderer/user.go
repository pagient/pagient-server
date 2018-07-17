package renderer

import (
	"github.com/pagient/pagient-server/pkg/model"
	"net/http"
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
