package templates

import (
	"fmt"
)

// ProjectTemplate defines a project template
type ProjectTemplate interface {
	// Name returns the template name
	Name() string

	// Description returns the template description
	Description() string

	// Generate generates the template in the given directory
	Generate(dir, projectName, modulePath string) error
}

// GetProjectTemplate returns a project template by name
func GetProjectTemplate(name string) (ProjectTemplate, error) {
	switch name {
	case "minimal":
		return &MinimalTemplate{}, nil
	case "standard":
		return &StandardTemplate{}, nil
	case "blog":
		return &BlogTemplate{}, nil
	case "dashboard":
		return &DashboardTemplate{}, nil
	case "api":
		return &APITemplate{}, nil
	default:
		return nil, fmt.Errorf("unknown template: %s", name)
	}
}
