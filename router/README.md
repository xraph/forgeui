# ForgeUI Router

Production-ready HTTP routing system for ForgeUI applications.

## Features

- **Pattern Matching**: Static routes, path parameters (`:id`), and wildcards (`*path`)
- **HTTP Methods**: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD
- **Middleware**: Global and route-specific middleware with chainable execution
- **Named Routes**: URL generation from route names
- **PageContext**: Rich request context with params, query, headers, cookies
- **Type Safety**: Full Go type safety with templ integration
- **Priority Routing**: Automatic route prioritization (static > params > wildcards)

## Quick Start

```go
package main

import (
    "github.com/a-h/templ"
    "github.com/xraph/forgeui"
    "github.com/xraph/forgeui/router"
)

func main() {
    app := forgeui.New()

    // Register routes
    app.Get("/", HomePage)
    app.Get("/users/:id", UserProfile)
    app.Post("/users", CreateUser)

    // Start server
    http.ListenAndServe(":8080", app.Handler())
}

func HomePage(ctx *router.PageContext) (templ.Component, error) {
    return HomePageView(), nil
}

func UserProfile(ctx *router.PageContext) (templ.Component, error) {
    id := ctx.Param("id")
    return UserProfileView(id), nil
}

func CreateUser(ctx *router.PageContext) (templ.Component, error) {
    return UserCreatedView(), nil
}
```

## Route Patterns

### Static Routes

```go
app.Get("/about", AboutPage)
app.Get("/contact", ContactPage)
```

### Path Parameters

```go
// Single parameter
app.Get("/users/:id", UserProfile)

// Multiple parameters
app.Get("/users/:userId/posts/:postId", PostDetail)
```

### Wildcards

```go
// Matches /files/docs/readme.md, /files/images/logo.png, etc.
app.Get("/files/*filepath", ServeFile)
```

## HTTP Methods

```go
app.Get("/users", ListUsers)
app.Post("/users", CreateUser)
app.Put("/users/:id", UpdateUser)
app.Patch("/users/:id", PatchUser)
app.Delete("/users/:id", DeleteUser)
app.Options("/users", OptionsUsers)
app.Head("/users", HeadUsers)

// Multiple methods on same route
app.Match([]string{"GET", "POST"}, "/api", APIHandler)
```

## PageContext

The `PageContext` provides access to request data and utilities:

### Path Parameters

```go
func UserProfile(ctx *router.PageContext) (templ.Component, error) {
    id := ctx.Param("id")                    // String
    userID, err := ctx.ParamInt("id")        // Integer
    userID64, err := ctx.ParamInt64("id")    // Int64

    return UserProfileView(id), nil
}
```

### Query Parameters

```go
func SearchUsers(ctx *router.PageContext) (templ.Component, error) {
    query := ctx.Query("q")                       // String
    page, err := ctx.QueryInt("page")             // Integer
    limit := ctx.QueryDefault("limit", "10")      // With default
    active := ctx.QueryBool("active")             // Boolean

    return SearchResultsView(query, page), nil
}
```

### Headers & Cookies

```go
func Protected(ctx *router.PageContext) (templ.Component, error) {
    auth := ctx.Header("Authorization")

    cookie, err := ctx.Cookie("session")
    if err != nil {
        // Handle missing cookie
    }

    // Set response headers
    ctx.SetHeader("X-Custom", "value")

    // Set cookies
    ctx.SetCookie(&http.Cookie{
        Name:  "session",
        Value: "abc123",
    })

    return ProtectedView(), nil
}
```

### Context Values

Store and retrieve values for the request lifecycle:

```go
// In middleware
func AuthMiddleware(next router.PageHandler) router.PageHandler {
    return func(ctx *router.PageContext) (templ.Component, error) {
        ctx.Set("user_id", 42)
        ctx.Set("username", "john")
        return next(ctx)
    }
}

// In handler
func Profile(ctx *router.PageContext) (templ.Component, error) {
    userID := ctx.GetInt("user_id")
    username := ctx.GetString("username")
    
    // ...
}
```

### Request Information

