package bridge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
)

func TestSignatureTypeName(t *testing.T) {
	tests := []struct {
		sig  SignatureType
		want string
	}{
		{SigInputOutput, "input-output"},
		{SigOutput, "output-only"},
		{SigInputOnly, "input-only"},
		{SigVoid, "void"},
		{SignatureType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := signatureTypeName(tt.sig)
			if got != tt.want {
				t.Errorf("signatureTypeName(%d) = %q, want %q", tt.sig, got, tt.want)
			}
		})
	}
}

func TestGetFunctionInfo(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.info", func(ctx Context, p testInput) (testOutput, error) {
		return testOutput{}, nil
	},
		WithDescription("test function"),
		WithHTTPMethod("GET", "POST"),
		WithRenderer(func(d testOutput) templ.Component {
			return testTemplComponent(d.Result)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	info, err := b.GetFunctionInfo("test.info")
	if err != nil {
		t.Fatalf("GetFunctionInfo() error = %v", err)
	}

	if info.Name != "test.info" {
		t.Errorf("Name = %q, want %q", info.Name, "test.info")
	}
	if info.Description != "test function" {
		t.Errorf("Description = %q, want %q", info.Description, "test function")
	}
	if info.SignatureType != "input-output" {
		t.Errorf("SignatureType = %q, want %q", info.SignatureType, "input-output")
	}
	if len(info.AllowedMethods) != 2 {
		t.Errorf("AllowedMethods = %v, want [GET POST]", info.AllowedMethods)
	}
	if info.ReturnsHTML {
		t.Error("ReturnsHTML = true, want false (testOutput is not templ.Component)")
	}
	if !info.HasRenderer {
		t.Error("HasRenderer = false, want true")
	}
	if info.HTMXEndpoint != "/api/bridge/fn/test.info" {
		t.Errorf("HTMXEndpoint = %q, want %q", info.HTMXEndpoint, "/api/bridge/fn/test.info")
	}
	if info.TypeInfo.InputType == "" {
		t.Error("TypeInfo.InputType is empty, want non-empty")
	}
	if info.TypeInfo.OutputType == "" {
		t.Error("TypeInfo.OutputType is empty, want non-empty")
	}
}

func TestGetFunctionInfo_NotFound(t *testing.T) {
	b := New(WithCSRF(false))

	_, err := b.GetFunctionInfo("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent function")
	}
}

func TestGetFunctionInfo_Void(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.void", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	info, err := b.GetFunctionInfo("test.void")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if info.SignatureType != "void" {
		t.Errorf("SignatureType = %q, want %q", info.SignatureType, "void")
	}
	if info.ReturnsHTML {
		t.Error("ReturnsHTML = true, want false")
	}
	if info.HasRenderer {
		t.Error("HasRenderer = true, want false")
	}
}

func TestGetFunctionInfo_ReturnsHTML(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.html", func(ctx Context) (templ.Component, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	info, err := b.GetFunctionInfo("test.html")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if !info.ReturnsHTML {
		t.Error("ReturnsHTML = false, want true")
	}
	if info.SignatureType != "output-only" {
		t.Errorf("SignatureType = %q, want %q", info.SignatureType, "output-only")
	}
}

func TestListFunctionInfo(t *testing.T) {
	b := New(WithCSRF(false))

	_ = b.Register("fn.one", func(ctx Context) error { return nil })
	_ = b.Register("fn.two", func(ctx Context) (testOutput, error) {
		return testOutput{}, nil
	})
	_ = b.Register("fn.three", func(ctx Context, p testInput) error { return nil })

	infos := b.ListFunctionInfo()

	if len(infos) != 3 {
		t.Fatalf("len(infos) = %d, want 3", len(infos))
	}

	// Verify names are present (sorted order from ListFunctions)
	names := make(map[string]bool)
	for _, info := range infos {
		names[info.Name] = true
	}
	for _, want := range []string{"fn.one", "fn.two", "fn.three"} {
		if !names[want] {
			t.Errorf("missing function %q in ListFunctionInfo", want)
		}
	}
}

func TestListFunctionInfo_Empty(t *testing.T) {
	b := New(WithCSRF(false))

	infos := b.ListFunctionInfo()
	if len(infos) != 0 {
		t.Errorf("len(infos) = %d, want 0", len(infos))
	}
}

func TestIntrospectionHandler(t *testing.T) {
	b := New(WithCSRF(false))

	_ = b.Register("test.fn1", func(ctx Context) error { return nil })
	_ = b.Register("test.fn2", func(ctx Context, p testInput) (testOutput, error) {
		return testOutput{}, nil
	}, WithHTTPMethod("POST"), WithDescription("second function"))

	handler := b.IntrospectionHandler()
	r := httptest.NewRequest(http.MethodGet, "/api/bridge/functions", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var result struct {
		Functions []FunctionInfo `json:"functions"`
		Count     int            `json:"count"`
	}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("JSON decode error: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("count = %d, want 2", result.Count)
	}
	if len(result.Functions) != 2 {
		t.Errorf("len(functions) = %d, want 2", len(result.Functions))
	}

	// Find fn2 and verify its properties
	for _, fn := range result.Functions {
		if fn.Name == "test.fn2" {
			if fn.Description != "second function" {
				t.Errorf("fn2.Description = %q, want %q", fn.Description, "second function")
			}
			if fn.SignatureType != "input-output" {
				t.Errorf("fn2.SignatureType = %q, want %q", fn.SignatureType, "input-output")
			}
			if fn.HTMXEndpoint != "/api/bridge/fn/test.fn2" {
				t.Errorf("fn2.HTMXEndpoint = %q, want %q", fn.HTMXEndpoint, "/api/bridge/fn/test.fn2")
			}
		}
	}
}
