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
	mainGo := fmt.Sprintf("package main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n\n\t\"github.com/xraph/forgeui\"\n\t\"github.com/xraph/forgeui/router\"\n\t\"%s/pages\"\n)\n\nfunc main() {\n\t// Initialize ForgeUI app\n\t:= forgeui.New(\n\t\tforgeui.WithDebug(true),\n\t)\n\n\t// Setup routes\n\tapp.Router.Get(\"/\", pages.Home)\n\tapp.Router.Get(\"/about\", pages.About)\n\tapp.Router.Get(\"/contact\", pages.Contact)\n\n\t// Serve static assets\n\thttp.Handle(\"/static/\", app.Assets.Handler())\n\n\t// Start server\n\tfmt.Println(\"Server starting on http://localhost:3000\")\n\tif err := http.ListenAndServe(\":3000\", app); err != nil {\n\t\tpanic(err)\n\t}", modulePath)

	if err := util.CreateFile(filepath.Join(dir, "main.go"), mainGo); err != nil {
		return err
	}

	// Create pages/home.go
	homeGo := fmt.Sprintf(`package pages

import (
	"github.com/xraph/forgeui"
	"%s/components"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Home(ctx *forgeui.PageContext) g.Node {
	return components.Layout(
		"Home",
		html.Div(
			html.Class("hero"),
			html.H1(g.Text("Welcome to ForgeUI")),
			html.P(g.Text("A modern Go UI framework for building beautiful web applications.")),
			html.Div(
				html.Class("buttons"),
				html.A(
					html.Href("/about"),
					html.Class("button primary"),
					g.Text("Learn More"),
				),
				html.A(
					html.Href("/contact"),
					html.Class("button"),
					g.Text("Get Started"),
				),
			),
		),
	)
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "pages", "home.go"), homeGo); err != nil {
		return err
	}

	// Create pages/about.go
	aboutGo := fmt.Sprintf(`package pages

import (
	"github.com/xraph/forgeui"
	"%s/components"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func About(ctx *forgeui.PageContext) g.Node {
	return components.Layout(
		"About",
		html.Div(
			html.Class("content"),
			html.H1(g.Text("About ForgeUI")),
			html.P(g.Text("ForgeUI is a modern Go UI framework that makes building web applications fast and enjoyable.")),
			html.H2(g.Text("Features")),
			html.Ul(
				html.Li(g.Text("Component-based architecture")),
				html.Li(g.Text("Type-safe HTML generation")),
				html.Li(g.Text("Built-in router and asset pipeline")),
				html.Li(g.Text("Hot reload development server")),
				html.Li(g.Text("Zero JavaScript required")),
			),
		),
	)
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "pages", "about.go"), aboutGo); err != nil {
		return err
	}

	// Create pages/contact.go
	contactGo := fmt.Sprintf(`package pages

import (
	"github.com/xraph/forgeui"
	"%s/components"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Contact(ctx *forgeui.PageContext) g.Node {
	return components.Layout(
		"Contact",
		html.Div(
			html.Class("content"),
			html.H1(g.Text("Contact Us")),
			html.P(g.Text("Get in touch with the ForgeUI team.")),
			html.Form(
				html.Div(
					html.Class("form-group"),
					html.Label(html.For("name"), g.Text("Name")),
					html.Input(html.Type("text"), html.ID("name"), html.Name("name"), html.Class("form-control")),
				),
				html.Div(
					html.Class("form-group"),
					html.Label(html.For("email"), g.Text("Email")),
					html.Input(html.Type("email"), html.ID("email"), html.Name("email"), html.Class("form-control")),
				),
				html.Div(
					html.Class("form-group"),
					html.Label(html.For("message"), g.Text("Message")),
					html.Textarea(html.ID("message"), html.Name("message"), html.Class("form-control"), html.Rows("5")),
				),
				html.Button(html.Type("submit"), html.Class("button primary"), g.Text("Send Message")),
			),
		),
	)
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "pages", "contact.go"), contactGo); err != nil {
		return err
	}

	// Create components/layout.go
	layoutGo := `package components

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Layout(title string, content ...g.Node) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text(title + " - ForgeUI")),
			html.StyleEl(g.Raw(layoutStyles)),
		),
		html.Body(
			Navigation(),
			html.Main(
				g.Group(content),
			),
			Footer(),
		),
	)
}

