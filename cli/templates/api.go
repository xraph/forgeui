package templates

import (
	"fmt"
	"path/filepath"

	"github.com/xraph/forgeui/cli/util"
)

// APITemplate is an API-first template with HTMX
type APITemplate struct{}

func (t *APITemplate) Name() string {
	return "api"
}

func (t *APITemplate) Description() string {
	return "API-first template with HTMX"
}

func (t *APITemplate) Generate(dir, projectName, modulePath string) error {
	// Create main.go
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"net/http"

	"github.com/xraph/forgeui"
	"%s/pages"
)

func main() {
	app := forgeui.New(forgeui.WithDebug(true))

	// Pages
	app.Router.Get("/", pages.Home)

	// API endpoints
	app.Router.Get("/api/users", pages.APIUsers)
	app.Router.Post("/api/users", pages.APICreateUser)

	http.Handle("/static/", app.Assets.Handler())

	fmt.Println("API server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", app); err != nil {
		panic(err)
	}
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "main.go"), mainGo); err != nil {
		return err
	}

	// Create pages/handlers.go
	handlersGo := `package pages

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
)

func Home(ctx *router.PageContext) (templ.Component, error) {
	return HomePage(), nil
}

func APIUsers(ctx *router.PageContext) (templ.Component, error) {
	users := []string{"Alice", "Bob", "Charlie"}
	return UserList(users), nil
}

func APICreateUser(ctx *router.PageContext) (templ.Component, error) {
	return UserCreated(), nil
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "handlers.go"), handlersGo); err != nil {
		return err
	}

	// Create pages/api.templ
	apiTempl := `package pages

import "github.com/xraph/forgeui/htmx"

templ HomePage() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<title>HTMX API Demo</title>
			@htmx.Scripts()
			@apiStyles()
		</head>
		<body>
			<h1>HTMX API Demo</h1>
			<div id="user-list" { htmx.HxGet("/api/users")... } { htmx.HxTrigger("load")... }>
				<p>Loading users...</p>
			</div>
		</body>
	</html>
}

templ UserList(users []string) {
	<ul>
		for _, name := range users {
			<li>{ name }</li>
		}
	</ul>
}

templ UserCreated() {
	<p>User created!</p>
}

templ apiStyles() {
	<style>
		body { font-family: system-ui, -apple-system, sans-serif; max-width: 800px; margin: 0 auto; padding: 2rem; }
		h1 { color: #2c3e50; margin-bottom: 2rem; }
		#user-list { background: #f5f5f5; padding: 1rem; border-radius: 8px; }
		ul { list-style: none; padding: 0; }
		li { padding: 0.5rem; background: white; margin: 0.5rem 0; border-radius: 4px; }
	</style>
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "api.templ"), apiTempl); err != nil {
		return err
	}

	// Create config
	config := fmt.Sprintf(`{"name":"%s","version":"1.0.0","dev":{"port":3000,"host":"localhost","auto_reload":true,"open_browser":false},"build":{"output_dir":"dist","public_dir":"public","minify":true,"binary":false,"embed_assets":true},"assets":{"css":[],"js":[]},"plugins":[],"router":{"base_path":"/","not_found":""}}`, projectName)

	if err := util.CreateFile(filepath.Join(dir, ".forgeui.json"), config); err != nil {
		return err
	}

	// Create .gitignore
	gitignore := `*.exe
*.dll
*.so
*.dylib
bin/
dist/
*.test
*.out
go.sum
*_templ.go
.vscode/
.idea/
*.swp
*.swo
*~
.DS_Store
Thumbs.db
`

	if err := util.CreateFile(filepath.Join(dir, ".gitignore"), gitignore); err != nil {
		return err
	}

	readme := fmt.Sprintf(`# %s

An API-first application with HTMX built with ForgeUI.

## Getting Started

`+"```"+`bash
forgeui dev
`+"```"+`

Or run directly:

`+"```"+`bash
templ generate && go run .
`+"```"+`

Visit http://localhost:3000 to see your app.
`, projectName)

	return util.CreateFile(filepath.Join(dir, "README.md"), readme)
}
