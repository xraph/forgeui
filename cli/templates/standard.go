package templates

import (
	"fmt"
	"path/filepath"

	"github.com/xraph/forgeui/cli/util"
)

// StandardTemplate is a standard project template with router and components
type StandardTemplate struct{}

func (t *StandardTemplate) Name() string {
	return "standard"
}

func (t *StandardTemplate) Description() string {
	return "Full setup with router, assets, examples"
}

func (t *StandardTemplate) Generate(dir, projectName, modulePath string) error {
	// Create main.go
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"net/http"

	"github.com/xraph/forgeui"
	"%s/pages"
)

func main() {
	// Initialize ForgeUI app
	app := forgeui.New(
		forgeui.WithDebug(true),
	)

	// Setup routes
	app.Router.Get("/", pages.Home)
	app.Router.Get("/about", pages.About)
	app.Router.Get("/contact", pages.Contact)

	// Serve static assets
	http.Handle("/static/", app.Assets.Handler())

	// Start server
	fmt.Println("Server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", app); err != nil {
		panic(err)
	}
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "main.go"), mainGo); err != nil {
		return err
	}

	// Create pages/handlers.go — page handler functions
	handlersGo := fmt.Sprintf(`package pages

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
	"%s/components"
)

func Home(ctx *router.PageContext) (templ.Component, error) {
	return components.Layout("Home", HomePage()), nil
}

func About(ctx *router.PageContext) (templ.Component, error) {
	return components.Layout("About", AboutPage()), nil
}

func Contact(ctx *router.PageContext) (templ.Component, error) {
	return components.Layout("Contact", ContactPage()), nil
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "pages", "handlers.go"), handlersGo); err != nil {
		return err
	}

	// Create pages/home.templ
	homeTempl := `package pages

templ HomePage() {
	<div class="hero">
		<h1>Welcome to ForgeUI</h1>
		<p>A modern Go UI framework for building beautiful web applications.</p>
		<div class="buttons">
			<a href="/about" class="button primary">Learn More</a>
			<a href="/contact" class="button">Get Started</a>
		</div>
	</div>
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "home.templ"), homeTempl); err != nil {
		return err
	}

	// Create pages/about.templ
	aboutTempl := `package pages

templ AboutPage() {
	<div class="content">
		<h1>About ForgeUI</h1>
		<p>ForgeUI is a modern Go UI framework that makes building web applications fast and enjoyable.</p>
		<h2>Features</h2>
		<ul>
			<li>Component-based architecture</li>
			<li>Type-safe HTML generation with templ</li>
			<li>Built-in router and asset pipeline</li>
			<li>Hot reload development server</li>
			<li>Zero JavaScript required</li>
		</ul>
	</div>
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "about.templ"), aboutTempl); err != nil {
		return err
	}

	// Create pages/contact.templ
	contactTempl := `package pages

