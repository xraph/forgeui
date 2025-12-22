package bridge

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := NewContext(req)

	if ctx.Request() != req {
		t.Error("Request() does not return the original request")
	}

	if ctx.Context() != req.Context() {
		t.Error("Context() does not return the request's context")
	}
}

func TestContext_Value(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := NewContext(req)

	key := "test-key"
	value := "test-value"

	ctx.SetValue(key, value)

	got := ctx.Value(key)
	if got != value {
		t.Errorf("Value() = %v, want %v", got, value)
	}
}

func TestWithSession(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := NewContext(req)

	session := NewSimpleSession("session-123")
	ctx = WithSession(ctx, session)

	if ctx.Session() == nil {
		t.Error("Session() returned nil after WithSession")
	}

	if ctx.Session().ID() != "session-123" {
		t.Errorf("Session ID = %s, want session-123", ctx.Session().ID())
	}
}

func TestWithUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := NewContext(req)

	user := &SimpleUser{
		UserID:    "user-123",
		UserEmail: "test@example.com",
		UserName:  "Test User",
		UserRoles: []string{"admin", "user"},
	}

	ctx = WithUser(ctx, user)

	if ctx.User() == nil {
		t.Error("User() returned nil after WithUser")
	}

	if ctx.User().ID() != "user-123" {
		t.Errorf("User ID = %s, want user-123", ctx.User().ID())
	}
}

func TestSimpleUser(t *testing.T) {
	user := &SimpleUser{
		UserID:    "123",
		UserEmail: "user@test.com",
		UserName:  "John Doe",
		UserRoles: []string{"admin", "editor"},
		UserData:  map[string]any{"department": "engineering"},
	}

	if user.ID() != "123" {
		t.Errorf("ID() = %s, want 123", user.ID())
	}

	if user.Email() != "user@test.com" {
		t.Errorf("Email() = %s, want user@test.com", user.Email())
	}

	if user.Name() != "John Doe" {
		t.Errorf("Name() = %s, want John Doe", user.Name())
	}

	if !user.HasRole("admin") {
		t.Error("HasRole(admin) = false, want true")
	}

	if user.HasRole("superuser") {
		t.Error("HasRole(superuser) = true, want false")
	}

	if len(user.Roles()) != 2 {
		t.Errorf("len(Roles()) = %d, want 2", len(user.Roles()))
	}
}

func TestSimpleSession(t *testing.T) {
	session := NewSimpleSession("sess-456")

	if session.ID() != "sess-456" {
		t.Errorf("ID() = %s, want sess-456", session.ID())
	}

	// Test Set and Get
	session.Set("key1", "value1")

	val, ok := session.Get("key1")
	if !ok {
		t.Error("Get(key1) returned false, want true")
	}

	if val != "value1" {
		t.Errorf("Get(key1) = %v, want value1", val)
	}

	// Test non-existent key
	_, ok = session.Get("nonexistent")
	if ok {
		t.Error("Get(nonexistent) returned true, want false")
	}

	// Test Delete
	session.Delete("key1")

	_, ok = session.Get("key1")
	if ok {
		t.Error("Get(key1) after Delete returned true, want false")
	}

	// Test Clear
	session.Set("key2", "value2")
	session.Set("key3", "value3")
	session.Clear()

	_, ok = session.Get("key2")
	if ok {
		t.Error("Get(key2) after Clear returned true, want false")
	}
}
