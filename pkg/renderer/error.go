package renderer

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"status"` // http renderer status code

	Message    string `json:"message"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render renders the ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	log.Error().Err(e.Err).Msg("")
	render.Status(r, e.HTTPStatusCode)
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

func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 409,
		Message:        "Resouce already exists.",
		ErrorText:      err.Error(),
	}
}

// ErrRender represents a 422 error caused by the renderer
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Message:        "Error rendering renderer.",
		ErrorText:      err.Error(),
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
		ErrorText:      err.Error(),
	}
}

// ErrGateway represents a 504 error caused by unresponsive auxiliary servers
func ErrGateway(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 504,
		Message:        "Error receiving response from server.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound represents a 404 error
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, Message: "Resource not found."}
