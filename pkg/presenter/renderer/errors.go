package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"`      // low-level runtime error
	HTTPStatusCode int   `json:"status"` // http renderer status code

	Message   string `json:"message"`         // user-level status message
	ErrorText string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render renders the ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, req *http.Request) error {
	if e.HTTPStatusCode == 500 {
		log.Error().
			Err(e.Err).
			Msg("")
	}

	render.Status(req, e.HTTPStatusCode)
	return nil
}

// ErrBadRequest represents a 400 error
func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		Message:        "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrConflict represents a 409 error
func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 409,
		Message:        "Resouce already exists.",
	}
}

// ErrValidation represents a 422 error caused by validation
func ErrValidation(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Message:        "Validation error.",
		ErrorText:      err.Error(),
	}
}

// ErrInternalServer represents a 500 error
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		Message:        "Internal server error.",
	}
}

// ErrGateway represents a 504 error caused by unresponsive auxiliary servers
func ErrGateway(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 504,
		Message:        "Error receiving response from server.",
	}
}

// ErrNotFound represents a 404 error
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, Message: "Resource not found."}
