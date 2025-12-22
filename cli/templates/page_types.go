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
	tmpl := `package {{.Package}}

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}} renders the {{.Name}} page
func {{.Name}}(ctx *forgeui.PageContext) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("{{.Name}}")),
		),
		html.Body(
			html.H1(g.Text("{{.Name}}")),
			html.P(g.Text("This is the {{.Name}} page.")),
		),
	)
}
`
	
	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
		"Path":    opts.Path,
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	return util.CreateFile(filePath, code)
}

// DynamicPageTemplate is a page with data loading
type DynamicPageTemplate struct{}

func (t *DynamicPageTemplate) Generate(filePath string, opts PageOptions) error {
	tmpl := `package {{.Package}}

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}}Data holds the page data
type {{.Name}}Data struct {
	Title   string
	Content string
}

// {{.Name}} renders the {{.Name}} page
func {{.Name}}(ctx *forgeui.PageContext) g.Node {
	// Load data
	data := load{{.Name}}Data(ctx)
	
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text(data.Title)),
		),
		html.Body(
			html.H1(g.Text(data.Title)),
			html.Div(g.Text(data.Content)),
		),
	)
}

func load{{.Name}}Data(ctx *forgeui.PageContext) {{.Name}}Data {
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
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	return util.CreateFile(filePath, code)
}

// FormPageTemplate is a page with a form
type FormPageTemplate struct{}

func (t *FormPageTemplate) Generate(filePath string, opts PageOptions) error {
	tmpl := `package {{.Package}}

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}} renders the {{.Name}} form page
func {{.Name}}(ctx *forgeui.PageContext) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("{{.Name}}")),
		),
		html.Body(
			html.H1(g.Text("{{.Name}}")),
			html.Form(
				html.Method("POST"),
				html.Action("{{.Path}}"),
				html.Div(
					html.Label(html.For("name"), g.Text("Name")),
					html.Input(html.Type("text"), html.ID("name"), html.Name("name")),
				),
				html.Div(
					html.Label(html.For("email"), g.Text("Email")),
					html.Input(html.Type("email"), html.ID("email"), html.Name("email")),
				),
				html.Button(html.Type("submit"), g.Text("Submit")),
			),
		),
	)
}
`
	
	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
		"Path":    opts.Path,
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	return util.CreateFile(filePath, code)
}

// ListPageTemplate is a list page with items
type ListPageTemplate struct{}

func (t *ListPageTemplate) Generate(filePath string, opts PageOptions) error {
	tmpl := `package {{.Package}}

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}}Item represents a list item
type {{.Name}}Item struct {
	ID    string
	Title string
}

// {{.Name}} renders the {{.Name}} list page
func {{.Name}}(ctx *forgeui.PageContext) g.Node {
	items := load{{.Name}}Items(ctx)
	
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text("{{.Name}}")),
		),
		html.Body(
			html.H1(g.Text("{{.Name}}")),
			html.Ul(
				g.Group(g.Map(items, func(item {{.Name}}Item) g.Node {
					return html.Li(
						html.A(
							html.Href("{{.Path}}/"+item.ID),
							g.Text(item.Title),
						),
					)
				})),
			),
		),
	)
}

func load{{.Name}}Items(ctx *forgeui.PageContext) []{{.Name}}Item {
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
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	return util.CreateFile(filePath, code)
}

// DetailPageTemplate is a detail page with params
type DetailPageTemplate struct{}

func (t *DetailPageTemplate) Generate(filePath string, opts PageOptions) error {
	tmpl := `package {{.Package}}

import (
	"github.com/xraph/forgeui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}}Detail holds the detail data
type {{.Name}}Detail struct {
	ID          string
	Title       string
	Description string
}

// {{.Name}} renders the {{.Name}} detail page
func {{.Name}}(ctx *forgeui.PageContext) g.Node {
	id := ctx.Params["id"]
	detail := load{{.Name}}Detail(ctx, id)
	
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text(detail.Title)),
		),
		html.Body(
			html.H1(g.Text(detail.Title)),
			html.P(g.Text(detail.Description)),
			html.A(html.Href("{{.Path}}"), g.Text("← Back to list")),
		),
	)
}

func load{{.Name}}Detail(ctx *forgeui.PageContext, id string) {{.Name}}Detail {
	// TODO: Load detail from database or API
	return {{.Name}}Detail{
		ID:          id,
		Title:       "Item " + id,
		Description: "Detail information for item " + id,
	}
}
`
	
	data := map[string]any{
		"Package": opts.Package,
		"Name":    util.ToPascalCase(opts.Name),
		"Path":    strings.TrimSuffix(opts.Path, "/:id"),
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	return util.CreateFile(filePath, code)
}

