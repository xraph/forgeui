package templates

import (
	"fmt"
	"strings"

	"github.com/xraph/forgeui/cli/util"
)

// PageTemplate defines a page template
type PageTemplate interface {
	Generate(filePath string, opts PageOptions) error
}

// PageOptions holds page generation options
type PageOptions struct {
	Name       string
	Package    string
	Path       string
	WithLoader bool
	WithMeta   bool
}

// GetPageTemplate returns a page template by type
func GetPageTemplate(pageType string) (PageTemplate, error) {
	switch pageType {
	case "simple":
		return &SimplePageTemplate{}, nil
	case "dynamic":
		return &DynamicPageTemplate{}, nil
	case "form":
		return &FormPageTemplate{}, nil
	case "list":
		return &ListPageTemplate{}, nil
	case "detail":
		return &DetailPageTemplate{}, nil
	default:
		return nil, fmt.Errorf("unknown page type: %s", pageType)
	}
}

// SimplePageTemplate is a static page
type SimplePageTemplate struct{}

func (t *SimplePageTemplate) Generate(filePath string, opts PageOptions) error {
	// Generate .templ file instead of .go
	templFilePath := strings.TrimSuffix(filePath, ".go") + ".templ"

	tmpl := `package {{.Package}}

import "github.com/xraph/forgeui/router"

templ {{.Name}}Page(ctx *router.PageContext) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{{.Name}}</title>
		</head>
		<body>
			<h1>{{.Name}}</h1>
			<p>This is the {{.Name}} page.</p>
		</body>
	</html>
}
`

	// Also generate handler Go file
	handlerTmpl := `package {{.Package}}

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
)

// {{.Name}} is the page handler for {{.Name}}
func {{.Name}}(ctx *router.PageContext) (templ.Component, error) {
	return {{.Name}}Page(ctx), nil
}
`

	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
		"Path":    opts.Path,
	}

	// Write templ file
	templCode, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(templFilePath, templCode); err != nil {
		return err
	}

	// Write handler Go file
	handlerCode, err := executeTemplate(handlerTmpl, data)
	if err != nil {
		return err
	}
	return util.CreateFile(filePath, handlerCode)
}

// DynamicPageTemplate is a page with data loading
type DynamicPageTemplate struct{}

func (t *DynamicPageTemplate) Generate(filePath string, opts PageOptions) error {
	templFilePath := strings.TrimSuffix(filePath, ".go") + ".templ"

	tmpl := `package {{.Package}}

import "github.com/xraph/forgeui/router"

templ {{.Name}}Page(ctx *router.PageContext, data {{.Name}}Data) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{ data.Title }</title>
		</head>
		<body>
			<h1>{ data.Title }</h1>
			<div>{ data.Content }</div>
		</body>
	</html>
}
`

	handlerTmpl := `package {{.Package}}

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
)

// {{.Name}}Data holds the page data
type {{.Name}}Data struct {
	Title   string
	Content string
}

// {{.Name}} is the page handler for {{.Name}}
func {{.Name}}(ctx *router.PageContext) (templ.Component, error) {
	data := load{{.Name}}Data(ctx)
	return {{.Name}}Page(ctx, data), nil
}

func load{{.Name}}Data(ctx *router.PageContext) {{.Name}}Data {
	// TODO: Load data from database or API
	return {{.Name}}Data{
		Title:   "{{.Name}}",
		Content: "Dynamic content goes here",
	}
}
`

	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
	}

	templCode, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(templFilePath, templCode); err != nil {
		return err
	}

	handlerCode, err := executeTemplate(handlerTmpl, data)
	if err != nil {
		return err
	}
	return util.CreateFile(filePath, handlerCode)
}

// FormPageTemplate is a page with a form
type FormPageTemplate struct{}