func Navigation() g.Node {
	return html.Nav(
		html.Class("nav"),
		html.Div(
			html.Class("nav-content"),
			html.A(html.Href("/"), html.Class("nav-brand"), g.Text("ForgeUI")),
			html.Div(
				html.Class("nav-links"),
				html.A(html.Href("/"), g.Text("Home")),
				html.A(html.Href("/about"), g.Text("About")),
				html.A(html.Href("/contact"), g.Text("Contact")),
			),
		),
	)
}

func Footer() g.Node {
	return html.Footer(
		html.Class("footer"),
		html.P(g.Text("Built with ForgeUI")),
	)
}

const layoutStyles = ` + "`" + `
* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

body {
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
	line-height: 1.6;
	color: #333;
}

.nav {
	background: #667eea;
	color: white;
	padding: 1rem 0;
	box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.nav-content {
	max-width: 1200px;
	margin: 0 auto;
	padding: 0 2rem;
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.nav-brand {
	font-size: 1.5rem;
	font-weight: bold;
	color: white;
	text-decoration: none;
}

.nav-links {
	display: flex;
	gap: 2rem;
}

.nav-links a {
	color: white;
	text-decoration: none;
	transition: opacity 0.2s;
}

.nav-links a:hover {
	opacity: 0.8;
}

main {
	min-height: calc(100vh - 120px);
}

.hero {
	background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
	color: white;
	padding: 6rem 2rem;
	text-align: center;
}

.hero h1 {
	font-size: 3rem;
	margin-bottom: 1rem;
}

.hero p {
	font-size: 1.25rem;
	margin-bottom: 2rem;
	opacity: 0.9;
}

.buttons {
	display: flex;
	gap: 1rem;
	justify-content: center;
}

.button {
	padding: 0.75rem 2rem;
	border-radius: 8px;
	text-decoration: none;
	font-weight: 500;
	transition: all 0.2s;
	border: 2px solid white;
	color: white;
	background: transparent;
	cursor: pointer;
}

.button.primary {
	background: white;
	color: #667eea;
}

.button:hover {
	transform: translateY(-2px);
	box-shadow: 0 4px 8px rgba(0,0,0,0.2);
}

.content {
	max-width: 800px;
	margin: 0 auto;
	padding: 4rem 2rem;
}

.content h1 {
	color: #667eea;
	margin-bottom: 1rem;
}

.content h2 {
	color: #667eea;
	margin: 2rem 0 1rem;
}

.content ul {
	list-style-position: inside;
	margin: 1rem 0;
}

.content li {
	margin: 0.5rem 0;
}

.form-group {
	margin-bottom: 1.5rem;
}

.form-group label {
	display: block;
	margin-bottom: 0.5rem;
	font-weight: 500;
}

.form-control {
	width: 100%%;
	padding: 0.75rem;
	border: 2px solid #e2e8f0;
	border-radius: 8px;
	font-size: 1rem;
	font-family: inherit;
}

.form-control:focus {
	outline: none;
	border-color: #667eea;
}

.footer {
	background: #2d3748;
	color: white;
	text-align: center;
	padding: 2rem;
}
` + "`" + `
`

	if err := util.CreateFile(filepath.Join(dir, "components", "layout.go"), layoutGo); err != nil {
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
go run main.go
`+"```"+`

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

`+"```"+`
.
├── main.go           # Application entry point
├── pages/            # Page handlers
│   ├── home.go
│   ├── about.go
│   └── contact.go
├── components/       # Reusable components
│   └── layout.go
└── public/           # Static assets
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
