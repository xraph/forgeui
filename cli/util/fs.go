package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// CreateDir creates a directory if it doesn't exist
func CreateDir(path string) error {
	if DirExists(path) {
		return nil
	}

	return os.MkdirAll(path, 0755)
}

// CreateFile creates a file with content
func CreateFile(path, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := CreateDir(dir); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0600)
}

// IsGoProject checks if the directory is a Go project
func IsGoProject(dir string) bool {
	return FileExists(filepath.Join(dir, "go.mod"))
}

// IsForgeUIProject checks if the directory is a ForgeUI project
func IsForgeUIProject(dir string) bool {
	return FileExists(filepath.Join(dir, ".forgeui.json")) ||
		FileExists(filepath.Join(dir, "forgeui.json"))
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}

			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteRune(unicode.ToUpper(rune(part[0])))
			result.WriteString(part[1:])
		}
	}

	return result.String()
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) == 0 {
		return ""
	}

	return string(unicode.ToLower(rune(pascal[0]))) + pascal[1:]
}

// GetProjectRoot finds the project root by looking for go.mod
func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree
	for {
		if IsGoProject(dir) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}

		dir = parent
	}

	return "", errors.New("not in a Go project (no go.mod found)")
}

// ValidateProjectName validates a project name
func ValidateProjectName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}

	if strings.Contains(name, " ") {
		return errors.New("project name cannot contain spaces")
	}

	// Check first character is letter or underscore
	if !unicode.IsLetter(rune(name[0])) && name[0] != '_' {
		return errors.New("project name must start with a letter or underscore")
	}

	// Check all characters are valid
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return errors.New("project name can only contain letters, digits, underscores, and hyphens")
		}
	}

	return nil
}

// ValidateGoModule validates a Go module path
func ValidateGoModule(module string) error {
	if module == "" {
		return errors.New("module path cannot be empty")
	}

	// Simple validation - should contain at least one slash
	if !strings.Contains(module, "/") {
		return errors.New("module path should be a valid Go module (e.g., github.com/user/project)")
	}

	return nil
}

// CopyDir recursively copies a directory
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return CreateDir(dstPath)
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}