```go
method := ctx.Method()        // HTTP method
path := ctx.Path()            // URL path
host := ctx.Host()            // Host header
isSecure := ctx.IsSecure()    // HTTPS check
clientIP := ctx.ClientIP()    // Client IP (respects X-Forwarded-For)
```

## Middleware

### Global Middleware

Applied to all routes:

```go
app.Use(router.Logger())
app.Use(router.Recovery())
app.Use(router.RequestID())
```

### Route-Specific Middleware

Applied to individual routes:

```go
route := app.Get("/admin", AdminDashboard)
route.WithMiddleware(AuthMiddleware)
```

### Built-in Middleware

```go
// Logging
app.Use(router.Logger())

// Panic recovery
app.Use(router.Recovery())

// CORS
app.Use(router.CORS("*"))

// Request ID
app.Use(router.RequestID())

// Basic authentication
app.Use(router.BasicAuth("admin", "secret"))

// Method restriction
app.Use(router.RequireMethod("GET", "POST"))

// Custom headers
app.Use(router.SetHeader("X-Frame-Options", "DENY"))

// Request timeout
app.Use(router.Timeout(30 * time.Second))
```

### Custom Middleware

```go
func AuthMiddleware(next router.PageHandler) router.PageHandler {
    return func(ctx *router.PageContext) (templ.Component, error) {
        // Check authentication
        token := ctx.Header("Authorization")
        if token == "" {
            ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
            return UnauthorizedView(), nil
        }
        
        // Continue to next handler
        return next(ctx)
    }
}

app.Use(AuthMiddleware)
```

### Chaining Middleware

```go
middleware := router.Chain(
    router.Logger(),
    router.Recovery(),
    router.RequestID(),
)

app.Use(middleware)
```

## Named Routes & URL Generation

```go
// Register named route
route := app.Get("/users/:id/posts/:postId", PostDetail)
app.Router().Name("user.post", route)

// Generate URL
url := app.Router().URL("user.post", 123, 456)
// Returns: "/users/123/posts/456"

// In templ files, generate URLs in handlers and pass to templates
```

## Router Options

```go
router := router.New(
    router.WithBasePath("/api/v1"),
    router.WithNotFound(Custom404Handler),
    router.WithErrorHandler(CustomErrorHandler),
)
```

### Custom 404 Handler

```go
func Custom404(ctx *router.PageContext) (templ.Component, error) {
    ctx.ResponseWriter.WriteHeader(404)
    return NotFoundView(), nil
}

router := router.New(router.WithNotFound(Custom404))
```

### Custom Error Handler

```go
func CustomError(ctx *router.PageContext, err error) templ.Component {
    ctx.ResponseWriter.WriteHeader(500)
    return ErrorView(err)
}

router := router.New(router.WithErrorHandler(CustomError))
```

## Testing

The router is designed for easy testing with `httptest`:

```go
func TestUserProfile(t *testing.T) {
    r := router.New()
    r.Get("/users/:id", UserProfile)
    
    req := httptest.NewRequest("GET", "/users/123", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Best Practices

1. **Use middleware for cross-cutting concerns**: Authentication, logging, CORS
2. **Keep handlers focused**: One responsibility per handler
3. **Return errors**: Let the error handler deal with them
4. **Use named routes**: Makes refactoring easier
5. **Test with httptest**: Fast, reliable, no network needed
6. **Set appropriate status codes**: Use `ctx.ResponseWriter.WriteHeader()`
7. **Validate input**: Check params and query values before use

## Performance

- Routes are sorted by priority at registration time
- Pattern matching uses compiled regex
- Middleware chains are built once per route
- Zero allocations for static routes
- Thread-safe route registration and lookup

## Integration with ForgeUI

The router is fully integrated with ForgeUI's App struct:

```go
app := forgeui.New()

// These are equivalent:
app.Get("/users", ListUsers)
app.Router().Get("/users", ListUsers)

// Access router directly for advanced features:
app.Router().Use(middleware...)
app.Router().Name("route.name", route)
```

## Examples

See the [example directory](../example/) for complete working examples.

