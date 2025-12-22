package bridge

import (
	"encoding/json"
	"testing"
)

func TestError_Error(t *testing.T) {
	err := NewError(ErrCodeInternal, "test error")
	if err.Error() != "test error" {
		t.Errorf("Error() = %s, want test error", err.Error())
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
		data    []any
		want    *Error
	}{
		{
			name:    "without data",
			code:    ErrCodeInternal,
			message: "internal error",
			data:    nil,
			want:    &Error{Code: ErrCodeInternal, Message: "internal error"},
		},
		{
			name:    "with data",
			code:    ErrCodeInvalidParams,
			message: "invalid params",
			data:    []any{"field: email"},
			want:    &Error{Code: ErrCodeInvalidParams, Message: "invalid params", Data: "field: email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.code, tt.message, tt.data...)
			if got.Code != tt.want.Code || got.Message != tt.want.Message {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_JSON(t *testing.T) {
	req := Request{
		JSONRPC: "2.0",
		ID:      "1",
		Method:  "testMethod",
		Params:  json.RawMessage(`{"key":"value"}`),
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var decoded Request
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if decoded.Method != req.Method {
		t.Errorf("Method = %s, want %s", decoded.Method, req.Method)
	}
}

func TestResponse_JSON(t *testing.T) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      "1",
		Result:  map[string]string{"status": "ok"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if decoded.JSONRPC != resp.JSONRPC {
		t.Errorf("JSONRPC = %s, want %s", decoded.JSONRPC, resp.JSONRPC)
	}
}

func TestBatchRequest(t *testing.T) {
	batch := BatchRequest{
		{JSONRPC: "2.0", ID: "1", Method: "method1"},
		{JSONRPC: "2.0", ID: "2", Method: "method2"},
	}

	data, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("Failed to marshal batch: %v", err)
	}

	var decoded BatchRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal batch: %v", err)
	}

	if len(decoded) != 2 {
		t.Errorf("len(batch) = %d, want 2", len(decoded))
	}
}