templ ContactPage() {
	<div class="content">
		<h1>Contact Us</h1>
		<p>Get in touch with the ForgeUI team.</p>
		<form>
			<div class="form-group">
				<label for="name">Name</label>
				<input type="text" id="name" name="name" class="form-control"/>
			</div>
			<div class="form-group">
				<label for="email">Email</label>
				<input type="email" id="email" name="email" class="form-control"/>
			</div>
			<div class="form-group">
				<label for="message">Message</label>
				<textarea id="message" name="message" class="form-control" rows="5"></textarea>
			</div>
			<button type="submit" class="button primary">Send Message</button>
		</form>
	</div>
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "contact.templ"), contactTempl); err != nil {
		return err
	}

	// Create components/layout.templ
	layoutTempl := `package components

import "github.com/a-h/templ"

templ Layout(title string, content templ.Component) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{ title } - ForgeUI</title>
			@layoutStyles()
		</head>
		<body>
			@Navigation()
			<main>
				@content
			</main>
			@Footer()
		</body>
	</html>
}

templ Navigation() {
	<nav class="nav">
		<div class="nav-content">
			<a href="/" class="nav-brand">ForgeUI</a>
			<div class="nav-links">
				<a href="/">Home</a>
				<a href="/about">About</a>
				<a href="/contact">Contact</a>
			</div>
		</div>
	</nav>
}

templ Footer() {
	<footer class="footer">
		<p>Built with ForgeUI</p>
	</footer>
}

templ layoutStyles() {
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
			line-height: 1.6;
			color: #333;
		}
		.nav { background: #667eea; color: white; padding: 1rem 0; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.nav-content { max-width: 1200px; margin: 0 auto; padding: 0 2rem; display: flex; justify-content: space-between; align-items: center; }
		.nav-brand { font-size: 1.5rem; font-weight: bold; color: white; text-decoration: none; }
		.nav-links { display: flex; gap: 2rem; }
		.nav-links a { color: white; text-decoration: none; transition: opacity 0.2s; }
		.nav-links a:hover { opacity: 0.8; }
		main { min-height: calc(100vh - 120px); }
		.hero { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 6rem 2rem; text-align: center; }
		.hero h1 { font-size: 3rem; margin-bottom: 1rem; }
		.hero p { font-size: 1.25rem; margin-bottom: 2rem; opacity: 0.9; }
		.buttons { display: flex; gap: 1rem; justify-content: center; }
		.button { padding: 0.75rem 2rem; border-radius: 8px; text-decoration: none; font-weight: 500; transition: all 0.2s; border: 2px solid white; color: white; background: transparent; cursor: pointer; }
		.button.primary { background: white; color: #667eea; }
		.button:hover { transform: translateY(-2px); box-shadow: 0 4px 8px rgba(0,0,0,0.2); }
		.content { max-width: 800px; margin: 0 auto; padding: 4rem 2rem; }
		.content h1 { color: #667eea; margin-bottom: 1rem; }
		.content h2 { color: #667eea; margin: 2rem 0 1rem; }
		.content ul { list-style-position: inside; margin: 1rem 0; }
		.content li { margin: 0.5rem 0; }
		.form-group { margin-bottom: 1.5rem; }
		.form-group label { display: block; margin-bottom: 0.5rem; font-weight: 500; }
		.form-control { width: 100%; padding: 0.75rem; border: 2px solid #e2e8f0; border-radius: 8px; font-size: 1rem; font-family: inherit; }
		.form-control:focus { outline: none; border-color: #667eea; }
		.footer { background: #2d3748; color: white; text-align: center; padding: 2rem; }
	</style>
}
`

	if err := util.CreateFile(filepath.Join(dir, "components", "layout.templ"), layoutTempl); err != nil {
		return err
	}

	// Create .forgeui.json
	config := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "dev": {
    "port": 3000,
    "host": "localhost",
    "auto_reload": true,
    "open_browser": false
  },
  "build": {
    "output_dir": "dist",
    "public_dir": "public",
    "minify": true,
    "binary": false,
    "embed_assets": true
  },
  "assets": {
    "css": [],
    "js": []
  },
  "plugins": [],
  "router": {
    "base_path": "/",
    "not_found": ""
  }
}
`, projectName)

	if err := util.CreateFile(filepath.Join(dir, ".forgeui.json"), config); err != nil {
		return err
	}

	// Create .gitignore
	gitignore := `# Binaries
*.exe
*.dll
*.so
*.dylib
bin/
dist/

# Go
*.test
*.out
go.sum

# Templ generated
*_templ.go

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`

	if err := util.CreateFile(filepath.Join(dir, ".gitignore"), gitignore); err != nil {
		return err
	}

	// Create README.md
	readme := fmt.Sprintf(`# %s

A standard ForgeUI application with router, components, and multiple pages.

## Getting Started

Start the development server:

`+"```"+`bash
forgeui dev
`+"```"+`

Or run directly with Go:

`+"```"+`bash
templ generate && go run .
`+"```"+`

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

`+"```"+`
.
├── main.go            # Application entry point
├── pages/             # Page handlers and templates
│   ├── handlers.go    # Page handler functions
│   ├── home.templ     # Home page template
│   ├── about.templ    # About page template
│   └── contact.templ  # Contact page template
├── components/        # Reusable components
│   └── layout.templ   # Layout component
└── public/            # Static assets
    ├── css/
    └── js/
`+"```"+`

## Learn More

- [ForgeUI Documentation](https://github.com/xraph/forgeui)
- [Go Documentation](https://go.dev/doc/)
`, projectName)

	if err := util.CreateFile(filepath.Join(dir, "README.md"), readme); err != nil {
		return err
	}

	return nil
}
