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
	if e.Err != nil {
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
		HTTPStatusCode: http.StatusBadRequest,
		Message:        http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
	}
}

// ErrConflict represents a 409 error
func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusConflict,
		Message:        http.StatusText(http.StatusConflict),
	}
}

// ErrValidation represents a 422 error caused by validation
func ErrValidation(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		Message:        http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
	}
}

// ErrInternalServer represents a 500 error
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        http.StatusText(http.StatusInternalServerError),
	}
}

// ErrGateway represents a 504 error caused by unresponsive auxiliary servers
func ErrGateway(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusGatewayTimeout,
		Message:        http.StatusText(http.StatusGatewayTimeout),
	}
}

// ErrUnauthorized represents a 401 error
var ErrUnauthorized = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}

// ErrNotFound represents a 404 error
var ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
