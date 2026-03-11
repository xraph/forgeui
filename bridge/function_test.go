package bridge

import (
	"reflect"
	"testing"
	"time"

	"github.com/a-h/templ"
)

type testInput struct {
	Name  string `json:"name"`
	Value int    `json:"value,omitempty"`
}

type testOutput struct {
	Result string `json:"result"`
}

func validHandler(ctx Context, input testInput) (testOutput, error) {
	return testOutput{Result: "ok"}, nil
}

func TestValidateFunction(t *testing.T) {
	tests := []struct {
		name    string
		fn      any
		wantSig SignatureType
		wantErr bool
	}{
		{
			name:    "valid function (input+output)",
			fn:      validHandler,
			wantSig: SigInputOutput,
			wantErr: false,
		},
		{
			name:    "not a function",
			fn:      "not a function",
			wantErr: true,
		},
		{
			name: "output only (no input)",
			fn: func(ctx Context) (testOutput, error) {
				return testOutput{}, nil
			},
			wantSig: SigOutput,
			wantErr: false,
		},
		{
			name: "input only (no output)",
			fn: func(ctx Context, input testInput) error {
				return nil
			},
			wantSig: SigInputOnly,
			wantErr: false,
		},
		{
			name: "void (no input, no output)",
			fn: func(ctx Context) error {
				return nil
			},
			wantSig: SigVoid,
			wantErr: false,
		},
		{
			name: "returns templ.Component",
			fn: func(ctx Context) (templ.Component, error) {
				return nil, nil
			},
			wantSig: SigOutput,
			wantErr: false,
		},
		{
			name: "returns templ.Component with input",
			fn: func(ctx Context, input testInput) (templ.Component, error) {
				return nil, nil
			},
			wantSig: SigInputOutput,
			wantErr: false,
		},
		{
			name: "second output not error",
			fn: func(ctx Context, input testInput) (testOutput, string) {
				return testOutput{}, ""
			},
			wantErr: true,
		},
		{
			name: "too many inputs",
			fn: func(ctx Context, a testInput, b testInput) (testOutput, error) {
				return testOutput{}, nil
			},
			wantErr: true,
		},
		{
			name: "too many outputs",
			fn: func(ctx Context) (testOutput, string, error) {
				return testOutput{}, "", nil
			},
			wantErr: true,
		},
		{
			name: "first param not Context",
			fn: func(s string) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sigType, err := validateFunction(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && sigType != tt.wantSig {
				t.Errorf("validateFunction() sigType = %v, want %v", sigType, tt.wantSig)
			}
		})
	}
}

func TestAnalyzeFunction(t *testing.T) {
	fn, err := analyzeFunction(validHandler)
	if err != nil {
		t.Fatalf("analyzeFunction() error = %v", err)
	}

	if fn.InputType.Name() != "testInput" {
		t.Errorf("InputType = %v, want testInput", fn.InputType.Name())
	}

	if fn.OutputType.Name() != "testOutput" {
		t.Errorf("OutputType = %v, want testOutput", fn.OutputType.Name())
	}

	if fn.SignatureType != SigInputOutput {
		t.Errorf("SignatureType = %v, want SigInputOutput", fn.SignatureType)
	}

	if !fn.HasInput {
		t.Error("HasInput = false, want true")
	}

	if !fn.HasOutput {
		t.Error("HasOutput = false, want true")
	}

	if fn.ReturnsHTML {
		t.Error("ReturnsHTML = true, want false")
	}

	// Timeout should be 0 by default - the bridge config timeout is used as fallback during execution
	if fn.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0 (unset)", fn.Timeout)
	}
}

func TestAnalyzeFunction_OutputOnly(t *testing.T) {
	fn, err := analyzeFunction(func(ctx Context) (testOutput, error) {
		return testOutput{}, nil
	})
	if err != nil {
		t.Fatalf("analyzeFunction() error = %v", err)
	}

	if fn.SignatureType != SigOutput {
		t.Errorf("SignatureType = %v, want SigOutput", fn.SignatureType)
	}
	if fn.HasInput {
		t.Error("HasInput = true, want false")
	}
	if !fn.HasOutput {
		t.Error("HasOutput = false, want true")
	}
	if fn.InputType != nil {
		t.Errorf("InputType = %v, want nil", fn.InputType)
	}
}

