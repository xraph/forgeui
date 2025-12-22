package templates

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	
	"github.com/xraph/forgeui/cli/util"
)

// ComponentTemplate defines a component template
type ComponentTemplate interface {
	Generate(dir string, opts ComponentOptions) error
}

// ComponentOptions holds component generation options
type ComponentOptions struct {
	Name         string
	Package      string
	WithVariants bool
	WithProps    bool
	WithTest     bool
}

// GetComponentTemplate returns a component template by type
func GetComponentTemplate(componentType string) (ComponentTemplate, error) {
	switch componentType {
	case "basic":
		return &BasicComponentTemplate{}, nil
	case "compound":
		return &CompoundComponentTemplate{}, nil
	case "form":
		return &FormComponentTemplate{}, nil
	case "layout":
		return &LayoutComponentTemplate{}, nil
	case "data":
		return &DataComponentTemplate{}, nil
	default:
		return nil, fmt.Errorf("unknown component type: %s", componentType)
	}
}

// BasicComponentTemplate is a simple functional component
type BasicComponentTemplate struct{}

func (t *BasicComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	tmpl := `package {{.Package}}

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
{{- if .WithProps}}
	"github.com/xraph/forgeui"
{{- end}}
)

{{- if .WithProps}}

// {{.Name}}Props defines properties for {{.Name}}
type {{.Name}}Props struct {
	Class string
	ID    string
}
{{- end}}

// {{.Name}} renders a {{.Name}} component
func {{.Name}}({{if .WithProps}}props {{.Name}}Props, {{end}}children ...g.Node) g.Node {
	return html.Div(
		{{- if .WithProps}}
		g.If(props.Class != "", html.Class(props.Class)),
		g.If(props.ID != "", html.ID(props.ID)),
		{{- else}}
		html.Class("{{.PackageName}}"),
		{{- end}}
		g.Group(children),
	)
}
`
	
	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
		"WithProps":   opts.WithProps,
		"WithVariants": opts.WithVariants,
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	fileName := opts.Package + ".go"
	if err := util.CreateFile(filepath.Join(dir, fileName), code); err != nil {
		return err
	}
	
	if opts.WithTest {
		testCode := generateTestFile(opts)
		testFileName := opts.Package + "_test.go"
		if err := util.CreateFile(filepath.Join(dir, testFileName), testCode); err != nil {
			return err
		}
	}
	
	return nil
}

// CompoundComponentTemplate is a compound component
type CompoundComponentTemplate struct{}

func (t *CompoundComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	tmpl := `package {{.Package}}

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}} renders a {{.Name}} component
func {{.Name}}(children ...g.Node) g.Node {
	return html.Div(
		html.Class("{{.PackageName}}"),
		g.Group(children),
	)
}

// {{.Name}}Header renders the header part
func {{.Name}}Header(children ...g.Node) g.Node {
	return html.Div(
		html.Class("{{.PackageName}}-header"),
		g.Group(children),
	)
}

// {{.Name}}Body renders the body part
func {{.Name}}Body(children ...g.Node) g.Node {
	return html.Div(
		html.Class("{{.PackageName}}-body"),
		g.Group(children),
	)
}

// {{.Name}}Footer renders the footer part
func {{.Name}}Footer(children ...g.Node) g.Node {
	return html.Div(
		html.Class("{{.PackageName}}-footer"),
		g.Group(children),
	)
}
`
	
	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	fileName := opts.Package + ".go"
	return util.CreateFile(filepath.Join(dir, fileName), code)
}

// FormComponentTemplate is a form component
type FormComponentTemplate struct{}

func (t *FormComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	tmpl := `package {{.Package}}

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}}Props defines properties for the form
type {{.Name}}Props struct {
	Action string
	Method string
}

// {{.Name}} renders a form component
func {{.Name}}(props {{.Name}}Props, children ...g.Node) g.Node {
	return html.Form(
		html.Class("{{.PackageName}}"),
		g.If(props.Action != "", html.Action(props.Action)),
		g.If(props.Method != "", html.Method(props.Method)),
		g.Group(children),
	)
}
`
	
	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	fileName := opts.Package + ".go"
	return util.CreateFile(filepath.Join(dir, fileName), code)
}

// LayoutComponentTemplate is a layout component
type LayoutComponentTemplate struct{}

func (t *LayoutComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	tmpl := `package {{.Package}}

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}} renders a layout component
func {{.Name}}(title string, children ...g.Node) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text(title)),
		),
		html.Body(
			html.Main(
				html.Class("{{.PackageName}}"),
				g.Group(children),
			),
		),
	)
}
`
	
	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	fileName := opts.Package + ".go"
	return util.CreateFile(filepath.Join(dir, fileName), code)
}

// DataComponentTemplate is a data display component
type DataComponentTemplate struct{}

func (t *DataComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	tmpl := `package {{.Package}}

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// {{.Name}}Item represents a data item
type {{.Name}}Item struct {
	ID    string
	Title string
	Value string
}

// {{.Name}} renders a data display component
func {{.Name}}(items []{{.Name}}Item) g.Node {
	return html.Div(
		html.Class("{{.PackageName}}"),
		g.Group(g.Map(items, func(item {{.Name}}Item) g.Node {
			return html.Div(
				html.Class("{{.PackageName}}-item"),
				html.Div(html.Class("title"), g.Text(item.Title)),
				html.Div(html.Class("value"), g.Text(item.Value)),
			)
		})),
	)
}
`
	
	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
	}
	
	code, err := executeTemplate(tmpl, data)
	if err != nil {
		return err
	}
	
	fileName := opts.Package + ".go"
	return util.CreateFile(filepath.Join(dir, fileName), code)
}

func executeTemplate(tmplStr string, data map[string]any) (string, error) {
	tmpl, err := template.New("component").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

func generateTestFile(opts ComponentOptions) string {
	return fmt.Sprintf(`package %s

import (
	"testing"

	g "maragu.dev/gomponents"
)

func Test%s(t *testing.T) {
	component := %s()
	
	if component == nil {
		t.Error("component should not be nil")
	}
	
	// Add more tests here
}
`, opts.Package, util.ToPascalCase(opts.Name), util.ToPascalCase(opts.Name))
}