func (t *FormPageTemplate) Generate(filePath string, opts PageOptions) error {
	templFilePath := strings.TrimSuffix(filePath, ".go") + ".templ"

	tmpl := `package {{.Package}}

import "github.com/xraph/forgeui/router"

templ {{.Name}}Page(ctx *router.PageContext) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{{.Name}}</title>
		</head>
		<body>
			<h1>{{.Name}}</h1>
			<form method="POST" action="{{.Path}}">
				<div>
					<label for="name">Name</label>
					<input type="text" id="name" name="name"/>
				</div>
				<div>
					<label for="email">Email</label>
					<input type="email" id="email" name="email"/>
				</div>
				<button type="submit">Submit</button>
			</form>
		</body>
	</html>
}
`

	handlerTmpl := `package {{.Package}}

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
)

// {{.Name}} is the page handler for {{.Name}}
func {{.Name}}(ctx *router.PageContext) (templ.Component, error) {
	return {{.Name}}Page(ctx), nil
}
`

	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
		"Path":    opts.Path,
	}

	templCode, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(templFilePath, templCode); err != nil {
		return err
	}

	handlerCode, err := executeTemplate(handlerTmpl, data)
	if err != nil {
		return err
	}
	return util.CreateFile(filePath, handlerCode)
}

// ListPageTemplate is a list page with items
type ListPageTemplate struct{}

func (t *ListPageTemplate) Generate(filePath string, opts PageOptions) error {
	templFilePath := strings.TrimSuffix(filePath, ".go") + ".templ"

	tmpl := `package {{.Package}}

import "github.com/xraph/forgeui/router"

templ {{.Name}}Page(ctx *router.PageContext, items []{{.Name}}Item) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{{.Name}}</title>
		</head>
		<body>
			<h1>{{.Name}}</h1>
			<ul>
				for _, item := range items {
					<li>
						<a href={ templ.SafeURL("{{.Path}}/" + item.ID) }>{ item.Title }</a>
					</li>
				}
			</ul>
		</body>
	</html>
}
`

	handlerTmpl := `package {{.Package}}

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
)

// {{.Name}}Item represents a list item
type {{.Name}}Item struct {
	ID    string
	Title string
}

// {{.Name}} is the page handler for {{.Name}}
func {{.Name}}(ctx *router.PageContext) (templ.Component, error) {
	items := load{{.Name}}Items(ctx)
	return {{.Name}}Page(ctx, items), nil
}

func load{{.Name}}Items(ctx *router.PageContext) []{{.Name}}Item {
	// TODO: Load items from database or API
	return []{{.Name}}Item{
		{ID: "1", Title: "Item 1"},
		{ID: "2", Title: "Item 2"},
		{ID: "3", Title: "Item 3"},
	}
}
`

	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
		"Path":    opts.Path,
	}

	templCode, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(templFilePath, templCode); err != nil {
		return err
	}

	handlerCode, err := executeTemplate(handlerTmpl, data)
	if err != nil {
		return err
	}
	return util.CreateFile(filePath, handlerCode)
}

// DetailPageTemplate is a detail page with params
type DetailPageTemplate struct{}

func (t *DetailPageTemplate) Generate(filePath string, opts PageOptions) error {
	templFilePath := strings.TrimSuffix(filePath, ".go") + ".templ"

	tmpl := `package {{.Package}}

import "github.com/xraph/forgeui/router"

templ {{.Name}}Page(ctx *router.PageContext, detail {{.Name}}Detail) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{ detail.Title }</title>
		</head>
		<body>
			<h1>{ detail.Title }</h1>
			<p>{ detail.Description }</p>
			<a href="{{.BackPath}}">&#8592; Back to list</a>
		</body>
	</html>
}
`

	handlerTmpl := `package {{.Package}}

import (
	"github.com/a-h/templ"
	"github.com/xraph/forgeui/router"
)

// {{.Name}}Detail holds the detail data
type {{.Name}}Detail struct {
	ID          string
	Title       string
	Description string
}

// {{.Name}} is the page handler for {{.Name}}
func {{.Name}}(ctx *router.PageContext) (templ.Component, error) {
	id := ctx.Params["id"]
	detail := load{{.Name}}Detail(ctx, id)
	return {{.Name}}Page(ctx, detail), nil
}

func load{{.Name}}Detail(ctx *router.PageContext, id string) {{.Name}}Detail {
	// TODO: Load detail from database or API
	return {{.Name}}Detail{
		ID:          id,
		Title:       "Item " + id,
		Description: "Detail information for item " + id,
	}
}
`

	data := map[string]any{
		"Package":  opts.Package,
		"Name":     util.ToPascalCase(opts.Name),
		"Path":     opts.Path,
		"BackPath": strings.TrimSuffix(opts.Path, "/:id"),
	}

	templCode, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(templFilePath, templCode); err != nil {
		return err
	}

	handlerCode, err := executeTemplate(handlerTmpl, data)
	if err != nil {
		return err
	}
	return util.CreateFile(filePath, handlerCode)
}
