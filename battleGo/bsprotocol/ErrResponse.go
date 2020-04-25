package bsprotocol

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse allows the renderer to send error responses
type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string      `json:"status"`
	AppCode    int         `json:"code,omitempty"`
	ErrorText  string      `json:"error,omitempty"`
	Body       interface{} `json:"request,omitempty"`
}

// ErrNotFound is the standard response object for a 404 case
var ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "Resource not found."}

// Render sets the status code header and satisfies the renderer interface
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Err: %v\n", e.Err)
	render.Status(r, e.HTTPStatusCode)
	return nil

}

// ErrInternalError a bad request response renderer
func ErrInternalError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
		ErrorText:      err.Error(),
	}
}

// ErrForbidden ...
func ErrForbidden(err error, body interface{}) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     http.StatusText(http.StatusForbidden),
		ErrorText:      err.Error(),
		Body:           body,
	}
}

// ErrBadRequest Returns a bad request response renderer
func ErrBadRequest(err error, body interface{}) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
		Body:           body,
	}
}

// ErrPreconditionFail ...
func ErrPreconditionFail(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusPreconditionFailed,
		StatusText:     http.StatusText(http.StatusPreconditionFailed),
		ErrorText:      err.Error(),
	}
}
