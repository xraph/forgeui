package bridge

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurity_CheckAuth(t *testing.T) {
	config := DefaultConfig()
	security := NewSecurity(config)

	// Function that doesn't require auth
	fn := &Function{RequireAuth: false}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(req)

	err := security.CheckAuth(ctx, fn)
	if err != nil {
		t.Errorf("CheckAuth() should pass for function without auth requirement")
	}

	// Function that requires auth
	fn = &Function{RequireAuth: true}

	err = security.CheckAuth(ctx, fn)
	if err == nil {
		t.Error("CheckAuth() should fail when user is nil")
	}

	// Add user to context
	user := &SimpleUser{UserID: "123"}
	ctx = WithUser(ctx, user)

	err = security.CheckAuth(ctx, fn)
	if err != nil {
		t.Errorf("CheckAuth() should pass with user: %v", err)
	}
}

func TestSecurity_CheckAuth_Roles(t *testing.T) {
	config := DefaultConfig()
	security := NewSecurity(config)

	fn := &Function{
		RequireAuth:  true,
		RequireRoles: []string{"admin"},
	}

	// User without required role
	user := &SimpleUser{
		UserID:    "123",
		UserRoles: []string{"user"},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := NewContext(req)
	ctx = WithUser(ctx, user)

	err := security.CheckAuth(ctx, fn)
	if err == nil {
		t.Error("CheckAuth() should fail without required role")
	}

	// User with required role
	user.UserRoles = []string{"admin", "user"}
	ctx = WithUser(ctx, user)

	err = security.CheckAuth(ctx, fn)
	if err != nil {
		t.Errorf("CheckAuth() should pass with required role: %v", err)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name:    "X-Forwarded-For",
			headers: map[string]string{"X-Forwarded-For": "192.168.1.1, 10.0.0.1"},
			want:    "192.168.1.1",
		},
		{
			name:    "X-Real-IP",
			headers: map[string]string{"X-Real-IP": "192.168.1.1"},
			want:    "192.168.1.1",
		},
		{
			name:       "RemoteAddr",
			remoteAddr: "192.168.1.1:12345",
			want:       "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr
			}

			got := GetClientIP(req)
			if got != tt.want {
				t.Errorf("GetClientIP() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestGenerateCSRFToken(t *testing.T) {
	token1, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf("GenerateCSRFToken() error = %v", err)
	}

	if len(token1) == 0 {
		t.Error("GenerateCSRFToken() returned empty string")
	}

	token2, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf("GenerateCSRFToken() error = %v", err)
	}

	// Tokens should be different
	if token1 == token2 {
		t.Error("GenerateCSRFToken() should generate unique tokens")
	}
}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid input", "normal text", false},
		{"null byte", "text\x00with null", true},
		{"empty", "", false},
		{"very long", string(make([]byte, 2_000_000)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
