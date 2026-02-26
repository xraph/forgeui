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

	// Create pages/handlers.go
	handlersGo := `package pages

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
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

func BlogHome(ctx *router.PageContext) (templ.Component, error) {
	return BlogHomePage(samplePosts), nil
}

func BlogPost(ctx *router.PageContext) (templ.Component, error) {
	slug := ctx.Params["slug"]
	return BlogPostPage(slug), nil
}

func BlogTag(ctx *router.PageContext) (templ.Component, error) {
	tag := ctx.Params["tag"]
	return BlogTagPage(tag), nil
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "handlers.go"), handlersGo); err != nil {
		return err
	}

	// Create pages/blog.templ
	blogTempl := `package pages

templ BlogHomePage(posts []Post) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<title>My Blog</title>
			@blogStyles()
		</head>
		<body>
			<header><h1>My Blog</h1></header>
			<main>
				<div class="posts">
					for _, p := range posts {
						<article class="post-preview">
							<h2><a href={ templ.SafeURL("/post/" + p.Slug) }>{ p.Title }</a></h2>
							<p>{ p.Content }</p>
						</article>
					}
				</div>
			</main>
		</body>
	</html>
}

templ BlogPostPage(slug string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<title>Post: { slug }</title>
			@blogStyles()
		</head>
		<body>
			<header><h1>Post: { slug }</h1></header>
			<main>
				<p>This is the blog post page.</p>
				<a href="/">&#8592; Back to home</a>
			</main>
		</body>
	</html>
}

templ BlogTagPage(tag string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<title>Tag: { tag }</title>
			@blogStyles()
		</head>
		<body>
			<header><h1>Posts tagged: { tag }</h1></header>
			<main>
				<p>Posts with this tag will appear here.</p>
				<a href="/">&#8592; Back to home</a>
			</main>
		</body>
	</html>
}

templ blogStyles() {
	<style>
		body { font-family: system-ui, -apple-system, sans-serif; max-width: 800px; margin: 0 auto; padding: 2rem; }
		header { margin-bottom: 2rem; }
		h1 { color: #2c3e50; }
		.posts { display: flex; flex-direction: column; gap: 2rem; }
		.post-preview { padding: 1rem; border-left: 3px solid #3498db; }
		.post-preview h2 { margin: 0 0 0.5rem; }
		.post-preview h2 a { color: #2c3e50; text-decoration: none; }
		.post-preview h2 a:hover { color: #3498db; }
	</style>
}
`

	if err := util.CreateFile(filepath.Join(dir, "pages", "blog.templ"), blogTempl); err != nil {
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

	// Create README.md
	readme := fmt.Sprintf(`# %s

A blog template built with ForgeUI.

## Getting Started

`+"```"+`bash
forgeui dev
`+"```"+`

Or run directly:

`+"```"+`bash
templ generate && go run .
`+"```"+`

Visit http://localhost:3000 to see your blog.
`, projectName)

	return util.CreateFile(filepath.Join(dir, "README.md"), readme)
}
