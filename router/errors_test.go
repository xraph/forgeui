package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestLoaderErrorType(t *testing.T) {
	err := &LoaderError{
		Status:  404,
		Message: "Not found",
	}

	if err.Error() != "Not found" {
		t.Errorf("Expected 'Not found', got '%s'", err.Error())
	}
}

func TestLoaderErrorWithWrappedError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := &LoaderError{
		Status:  500,
		Message: "Server error",
		Err:     innerErr,
	}

	if !strings.Contains(err.Error(), "Server error") {
		t.Error("Expected error message to contain 'Server error'")
	}

	if !strings.Contains(err.Error(), "inner error") {
		t.Error("Expected error message to contain 'inner error'")
	}

	// Test Unwrap
	if errors.Unwrap(err) != innerErr {
		t.Error("Expected Unwrap to return inner error")
	}
}

func TestError404(t *testing.T) {
	err := Error404("Page not found")

	loaderErr, ok := err.(*LoaderError)
	if !ok {
		t.Fatal("Expected LoaderError")
	}

	if loaderErr.Status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", loaderErr.Status)
	}

	if loaderErr.Message != "Page not found" {
		t.Errorf("Expected message 'Page not found', got '%s'", loaderErr.Message)
	}
}

func TestError403(t *testing.T) {
	err := Error403("Access denied")

	loaderErr, ok := err.(*LoaderError)
	if !ok {
		t.Fatal("Expected LoaderError")
	}

	if loaderErr.Status != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", loaderErr.Status)
	}
}

func TestError500(t *testing.T) {
	innerErr := errors.New("database error")
	err := Error500("Internal error", innerErr)

	loaderErr, ok := err.(*LoaderError)
	if !ok {
		t.Fatal("Expected LoaderError")
	}

	if loaderErr.Status != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", loaderErr.Status)
	}

	if loaderErr.Err != innerErr {
		t.Error("Expected wrapped error to be preserved")
	}
}

func TestDefaultErrorPage(t *testing.T) {
	testCases := []struct {
		status  int
		title   string
		message string
	}{
		{404, "404 - Page Not Found", "looking for"},
		{403, "403 - Forbidden", "permission"},
		{401, "401 - Unauthorized", "logged in"},
		{500, "500 - Internal Server Error", "went wrong"},
		{408, "408 - Request Timeout", "took too long"},
		{418, "418 - Error", "error occurred"},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			handler := DefaultErrorPage(tc.status)
			ctx := &PageContext{
				ResponseWriter: httptest.NewRecorder(),
			}

			node, err := handler(ctx)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Render and check output
			var buf strings.Builder

			_ = node.Render(&buf)
			output := buf.String()

			if !strings.Contains(output, tc.title) {
				t.Errorf("Expected output to contain '%s', got: %s", tc.title, output)
			}

			if !strings.Contains(output, tc.message) {
				t.Errorf("Expected output to contain '%s', got: %s", tc.message, output)
			}
		})
	}
}

func TestSetErrorPage(t *testing.T) {
	r := New()

	// Set custom error page
	customHandler := func(ctx *PageContext) (g.Node, error) {
		return g.Text("Custom Error"), nil
	}
	r.SetErrorPage(404, customHandler)

	// Verify it was set
	handler := r.getErrorPage(404)
	ctx := &PageContext{
		ResponseWriter: httptest.NewRecorder(),
	}

	node, _ := handler(ctx)

	var buf strings.Builder

	_ = node.Render(&buf)

	if buf.String() != "Custom Error" {
		t.Errorf("Expected 'Custom Error', got '%s'", buf.String())
	}
}

func TestGetErrorPageDefault(t *testing.T) {
	r := New()

	// Get error page for status that wasn't set
	handler := r.getErrorPage(404)
	if handler == nil {
		t.Error("Expected default error handler to be returned")
	}

	ctx := &PageContext{
		ResponseWriter: httptest.NewRecorder(),
	}

	node, _ := handler(ctx)

	var buf strings.Builder

	_ = node.Render(&buf)

	if !strings.Contains(buf.String(), "404") {
		t.Error("Expected default 404 error page")
	}
}

func TestCustomErrorPageIntegration(t *testing.T) {
	r := New()

	// Set custom 500 error page
	r.SetErrorPage(500, func(ctx *PageContext) (g.Node, error) {
		return g.Text("Oops! Something went wrong."), nil
	})

	// Create a route that triggers the error
	r.Get("/error", func(ctx *PageContext) (g.Node, error) {
		return nil, Error500("Test error", nil)
	})

	// Test the route
	req := httptest.NewRequest(MethodGet, "/error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// The error should be handled by the error handler, not the custom error page
	// because the error is returned from the handler, not the loader
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
