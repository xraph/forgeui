# ForgeUI Example

This example demonstrates all ForgeUI components in action.

## Quick Start

1. **Run the example:**

```bash
cd example
go run main.go
```

2. **Open your browser:**

Visit [http://localhost:8080](http://localhost:8080)

## What's Included

The example showcases:

- ✅ **Buttons** - All variants and sizes
- ✅ **Cards** - Compound components with header, content, footer
- ✅ **Badges** - Status indicators
- ✅ **Avatars** - User avatars with fallbacks
- ✅ **Alerts** - Message displays
- ✅ **Forms** - Complete form with inputs, checkboxes, radio buttons
- ✅ **Loading States** - Spinners, skeletons, progress bars
- ✅ **Layout** - Container, VStack, HStack, Grid

## How It Works

The example uses:

1. **Go's http.Server** - Standard library HTTP server
2. **Gomponents** - HTML generation in pure Go
3. **ForgeUI Components** - Type-safe UI components
4. **Tailwind CSS Play CDN** - JIT compilation for instant styling

> **Note**: The example uses `https://cdn.tailwindcss.com` which includes JIT compilation. This is perfect for demos but **not recommended for production**. For production, set up a proper Tailwind build process.

## Code Structure

```go
// Create a page
page := html.Html(
    html.Head(
        // Tailwind CSS for styling
        html.Link(html.Rel("stylesheet"), html.Href("...")),
    ),
    html.Body(
        // Use ForgeUI components
        primitives.Container(
            card.Card(
                card.Header(
                    card.Title("Hello"),
                ),
                card.Content(
                    button.Primary(g.Text("Click me")),
                ),
            ),
        ),
    ),
)

// Render to response
page.Render(w)
```

## Customization

### Add More Examples

Edit `main.go` and add new sections using `componentSection()`:

```go
componentSection(
    "Your Section",
    "Description",
    yourComponentDemo,
)
```

### Use Custom Styles

For production, set up a proper Tailwind build:

1. **Install Tailwind:**
```bash
npm install -D tailwindcss
npx tailwindcss init
```

2. **Configure `tailwind.config.js`:**
```js
module.exports = {
  content: ["./main.go", "../components/**/*.go", "../primitives/**/*.go"],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

3. **Build CSS:**
```bash
npx tailwindcss -i ./input.css -o ./static/output.css --watch
```

4. **Replace CDN in `main.go`:**
```go
// Replace:
html.Script(html.Src("https://cdn.tailwindcss.com")),

// With:
html.Link(html.Rel("stylesheet"), html.Href("/static/output.css")),
```

### Add Routes

Add more routes in `main()`:

```go
http.HandleFunc("/about", handleAbout)
```

## Production Tips

For production use:

1. **Use a proper Tailwind build** instead of CDN
2. **Add error handling** for rendering
3. **Implement proper routing** (e.g., with chi, gorilla/mux)
4. **Add HTMX** for interactivity without JavaScript
5. **Use Alpine.js** for client-side interactions (Phase 8)

## Next Steps

- Explore the [main README](../README.md) for API documentation
- Check out component implementations in `../components/`
- Review coding standards in workspace rules

