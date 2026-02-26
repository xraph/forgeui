package templates

import (
	"fmt"
	"path/filepath"

	"github.com/xraph/forgeui/cli/util"
)

// DashboardTemplate is an admin dashboard template
type DashboardTemplate struct{}

func (t *DashboardTemplate) Name() string {
	return "dashboard"
}

func (t *DashboardTemplate) Description() string {
	return "Admin dashboard with charts, tables"
}

func (t *DashboardTemplate) Generate(dir, projectName, modulePath string) error {
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

	app.Router.Get("/", pages.Dashboard)
	app.Router.Get("/users", pages.Users)
	app.Router.Get("/analytics", pages.Analytics)

	http.Handle("/static/", app.Assets.Handler())

	fmt.Println("Dashboard server starting on http://localhost:3000")
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

func Dashboard(ctx *router.PageContext) (templ.Component, error) {
	return DashboardPage(), nil
}

func Users(ctx *router.PageContext) (templ.Component, error) {
	return UsersPage(), nil
}

func Analytics(ctx *router.PageContext) (templ.Component, error) {
	return AnalyticsPage(), nil
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "handlers.go"), handlersGo); err != nil {
		return err
	}

	// Create pages/dashboard.templ
	dashTempl := `package pages

templ DashboardPage() {
	@dashboardLayout("Dashboard") {
		<div class="stats">
			@statCard("Users", "1,234")
			@statCard("Revenue", "$56,789")
			@statCard("Orders", "890")
			@statCard("Growth", "+12.5%")
		</div>
	}
}

templ UsersPage() {
	@dashboardLayout("Users") {
		<div>
			<h2>User Management</h2>
			<p>User list and management tools will appear here.</p>
		</div>
	}
}

templ AnalyticsPage() {
	@dashboardLayout("Analytics") {
		<div>
			<h2>Analytics Dashboard</h2>
			<p>Charts and analytics data will appear here.</p>
		</div>
	}
}

templ dashboardLayout(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<title>{ title } - Dashboard</title>
			@dashStyles()
		</head>
		<body>
			<div class="dashboard">
				<nav class="sidebar">
					<h2>Dashboard</h2>
					<ul>
						<li><a href="/">Overview</a></li>
						<li><a href="/users">Users</a></li>
						<li><a href="/analytics">Analytics</a></li>
					</ul>
				</nav>
				<main class="content">
					<h1>{ title }</h1>
					{ children... }
				</main>
			</div>
		</body>
	</html>
}

templ statCard(label, value string) {
	<div class="stat-card">
		<div class="stat-label">{ label }</div>
		<div class="stat-value">{ value }</div>
	</div>
}

templ dashStyles() {
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: system-ui, -apple-system, sans-serif; background: #f5f5f5; }
		.dashboard { display: flex; min-height: 100vh; }
		.sidebar { width: 250px; background: #2c3e50; color: white; padding: 2rem; }
		.sidebar h2 { margin-bottom: 2rem; }
		.sidebar ul { list-style: none; }
		.sidebar li { margin: 0.5rem 0; }
		.sidebar a { color: white; text-decoration: none; display: block; padding: 0.5rem; border-radius: 4px; }
		.sidebar a:hover { background: rgba(255,255,255,0.1); }
		.content { flex: 1; padding: 2rem; }
		.content h1 { margin-bottom: 2rem; color: #2c3e50; }
		.stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; }
		.stat-card { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.stat-label { color: #666; font-size: 0.875rem; margin-bottom: 0.5rem; }
		.stat-value { font-size: 2rem; font-weight: bold; color: #2c3e50; }
	</style>
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "dashboard.templ"), dashTempl); err != nil {
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

An admin dashboard built with ForgeUI.

## Getting Started

`+"```"+`bash
forgeui dev
`+"```"+`

Or run directly:

`+"```"+`bash
templ generate && go run .
`+"```"+`

Visit http://localhost:3000 to see your dashboard.
`, projectName)

	return util.CreateFile(filepath.Join(dir, "README.md"), readme)
}
