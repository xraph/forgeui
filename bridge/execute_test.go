package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestBridge_Execute(t *testing.T) {
	b := New()

	// Register a test function
	_ = b.Register("echo", func(ctx Context, input testInput) (testOutput, error) {
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
	_ = b.Register("slow", func(ctx Context, input emptyInput) (testOutput, error) {
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
	_ = b.Register("panic", func(ctx Context, input emptyInput) (testOutput, error) {
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
	_ = b.Register("add", func(ctx Context, input struct {
		A int `json:"a"`
		B int `json:"b"`
	}) (struct {
		Sum int `json:"sum"`
	}, error) {
		return struct {
			Sum int `json:"sum"`
		}{Sum: input.A + input.B}, nil
	})

	_ = b.Register("multiply", func(ctx Context, input struct {
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

// --- executeDirect tests ---

func TestExecuteDirect_SigInputOutput(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.echo", func(ctx Context, params testInput) (testOutput, error) {
		return testOutput{Result: "echo:" + params.Name}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.echo")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	paramValue := reflect.ValueOf(testInput{Name: "hello"})
	result := b.executeDirect(ctx, fn, paramValue)

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	out := result.Result.(testOutput)
	if out.Result != "echo:hello" {
		t.Errorf("result = %q, want %q", out.Result, "echo:hello")
	}
}

func TestExecuteDirect_SigOutput(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.getAll", func(ctx Context) (testOutput, error) {
		return testOutput{Result: "all"}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.getAll")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeDirect(ctx, fn, reflect.Value{})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	out := result.Result.(testOutput)
	if out.Result != "all" {
		t.Errorf("result = %q, want %q", out.Result, "all")
	}
}

func TestExecuteDirect_SigInputOnly(t *testing.T) {
	b := New(WithCSRF(false))

	var received string
	err := b.Register("test.action", func(ctx Context, params testInput) error {
		received = params.Name
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.action")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	paramValue := reflect.ValueOf(testInput{Name: "doIt"})
	result := b.executeDirect(ctx, fn, paramValue)

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if result.Result != nil {
		t.Errorf("result = %v, want nil", result.Result)
	}
	if received != "doIt" {
		t.Errorf("received = %q, want %q", received, "doIt")
	}
}

func TestExecuteDirect_SigVoid(t *testing.T) {
	b := New(WithCSRF(false))

	var called bool
	err := b.Register("test.ping", func(ctx Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.ping")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeDirect(ctx, fn, reflect.Value{})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if !called {
		t.Error("function was not called")
	}
}

func TestExecuteDirect_ReturnsError(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.fail", func(ctx Context) error {
		return fmt.Errorf("something broke")
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.fail")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeDirect(ctx, fn, reflect.Value{})

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Error.Code != ErrCodeInternal {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeInternal)
	}
	if result.Error.Message != "something broke" {
		t.Errorf("error message = %q, want %q", result.Error.Message, "something broke")
	}
}

func TestExecuteDirect_ReturnsBridgeError(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.bridgeErr", func(ctx Context) (testOutput, error) {
		return testOutput{}, NewError(ErrCodeForbidden, "access denied")
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.bridgeErr")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeDirect(ctx, fn, reflect.Value{})

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Error.Code != ErrCodeForbidden {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeForbidden)
	}
}

func TestExecuteDirect_StrictValidation(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.strict", func(ctx Context, params testInput) (testOutput, error) {
		return testOutput{Result: params.Name}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.strict")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))

	// Empty Name should fail strict validation (no omitempty)
	paramValue := reflect.ValueOf(testInput{Name: "", Value: 42})
	result := b.executeDirect(ctx, fn, paramValue)

	if result.Error == nil {
		t.Fatal("expected validation error, got nil")
	}
	if result.Error.Code != ErrCodeInvalidParams {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeInvalidParams)
	}
}

func TestExecuteDirect_LaxValidation(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.lax", func(ctx Context, params testInput) (testOutput, error) {
		return testOutput{Result: "ok"}, nil
	}, WithLaxValidation())
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.lax")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))

	// Empty fields should pass with lax validation
	paramValue := reflect.ValueOf(testInput{Name: "", Value: 0})
	result := b.executeDirect(ctx, fn, paramValue)

	if result.Error != nil {
		t.Fatalf("unexpected error (lax should pass): %v", result.Error)
	}
}

func TestExecuteDirect_Hooks(t *testing.T) {
	b := New(WithCSRF(false))

	var beforeCalled, afterCalled, successCalled int32
	b.GetHooks().Register(BeforeCall, func(ctx Context, data HookData) {
		atomic.AddInt32(&beforeCalled, 1)
	})
	b.GetHooks().Register(AfterCall, func(ctx Context, data HookData) {
		atomic.AddInt32(&afterCalled, 1)
	})
	b.GetHooks().Register(OnSuccess, func(ctx Context, data HookData) {
		atomic.AddInt32(&successCalled, 1)
	})

	err := b.Register("test.hooked", func(ctx Context) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.hooked")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	b.executeDirect(ctx, fn, reflect.Value{})

	// Hooks run in goroutines, give them time to complete
	time.Sleep(50 * time.Millisecond)

	if atomic.LoadInt32(&beforeCalled) != 1 {
		t.Errorf("BeforeCall called %d times, want 1", beforeCalled)
	}
	if atomic.LoadInt32(&afterCalled) != 1 {
		t.Errorf("AfterCall called %d times, want 1", afterCalled)
	}
	if atomic.LoadInt32(&successCalled) != 1 {
		t.Errorf("OnSuccess called %d times, want 1", successCalled)
	}
}

func TestExecuteDirect_ErrorHook(t *testing.T) {
	b := New(WithCSRF(false))

	var errorCalled int32
	b.GetHooks().Register(OnError, func(ctx Context, data HookData) {
		atomic.AddInt32(&errorCalled, 1)
	})

	err := b.Register("test.errHook", func(ctx Context) error {
		return fmt.Errorf("boom")
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.errHook")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	b.executeDirect(ctx, fn, reflect.Value{})

	// Hooks run in goroutines, give them time to complete
	time.Sleep(50 * time.Millisecond)

	if atomic.LoadInt32(&errorCalled) != 1 {
		t.Errorf("OnError called %d times, want 1", errorCalled)
	}
}

func TestExecuteDirect_ReturnsTemplComponent(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.html", func(ctx Context) (templ.Component, error) {
		return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			_, err := fmt.Fprintf(w, "<div>hello</div>")
			return err
		}), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.html")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeDirect(ctx, fn, reflect.Value{})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if _, ok := result.Result.(templ.Component); !ok {
		t.Errorf("result type = %T, want templ.Component", result.Result)
	}
}

// --- executeWithTimeout tests ---

func TestExecuteWithTimeout_FunctionTimeoutOverride(t *testing.T) {
	b := New(WithCSRF(false), WithTimeout(5*time.Second))

	err := b.Register("test.fastTimeout", func(ctx Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	}, WithFunctionTimeout(50*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.fastTimeout")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeWithTimeout(ctx, fn, reflect.Value{})

	if result.Error == nil {
		t.Fatal("expected timeout error from function-specific timeout")
	}
	if result.Error.Code != ErrCodeTimeout {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeTimeout)
	}
}

func TestExecuteWithTimeout_NilPtrResult(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.nilResult", func(ctx Context) (*testOutput, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.nilResult")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeWithTimeout(ctx, fn, reflect.Value{})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if result.Result != nil {
		t.Errorf("result = %v, want nil", result.Result)
	}
}

func TestExecuteWithTimeout_SigInputOnly_Error(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.inputErr", func(ctx Context, p testInput) error {
		return fmt.Errorf("input failed: %s", p.Name)
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.inputErr")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeWithTimeout(ctx, fn, reflect.ValueOf(testInput{Name: "bad"}))

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Error.Message != "input failed: bad" {
		t.Errorf("error message = %q, want %q", result.Error.Message, "input failed: bad")
	}
}

func TestExecuteWithTimeout_SigInputOnly_BridgeError(t *testing.T) {
	b := New(WithCSRF(false))

	err := b.Register("test.inputBridgeErr", func(ctx Context, p testInput) error {
		return NewError(ErrCodeForbidden, "not allowed")
	})
	if err != nil {
		t.Fatal(err)
	}

	fn, _ := b.GetFunction("test.inputBridgeErr")
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.executeWithTimeout(ctx, fn, reflect.ValueOf(testInput{Name: "x"}))

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Error.Code != ErrCodeForbidden {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeForbidden)
	}
}

// --- canBeNil tests ---

func TestCanBeNil(t *testing.T) {
	tests := []struct {
		name string
		val  reflect.Value
		want bool
	}{
		{"pointer", reflect.ValueOf((*string)(nil)), true},
		{"slice", reflect.ValueOf([]string(nil)), true},
		{"map", reflect.ValueOf(map[string]int(nil)), true},
		{"chan", reflect.ValueOf((chan int)(nil)), true},
		{"func", reflect.ValueOf((func())(nil)), true},
		{"string", reflect.ValueOf("hello"), false},
		{"int", reflect.ValueOf(42), false},
		{"bool", reflect.ValueOf(true), false},
		{"struct", reflect.ValueOf(testInput{}), false},
		{"float", reflect.ValueOf(3.14), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canBeNil(tt.val)
			if got != tt.want {
				t.Errorf("canBeNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- execute (JSON-RPC path) additional tests ---

func TestExecute_MethodNotFound(t *testing.T) {
	b := New(WithCSRF(false))
	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))

	result := b.execute(ctx, "nonexistent", nil)

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if result.Error.Code != ErrCodeMethodNotFound {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeMethodNotFound)
	}
}

func TestExecute_InvalidJSON(t *testing.T) {
	b := New(WithCSRF(false))
	_ = b.Register("test.fn", func(ctx Context, p testInput) (testOutput, error) {
		return testOutput{}, nil
	})

	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.execute(ctx, "test.fn", json.RawMessage(`{invalid`))

	if result.Error == nil {
		t.Fatal("expected parse error, got nil")
	}
	if result.Error.Code != ErrCodeInvalidParams {
		t.Errorf("error code = %d, want %d", result.Error.Code, ErrCodeInvalidParams)
	}
}

func TestExecute_NoInput_SkipsParseValidate(t *testing.T) {
	b := New(WithCSRF(false))
	_ = b.Register("test.noInput", func(ctx Context) (testOutput, error) {
		return testOutput{Result: "works"}, nil
	})

	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.execute(ctx, "test.noInput", nil)

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	out := result.Result.(testOutput)
	if out.Result != "works" {
		t.Errorf("result = %q, want %q", out.Result, "works")
	}
}

func TestExecute_LaxValidation_AllowsEmpty(t *testing.T) {
	b := New(WithCSRF(false))
	_ = b.Register("test.laxFn", func(ctx Context, p testInput) (testOutput, error) {
		return testOutput{Result: "ok"}, nil
	}, WithLaxValidation())

	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	result := b.execute(ctx, "test.laxFn", json.RawMessage(`{}`))

	if result.Error != nil {
		t.Fatalf("unexpected error (lax mode should pass): %v", result.Error)
	}
}

// --- HookManager utility tests ---

func TestHookManager_Count(t *testing.T) {
	hm := NewHookManager()

	if got := hm.Count(BeforeCall); got != 0 {
		t.Errorf("initial Count = %d, want 0", got)
	}

	hm.Register(BeforeCall, func(ctx Context, data HookData) {})
	hm.Register(BeforeCall, func(ctx Context, data HookData) {})

	if got := hm.Count(BeforeCall); got != 2 {
		t.Errorf("after 2 registers, Count = %d, want 2", got)
	}

	if got := hm.Count(AfterCall); got != 0 {
		t.Errorf("AfterCall Count = %d, want 0 (different type)", got)
	}
}

func TestHookManager_Clear(t *testing.T) {
	hm := NewHookManager()

	hm.Register(BeforeCall, func(ctx Context, data HookData) {})
	hm.Register(AfterCall, func(ctx Context, data HookData) {})

	hm.Clear(BeforeCall)

	if got := hm.Count(BeforeCall); got != 0 {
		t.Errorf("after Clear, BeforeCall Count = %d, want 0", got)
	}
	if got := hm.Count(AfterCall); got != 1 {
		t.Errorf("AfterCall Count = %d, want 1 (not cleared)", got)
	}
}

func TestHookManager_ClearAll(t *testing.T) {
	hm := NewHookManager()

	hm.Register(BeforeCall, func(ctx Context, data HookData) {})
	hm.Register(AfterCall, func(ctx Context, data HookData) {})
	hm.Register(OnError, func(ctx Context, data HookData) {})

	hm.ClearAll()

	if got := hm.Count(BeforeCall); got != 0 {
		t.Errorf("BeforeCall Count = %d, want 0", got)
	}
	if got := hm.Count(AfterCall); got != 0 {
		t.Errorf("AfterCall Count = %d, want 0", got)
	}
	if got := hm.Count(OnError); got != 0 {
		t.Errorf("OnError Count = %d, want 0", got)
	}
}

func TestHookManager_TriggerPanicRecovery(t *testing.T) {
	hm := NewHookManager()

	var safeCalled int32
	hm.Register(BeforeCall, func(ctx Context, data HookData) {
		panic("test panic in hook")
	})
	hm.Register(BeforeCall, func(ctx Context, data HookData) {
		atomic.AddInt32(&safeCalled, 1)
	})

	ctx := NewContext(httptest.NewRequest(http.MethodGet, "/", nil))
	hm.Trigger(BeforeCall, ctx, HookData{FunctionName: "test"})

	// Give goroutines time to complete
	time.Sleep(50 * time.Millisecond)

	// The panicking hook should be recovered; the other should still run
	if atomic.LoadInt32(&safeCalled) != 1 {
		t.Errorf("safe hook called %d times, want 1", safeCalled)
	}
}
