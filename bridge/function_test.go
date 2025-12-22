package bridge

import (
	"reflect"
	"testing"
	"time"
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
		wantErr bool
	}{
		{
			name:    "valid function",
			fn:      validHandler,
			wantErr: false,
		},
		{
			name:    "not a function",
			fn:      "not a function",
			wantErr: true,
		},
		{
			name: "wrong number of inputs",
			fn: func(ctx Context) (testOutput, error) {
				return testOutput{}, nil
			},
			wantErr: true,
		},
		{
			name: "wrong number of outputs",
			fn: func(ctx Context, input testInput) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "second output not error",
			fn: func(ctx Context, input testInput) (testOutput, string) {
				return testOutput{}, ""
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFunction(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFunction() error = %v, wantErr %v", err, tt.wantErr)
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

	// Timeout should be 0 by default - the bridge config timeout is used as fallback during execution
	if fn.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0 (unset)", fn.Timeout)
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
