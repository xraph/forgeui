package router

import (
	"fmt"
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// LoaderError represents an error that occurred during data loading
type LoaderError struct {
	Status  int
	Message string
	Err     error
}

// Error implements the error interface
func (e *LoaderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *LoaderError) Unwrap() error {
	return e.Err
}

// SetErrorPage registers a custom error page for a specific status code
func (r *Router) SetErrorPage(status int, handler PageHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.errorPages == nil {
		r.errorPages = make(map[int]PageHandler)
	}

	r.errorPages[status] = handler
}

// getErrorPage retrieves the error page handler for a status code
func (r *Router) getErrorPage(status int) PageHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.errorPages != nil {
		if handler, ok := r.errorPages[status]; ok {
			return handler
		}
	}

	// Return default error page
	return DefaultErrorPage(status)
}

// DefaultErrorPage returns a default error page for the given status code
func DefaultErrorPage(status int) PageHandler {
	return func(ctx *PageContext) (g.Node, error) {
		ctx.ResponseWriter.WriteHeader(status)

		var title, message string
		switch status {
		case http.StatusNotFound:
			title = "404 - Page Not Found"
			message = "The page you're looking for doesn't exist."
		case http.StatusForbidden:
			title = "403 - Forbidden"
			message = "You don't have permission to access this resource."
		case http.StatusUnauthorized:
			title = "401 - Unauthorized"
			message = "You must be logged in to access this page."
		case http.StatusInternalServerError:
			title = "500 - Internal Server Error"
			message = "Something went wrong on our end."
		case http.StatusRequestTimeout:
			title = "408 - Request Timeout"
			message = "The request took too long to complete."
		default:
			title = fmt.Sprintf("%d - Error", status)
			message = "An error occurred while processing your request."
		}

		return html.Div(
			html.Class("error-page"),
			html.H1(g.Text(title)),
			html.P(g.Text(message)),
		), nil
	}
}

// Error404 is a convenience function for returning a 404 error
func Error404(message string) error {
	return &LoaderError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

// Error403 is a convenience function for returning a 403 error
func Error403(message string) error {
	return &LoaderError{
		Status:  http.StatusForbidden,
		Message: message,
	}
}

// Error500 is a convenience function for returning a 500 error
func Error500(message string, err error) error {
	return &LoaderError{
		Status:  http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

