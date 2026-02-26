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
	propsGoTmpl := `package {{.Package}}

// {{.Name}}Props defines properties for {{.Name}}
type {{.Name}}Props struct {
	Class string
	ID    string
}
`

	templWithPropsTmpl := `package {{.Package}}

templ {{.Name}}(props {{.Name}}Props) {
	<div
		if props.ID != "" {
			id={ props.ID }
		}
		if props.Class != "" {
			class={ props.Class }
		} else {
			class="{{.PackageName}}"
		}
	>
		{ children... }
	</div>
}
`

	templNoPropsTmpl := `package {{.Package}}

templ {{.Name}}() {
	<div class="{{.PackageName}}">
		{ children... }
	</div>
}
`

	data := map[string]any{
		"Package":      opts.Package,
		"PackageName":  opts.Package,
		"Name":         util.ToPascalCase(opts.Name),
		"WithProps":    opts.WithProps,
		"WithVariants": opts.WithVariants,
	}

	// Write props Go file if needed
	if opts.WithProps {
		code, err := executeTemplate(propsGoTmpl, data)
		if err != nil {
			return err
		}
		propsFileName := opts.Package + "_props.go"
		if err := util.CreateFile(filepath.Join(dir, propsFileName), code); err != nil {
			return err
		}
	}

	// Write templ file
	var templTmpl string
	if opts.WithProps {
		templTmpl = templWithPropsTmpl
	} else {
		templTmpl = templNoPropsTmpl
	}

	templCode, err := executeTemplate(templTmpl, data)
	if err != nil {
		return err
	}

	templFileName := opts.Package + ".templ"
	if err := util.CreateFile(filepath.Join(dir, templFileName), templCode); err != nil {
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

templ {{.Name}}() {
	<div class="{{.PackageName}}">
		{ children... }
	</div>
}

templ {{.Name}}Header() {
	<div class="{{.PackageName}}-header">
		{ children... }
	</div>
}

templ {{.Name}}Body() {
	<div class="{{.PackageName}}-body">
		{ children... }
	</div>
}

templ {{.Name}}Footer() {
	<div class="{{.PackageName}}-footer">
		{ children... }
	</div>
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

	fileName := opts.Package + ".templ"

	return util.CreateFile(filepath.Join(dir, fileName), code)
}

// FormComponentTemplate is a form component
type FormComponentTemplate struct{}

func (t *FormComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	goTmpl := `package {{.Package}}

// {{.Name}}Props defines properties for the form
type {{.Name}}Props struct {
	Action string
	Method string
}
`

	templTmpl := `package {{.Package}}

templ {{.Name}}(props {{.Name}}Props) {
	<form class="{{.PackageName}}"
		if props.Action != "" {
			action={ templ.SafeURL(props.Action) }
		}
		if props.Method != "" {
			method={ props.Method }
		}
	>
		{ children... }
	</form>
}
`

	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
	}

	// Write props Go file
	goCode, err := executeTemplate(goTmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(filepath.Join(dir, opts.Package+"_props.go"), goCode); err != nil {
		return err
	}

	// Write templ file
	templCode, err := executeTemplate(templTmpl, data)
	if err != nil {
		return err
	}

	return util.CreateFile(filepath.Join(dir, opts.Package+".templ"), templCode)
}

// LayoutComponentTemplate is a layout component
type LayoutComponentTemplate struct{}

func (t *LayoutComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	tmpl := `package {{.Package}}

templ {{.Name}}(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1"/>
			<title>{ title }</title>
		</head>
		<body>
			<main class="{{.PackageName}}">
				{ children... }
			</main>
		</body>
	</html>
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

	fileName := opts.Package + ".templ"

	return util.CreateFile(filepath.Join(dir, fileName), code)
}

// DataComponentTemplate is a data display component
type DataComponentTemplate struct{}

func (t *DataComponentTemplate) Generate(dir string, opts ComponentOptions) error {
	goTmpl := `package {{.Package}}

// {{.Name}}Item represents a data item
type {{.Name}}Item struct {
	ID    string
	Title string
	Value string
}
`

	templTmpl := `package {{.Package}}

templ {{.Name}}(items []{{.Name}}Item) {
	<div class="{{.PackageName}}">
		for _, item := range items {
			<div class="{{.PackageName}}-item">
				<div class="title">{ item.Title }</div>
				<div class="value">{ item.Value }</div>
			</div>
		}
	</div>
}
`

	data := map[string]any{
		"Package":     opts.Package,
		"PackageName": opts.Package,
		"Name":        util.ToPascalCase(opts.Name),
	}

	// Write Go types file
	goCode, err := executeTemplate(goTmpl, data)
	if err != nil {
		return err
	}
	if err := util.CreateFile(filepath.Join(dir, opts.Package+"_types.go"), goCode); err != nil {
		return err
	}

	// Write templ file
	templCode, err := executeTemplate(templTmpl, data)
	if err != nil {
		return err
	}

	return util.CreateFile(filepath.Join(dir, opts.Package+".templ"), templCode)
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
	"bytes"
	"context"
	"testing"
)

func Test%s(t *testing.T) {
	var buf bytes.Buffer
	err := %s().Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render error: %%v", err)
	}

	if buf.Len() == 0 {
		t.Error("component rendered empty output")
	}
}
`, opts.Package, util.ToPascalCase(opts.Name), util.ToPascalCase(opts.Name))
}
