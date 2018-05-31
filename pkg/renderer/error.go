package renderer

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http renderer status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render renders the ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest represents a 400 error
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender represents a 422 error caused by the renderer
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering renderer.",
		ErrorText:      err.Error(),
	}
}

// ErrValidation represents a 422 error caused by validation
func ErrValidation(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Validation error.",
		ErrorText:      err.Error(),
	}
}

// ErrInternalServer represents a 500 error
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

// ErrGateway represents a 504 error caused by unresponsive auxiliary servers
func ErrGateway(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 504,
		StatusText:     "Error receiving response from server.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound represents a 404 error
var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
