package bridge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBridge_Execute(t *testing.T) {
	b := New()

	// Register a test function
	b.Register("echo", func(ctx Context, input testInput) (testOutput, error) {
		return testOutput{Result: input.Name}, nil
	})

	// Create context
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := NewContext(req)

	// Execute function
	params := json.RawMessage(`{"name":"test"}`)
	result := b.execute(ctx, "echo", params)

	if result.Error != nil {
		t.Fatalf("execute() error = %v", result.Error)
	}

	output, ok := result.Result.(testOutput)
	if !ok {
		t.Fatalf("result type = %T, want testOutput", result.Result)
	}

	if output.Result != "test" {
		t.Errorf("output.Result = %s, want test", output.Result)
	}
}

func TestBridge_Execute_Timeout(t *testing.T) {
	b := New(WithTimeout(50 * time.Millisecond))

	// Use a struct with no required fields for timeout test
	type emptyInput struct {
		Value int `json:"value,omitempty"`
	}

	// Register function that takes too long
	b.Register("slow", func(ctx Context, input emptyInput) (testOutput, error) {
		time.Sleep(500 * time.Millisecond)
		return testOutput{Result: "done"}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := NewContext(req)

	// Pass empty JSON object to avoid required field validation
	params := json.RawMessage(`{}`)
	result := b.execute(ctx, "slow", params)

	if result.Error == nil {
		t.Fatal("execute() should timeout")
	}

	if result.Error.Code != ErrCodeTimeout {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeTimeout)
	}
}

func TestBridge_Execute_Panic(t *testing.T) {
	b := New()

	// Use a struct with no required fields for panic test
	type emptyInput struct {
		Value int `json:"value,omitempty"`
	}

	// Register function that panics
	b.Register("panic", func(ctx Context, input emptyInput) (testOutput, error) {
		panic("intentional panic")
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := NewContext(req)

	// Pass empty JSON object to avoid required field validation
	params := json.RawMessage(`{}`)
	result := b.execute(ctx, "panic", params)

	if result.Error == nil {
		t.Error("execute() should catch panic")
	}

	if result.Error.Code != ErrCodeInternal {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeInternal)
	}
}

func TestBridge_CallBatch(t *testing.T) {
	b := New()

	// Register test functions
	b.Register("add", func(ctx Context, input struct {
		A int `json:"a"`
		B int `json:"b"`
	}) (struct {
		Sum int `json:"sum"`
	}, error) {
		return struct {
			Sum int `json:"sum"`
		}{Sum: input.A + input.B}, nil
	})

	b.Register("multiply", func(ctx Context, input struct {
		A int `json:"a"`
		B int `json:"b"`
	}) (struct {
		Product int `json:"product"`
	}, error) {
		return struct {
			Product int `json:"product"`
		}{Product: input.A * input.B}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := NewContext(req)

	requests := []Request{
		{
			JSONRPC: "2.0",
			ID:      "1",
			Method:  "add",
			Params:  json.RawMessage(`{"a":2,"b":3}`),
		},
		{
			JSONRPC: "2.0",
			ID:      "2",
			Method:  "multiply",
			Params:  json.RawMessage(`{"a":4,"b":5}`),
		},
	}

	responses := b.CallBatch(ctx, requests)

	if len(responses) != 2 {
		t.Fatalf("len(responses) = %d, want 2", len(responses))
	}

	// Check first response
	if responses[0].Error != nil {
		t.Errorf("responses[0] has error: %v", responses[0].Error)
	}

	// Check second response
	if responses[1].Error != nil {
		t.Errorf("responses[1] has error: %v", responses[1].Error)
	}
}

func TestBridge_CallBatch_ExceedsMaxSize(t *testing.T) {
	b := New(WithMaxBatchSize(2))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := NewContext(req)

	requests := []Request{
		{Method: "test1"},
		{Method: "test2"},
		{Method: "test3"}, // Exceeds max
	}

	responses := b.CallBatch(ctx, requests)

	if len(responses) != 1 {
		t.Fatalf("len(responses) = %d, want 1 (error response)", len(responses))
	}

	if responses[0].Error == nil {
		t.Error("expected error for exceeding batch size")
	}
}
