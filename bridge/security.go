package bridge

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"slices"
	"strings"
)

// Security handles authentication, authorization, and CSRF
type Security struct {
	rateLimiter *RateLimiter
	csrfEnabled bool
	csrfHeader  string
	csrfCookie  string
}

// NewSecurity creates a new security instance
func NewSecurity(config *Config) *Security {
	return &Security{
		rateLimiter: NewRateLimiter(config.DefaultRateLimit, config.DefaultRateLimit*2),
		csrfEnabled: config.EnableCSRF,
		csrfHeader:  config.CSRFTokenHeader,
		csrfCookie:  config.CSRFCookieName,
	}
}

// CheckAuth verifies authentication
func (s *Security) CheckAuth(ctx Context, fn *Function) error {
	if !fn.RequireAuth {
		return nil
	}

	user := ctx.User()
	if user == nil {
		return ErrUnauthorized
	}

	// Check required roles
	if len(fn.RequireRoles) > 0 {
		hasRole := slices.ContainsFunc(fn.RequireRoles, user.HasRole)

		if !hasRole {
			return ErrForbidden
		}
	}

	return nil
}

// CheckRateLimit verifies rate limiting
func (s *Security) CheckRateLimit(key string, fn *Function) error {
	// Use function-specific rate limit if set, otherwise use default
	rateLimit := fn.RateLimit
	if rateLimit == 0 {
		// No rate limit for this function
		return nil
	}

	if !s.rateLimiter.Allow(key) {
		return ErrRateLimit
	}

	return nil
}

// CheckCSRF verifies CSRF token
func (s *Security) CheckCSRF(r *http.Request) error {
	if !s.csrfEnabled {
		return nil
	}

	// Skip CSRF for GET requests
	if r.Method == http.MethodGet {
		return nil
	}

	// Get token from header
	headerToken := r.Header.Get(s.csrfHeader)

	// Get token from cookie
	cookie, err := r.Cookie(s.csrfCookie)
	if err != nil {
		return NewError(ErrCodeBadRequest, "CSRF token missing")
	}

	cookieToken := cookie.Value

	// Compare tokens using constant-time comparison
	if !secureCompare(headerToken, cookieToken) {
		return NewError(ErrCodeBadRequest, "CSRF token mismatch")
	}

	return nil
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// SetCSRFCookie sets the CSRF token cookie
func SetCSRFCookie(w http.ResponseWriter, token string, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		MaxAge:   86400, // 24 hours
	})
}

// secureCompare performs constant-time string comparison
func secureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

// getRateLimitKey generates a rate limit key based on IP and user
func getRateLimitKey(ctx Context) string {
	user := ctx.User()
	if user != nil {
		return "user:" + user.ID()
	}

	req := ctx.Request()

	return "ip:" + GetClientIP(req)
}

// ValidateInput performs basic input sanitization
func ValidateInput(input string) error {
	// Check for null bytes
	if strings.Contains(input, "\x00") {
		return NewError(ErrCodeBadRequest, "Input contains null bytes")
	}

	// Check maximum length
	if len(input) > 1_000_000 { // 1MB
		return NewError(ErrCodeBadRequest, "Input too large")
	}

	return nil
}
