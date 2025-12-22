package templates

import (
	"fmt"
	"path/filepath"
	
	"github.com/xraph/forgeui/cli/util"
)

// MinimalTemplate is a minimal project template
type MinimalTemplate struct{}

func (t *MinimalTemplate) Name() string {
	return "minimal"
}

func (t *MinimalTemplate) Description() string {
	return "Basic setup with one page"
}

func (t *MinimalTemplate) Generate(dir, projectName, modulePath string) error {
	// Create main.go
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"net/http"

	"github.com/xraph/forgeui"
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func main() {
	// Initialize ForgeUI app
	app := forgeui.New(
		forgeui.WithDebug(true),
	)

	// Define home page
	app.Router.Get("/", homePage)

	// Serve static assets
	http.Handle("/static/", app.Assets.Handler())

	// Start server
	fmt.Println("Server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", app); err != nil {
		panic(err)
	}
}

func homePage(ctx *forgeui.PageContext) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("Welcome to %s")),
			html.StyleEl(g.Raw(pageStyles)),
		),
		html.Body(
			html.Div(
				html.Class("container"),
				html.H1(g.Text("Welcome to ForgeUI")),
				html.P(g.Text("Your minimal ForgeUI application is up and running!")),
				html.P(
					g.Text("Edit "),
					html.Code(g.Text("main.go")),
					g.Text(" to get started."),
				),
			),
		),
	)
}

const pageStyles = ` + "`" + `
body {
	margin: 0;
	padding: 0;
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
	background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
}

.container {
	background: white;
	border-radius: 12px;
	padding: 3rem;
	box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
	text-align: center;
	max-width: 600px;
}

h1 {
	color: #667eea;
	margin-top: 0;
	font-size: 2.5rem;
}

p {
	color: #4a5568;
	line-height: 1.6;
	font-size: 1.1rem;
}

code {
	background: #f7fafc;
	padding: 0.2rem 0.5rem;
	border-radius: 4px;
	font-family: 'Courier New', monospace;
	color: #667eea;
}
` + "`" + `
`, projectName)
	
	if err := util.CreateFile(filepath.Join(dir, "main.go"), mainGo); err != nil {
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

A minimal ForgeUI application.

## Getting Started

Start the development server:

` + "```" + `bash
forgeui dev
` + "```" + `

Or run directly with Go:

` + "```" + `bash
go run main.go
` + "```" + `

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Learn More

- [ForgeUI Documentation](https://github.com/xraph/forgeui)
- [Go Documentation](https://go.dev/doc/)
`, projectName)
	
	if err := util.CreateFile(filepath.Join(dir, "README.md"), readme); err != nil {
		return err
	}
	
	return nil
}

