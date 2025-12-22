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

	// Create pages/dashboard.go
	dashGo := `package pages

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func Dashboard(ctx *forgeui.PageContext) g.Node {
	return dashboardLayout("Dashboard", html.Div(
		html.Class("stats"),
		statCard("Users", "1,234"),
		statCard("Revenue", "$56,789"),
		statCard("Orders", "890"),
		statCard("Growth", "+12.5%%"),
	))
}

func Users(ctx *forgeui.PageContext) g.Node {
	return dashboardLayout("Users", html.Div(
		html.H2(g.Text("User Management")),
		html.P(g.Text("User list and management tools will appear here.")),
	))
}

func Analytics(ctx *forgeui.PageContext) g.Node {
	return dashboardLayout("Analytics", html.Div(
		html.H2(g.Text("Analytics Dashboard")),
		html.P(g.Text("Charts and analytics data will appear here.")),
	))
}

func dashboardLayout(title string, content g.Node) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.TitleEl(g.Text(title+" - Dashboard")),
			html.StyleEl(g.Raw(dashStyles)),
		),
		html.Body(
			html.Div(
				html.Class("dashboard"),
				html.Nav(
					html.Class("sidebar"),
					html.H2(g.Text("Dashboard")),
					html.Ul(
						html.Li(html.A(html.Href("/"), g.Text("Overview"))),
						html.Li(html.A(html.Href("/users"), g.Text("Users"))),
						html.Li(html.A(html.Href("/analytics"), g.Text("Analytics"))),
					),
				),
				html.Main(
					html.Class("content"),
					html.H1(g.Text(title)),
					content,
				),
			),
		),
	)
}

func statCard(label, value string) g.Node {
	return html.Div(
		html.Class("stat-card"),
		html.Div(html.Class("stat-label"), g.Text(label)),
		html.Div(html.Class("stat-value"), g.Text(value)),
	)
}

const dashStyles = ` + "`" + `
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
` + "`" + `
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "dashboard.go"), dashGo); err != nil {
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

Visit http://localhost:3000 to see your dashboard.
`, projectName)

	return util.CreateFile(filepath.Join(dir, "README.md"), readme)
}
