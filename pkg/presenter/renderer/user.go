package renderer

import (
	"net/http"

	"github.com/pagient/pagient-api/pkg/model"
)

// UserRequest is the request payload for user data model
type UserRequest struct {
	*model.User
}

// Bind postprocesses the decoding of the request body
func (pr *UserRequest) Bind(r *http.Request) error {
	return nil
}
