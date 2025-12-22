package bridge

import (
	"context"
	"net/http"
	"slices"
)

// Context provides access to request-scoped data
type Context interface {
	// Context returns the underlying context.Context
	Context() context.Context

	// Request returns the HTTP request
	Request() *http.Request

	// Value retrieves a value from the context
	Value(key any) any

	// Session returns session data (implementation-specific)
	Session() Session

	// User returns the authenticated user (if any)
	User() User

	// SetValue stores a value in the context
	SetValue(key, val any)
}

// Session represents user session data
type Session interface {
	// Get retrieves a value from the session
	Get(key string) (any, bool)

	// Set stores a value in the session
	Set(key string, value any)

	// ID returns the session ID
	ID() string

	// Delete removes a value from the session
	Delete(key string)

	// Clear removes all values from the session
	Clear()
}

// User represents an authenticated user
type User interface {
	// ID returns the user's unique identifier
	ID() string

	// Email returns the user's email address
	Email() string

	// Name returns the user's display name
	Name() string

	// Roles returns the user's roles/permissions
	Roles() []string

	// HasRole checks if the user has a specific role
	HasRole(role string) bool

	// Data returns additional user data
	Data() map[string]any
}

// bridgeContext is the default implementation of Context
type bridgeContext struct {
	ctx     context.Context
	req     *http.Request
	session Session
	user    User
	values  map[any]any
}

// NewContext creates a new bridge context
func NewContext(r *http.Request) Context {
	return &bridgeContext{
		ctx:    r.Context(),
		req:    r,
		values: make(map[any]any),
	}
}

// Context returns the underlying context.Context
func (c *bridgeContext) Context() context.Context {
	return c.ctx
}

// Request returns the HTTP request
func (c *bridgeContext) Request() *http.Request {
	return c.req
}

// Value retrieves a value from the context
func (c *bridgeContext) Value(key any) any {
	if val, ok := c.values[key]; ok {
		return val
	}

	return c.ctx.Value(key)
}

// Session returns session data
func (c *bridgeContext) Session() Session {
	return c.session
}

// User returns the authenticated user
func (c *bridgeContext) User() User {
	return c.user
}

// SetValue stores a value in the context
func (c *bridgeContext) SetValue(key, val any) {
	c.values[key] = val
}

// WithSession returns a new context with the given session
func WithSession(ctx Context, session Session) Context {
	if bc, ok := ctx.(*bridgeContext); ok {
		bc.session = session
	}

	return ctx
}

// WithUser returns a new context with the given user
func WithUser(ctx Context, user User) Context {
	if bc, ok := ctx.(*bridgeContext); ok {
		bc.user = user
	}

	return ctx
}

// SimpleUser is a basic implementation of User
type SimpleUser struct {
	UserID    string         `json:"id"`
	UserEmail string         `json:"email"`
	UserName  string         `json:"name"`
	UserRoles []string       `json:"roles"`
	UserData  map[string]any `json:"data"`
}

// ID returns the user's unique identifier
func (u *SimpleUser) ID() string {
	return u.UserID
}

// Email returns the user's email address
func (u *SimpleUser) Email() string {
	return u.UserEmail
}

// Name returns the user's display name
func (u *SimpleUser) Name() string {
	return u.UserName
}

// Roles returns the user's roles/permissions
func (u *SimpleUser) Roles() []string {
	return u.UserRoles
}

// HasRole checks if the user has a specific role
func (u *SimpleUser) HasRole(role string) bool {
	return slices.Contains(u.UserRoles, role)
}

// Data returns additional user data
func (u *SimpleUser) Data() map[string]any {
	return u.UserData
}

// SimpleSession is a basic implementation of Session
type SimpleSession struct {
	SessionID string
	data      map[string]any
}

// NewSimpleSession creates a new simple session
func NewSimpleSession(id string) *SimpleSession {
	return &SimpleSession{
		SessionID: id,
		data:      make(map[string]any),
	}
}

// Get retrieves a value from the session
func (s *SimpleSession) Get(key string) (any, bool) {
	val, ok := s.data[key]
	return val, ok
}

// Set stores a value in the session
func (s *SimpleSession) Set(key string, value any) {
	s.data[key] = value
}

// ID returns the session ID
func (s *SimpleSession) ID() string {
	return s.SessionID
}

// Delete removes a value from the session
func (s *SimpleSession) Delete(key string) {
	delete(s.data, key)
}

// Clear removes all values from the session
func (s *SimpleSession) Clear() {
	s.data = make(map[string]any)
}