func TestAnalyzeFunction_Void(t *testing.T) {
	fn, err := analyzeFunction(func(ctx Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("analyzeFunction() error = %v", err)
	}

	if fn.SignatureType != SigVoid {
		t.Errorf("SignatureType = %v, want SigVoid", fn.SignatureType)
	}
	if fn.HasInput {
		t.Error("HasInput = true, want false")
	}
	if fn.HasOutput {
		t.Error("HasOutput = true, want false")
	}
}

func TestAnalyzeFunction_ReturnsHTML(t *testing.T) {
	fn, err := analyzeFunction(func(ctx Context) (templ.Component, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("analyzeFunction() error = %v", err)
	}

	if !fn.ReturnsHTML {
		t.Error("ReturnsHTML = false, want true")
	}
	if fn.SignatureType != SigOutput {
		t.Errorf("SignatureType = %v, want SigOutput", fn.SignatureType)
	}
}

func TestFunctionOptions(t *testing.T) {
	fn := &Function{}

	RequireAuth()(fn)

	if !fn.RequireAuth {
		t.Error("RequireAuth() did not set RequireAuth to true")
	}

	RequireRoles("admin", "editor")(fn)

	if len(fn.RequireRoles) != 2 {
		t.Errorf("RequireRoles() set %d roles, want 2", len(fn.RequireRoles))
	}

	WithFunctionTimeout(10 * time.Second)(fn)

	if fn.Timeout != 10*time.Second {
		t.Errorf("WithFunctionTimeout() set %v, want 10s", fn.Timeout)
	}

	WithRateLimit(100)(fn)

	if fn.RateLimit != 100 {
		t.Errorf("WithRateLimit() set %d, want 100", fn.RateLimit)
	}

	WithFunctionCache(5 * time.Minute)(fn)

	if !fn.Cacheable {
		t.Error("WithFunctionCache() did not set Cacheable to true")
	}

	if fn.CacheTTL != 5*time.Minute {
		t.Errorf("WithFunctionCache() set TTL %v, want 5m", fn.CacheTTL)
	}

	WithDescription("test description")(fn)

	if fn.Description != "test description" {
		t.Errorf("WithDescription() set %s, want 'test description'", fn.Description)
	}
}

func TestNewFunctionOptions(t *testing.T) {
	fn := &Function{}

	WithHTTPMethod("GET", "POST")(fn)
	if len(fn.AllowedMethods) != 2 || fn.AllowedMethods[0] != "GET" {
		t.Errorf("WithHTTPMethod() = %v, want [GET POST]", fn.AllowedMethods)
	}

	WithHTMXTrigger("itemCreated", "listUpdated")(fn)
	if len(fn.HTMXTriggers) != 2 {
		t.Errorf("WithHTMXTrigger() set %d triggers, want 2", len(fn.HTMXTriggers))
	}

	WithHTMXRedirect("/dashboard")(fn)
	if fn.HTMXRedirect != "/dashboard" {
		t.Errorf("WithHTMXRedirect() = %s, want /dashboard", fn.HTMXRedirect)
	}

	WithHTMXReswap("none")(fn)
	if fn.HTMXReswap != "none" {
		t.Errorf("WithHTMXReswap() = %s, want none", fn.HTMXReswap)
	}

	WithHTMXRetarget("#content")(fn)
	if fn.HTMXRetarget != "#content" {
		t.Errorf("WithHTMXRetarget() = %s, want #content", fn.HTMXRetarget)
	}

	WithLaxValidation()(fn)
	if !fn.LaxValidation {
		t.Error("WithLaxValidation() did not set LaxValidation to true")
	}
}

func TestWithRenderer(t *testing.T) {
	fn := &Function{}

	WithRenderer(func(data testOutput) templ.Component {
		return nil // just testing the option sets the renderer
	})(fn)

	if fn.Renderer == nil {
		t.Error("WithRenderer() did not set Renderer")
	}
}

func TestFunction_GetTypeInfo(t *testing.T) {
	fn, err := analyzeFunction(validHandler)
	if err != nil {
		t.Fatalf("analyzeFunction() error = %v", err)
	}

	fn.Name = "testFunc"
	info := fn.GetTypeInfo()

	if info.Name != "testFunc" {
		t.Errorf("TypeInfo.Name = %s, want testFunc", info.Name)
	}

	if len(info.Fields) != 2 {
		t.Errorf("len(Fields) = %d, want 2", len(info.Fields))
	}

	// Check Name field
	if info.Fields[0].Name != "Name" {
		t.Errorf("Fields[0].Name = %s, want Name", info.Fields[0].Name)
	}

	if info.Fields[0].JSONName != "name" {
		t.Errorf("Fields[0].JSONName = %s, want name", info.Fields[0].JSONName)
	}

	if !info.Fields[0].Required {
		t.Error("Fields[0].Required = false, want true")
	}

	// Check Value field (has omitempty)
	if info.Fields[1].Name != "Value" {
		t.Errorf("Fields[1].Name = %s, want Value", info.Fields[1].Name)
	}

	if info.Fields[1].Required {
		t.Error("Fields[1].Required = true, want false (has omitempty)")
	}
}

func TestFunction_GetTypeInfo_NoInput(t *testing.T) {
	fn, err := analyzeFunction(func(ctx Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("analyzeFunction() error = %v", err)
	}

	fn.Name = "voidFunc"
	info := fn.GetTypeInfo()

	if info.InputType != "" {
		t.Errorf("TypeInfo.InputType = %s, want empty", info.InputType)
	}
	if info.OutputType != "" {
		t.Errorf("TypeInfo.OutputType = %s, want empty", info.OutputType)
	}
	if len(info.Fields) != 0 {
		t.Errorf("len(Fields) = %d, want 0", len(info.Fields))
	}
}

func TestExtractFields(t *testing.T) {
	type localTestStruct struct {
		PublicField  string `json:"public"`
		privateField string
		OmitEmpty    int    `json:"omit,omitempty"`
		IgnoredField string `json:"-"`
		NoTagField   bool
	}

	// Suppress unused variable warning for privateField
	_ = localTestStruct{privateField: ""}

	fields := extractFields(reflect.TypeFor[localTestStruct]())

	// Should have 3 fields (PublicField, OmitEmpty, NoTagField)
	// privateField is not exported
	// IgnoredField has json:"-"
	if len(fields) != 3 {
		t.Errorf("len(fields) = %d, want 3", len(fields))
	}
}
