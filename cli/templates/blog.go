package templates

import (
	"fmt"
	"path/filepath"

	"github.com/xraph/forgeui/cli/util"
)

// BlogTemplate is a blog template with posts and tags
type BlogTemplate struct{}

func (t *BlogTemplate) Name() string {
	return "blog"
}

func (t *BlogTemplate) Description() string {
	return "Blog template with posts, tags"
}

func (t *BlogTemplate) Generate(dir, projectName, modulePath string) error {
	// Create main.go (simplified blog version)
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"
	"net/http"

	"github.com/xraph/forgeui"
	"%s/pages"
)

func main() {
	app := forgeui.New(forgeui.WithDebug(true))

	// Routes
	app.Router.Get("/", pages.BlogHome)
	app.Router.Get("/post/:slug", pages.BlogPost)
	app.Router.Get("/tag/:tag", pages.BlogTag)

	http.Handle("/static/", app.Assets.Handler())

	fmt.Println("Blog server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", app); err != nil {
		panic(err)
	}
}
`, modulePath)

	if err := util.CreateFile(filepath.Join(dir, "main.go"), mainGo); err != nil {
		return err
	}

	// Create pages/blog.go
	blogGo := `package pages

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

type Post struct {
	Slug    string
	Title   string
	Content string
	Tags    []string
}

var samplePosts = []Post{
	{Slug: "hello-world", Title: "Hello World", Content: "Welcome to my blog!", Tags: []string{"introduction"}},
	{Slug: "go-web-dev", Title: "Go Web Development", Content: "Building web apps with Go is awesome!", Tags: []string{"go", "web"}},
}

func BlogHome(ctx *forgeui.PageContext) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.TitleEl(g.Text("My Blog")),
			html.StyleEl(g.Raw(blogStyles)),
		),
		html.Body(
			html.Header(html.H1(g.Text("My Blog"))),
			html.Main(
				html.Div(
					html.Class("posts"),
					g.Group(g.Map(samplePosts, func(p Post) g.Node {
						return html.Article(
							html.Class("post-preview"),
							html.H2(html.A(html.Href("/post/"+p.Slug), g.Text(p.Title))),
							html.P(g.Text(p.Content)),
						)
					})),
				),
			),
		),
	)
}

func BlogPost(ctx *forgeui.PageContext) g.Node {
	slug := ctx.Params["slug"]
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.TitleEl(g.Text("Post: "+slug)),
			html.StyleEl(g.Raw(blogStyles)),
		),
		html.Body(
			html.Header(html.H1(g.Text("Post: "+slug))),
			html.Main(
				html.P(g.Text("This is the blog post page.")),
				html.A(html.Href("/"), g.Text("← Back to home")),
			),
		),
	)
}

func BlogTag(ctx *forgeui.PageContext) g.Node {
	tag := ctx.Params["tag"]
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.TitleEl(g.Text("Tag: "+tag)),
			html.StyleEl(g.Raw(blogStyles)),
		),
		html.Body(
			html.Header(html.H1(g.Text("Posts tagged: "+tag))),
			html.Main(
				html.P(g.Text("Posts with this tag will appear here.")),
				html.A(html.Href("/"), g.Text("← Back to home")),
			),
		),
	)
}

const blogStyles = ` + "`" + `
body { font-family: system-ui, -apple-system, sans-serif; max-width: 800px; margin: 0 auto; padding: 2rem; }
header { margin-bottom: 2rem; }
h1 { color: #2c3e50; }
.posts { display: flex; flex-direction: column; gap: 2rem; }
.post-preview { padding: 1rem; border-left: 3px solid #3498db; }
.post-preview h2 { margin: 0 0 0.5rem; }
.post-preview h2 a { color: #2c3e50; text-decoration: none; }
.post-preview h2 a:hover { color: #3498db; }
` + "`" + `
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "blog.go"), blogGo); err != nil {
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

	// Create README.md
	readme := fmt.Sprintf(`# %s

A blog template built with ForgeUI.

## Getting Started

`+"```"+`bash
forgeui dev
`+"```"+`

Visit http://localhost:3000 to see your blog.
`, projectName)

	return util.CreateFile(filepath.Join(dir, "README.md"), readme)
}
