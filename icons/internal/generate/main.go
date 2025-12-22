//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

const (
	// Lucide icons metadata from jsdelivr CDN
	lucideIconNodesURL = "https://cdn.jsdelivr.net/npm/lucide-static@latest/icon-nodes.json"
	lucideTagsURL      = "https://cdn.jsdelivr.net/npm/lucide-static@latest/tags.json"
	outputFile         = "lucide_generated.go"
)

// IconNode represents a single SVG element in the icon
type IconNode []interface{}

// IconData represents processed icon information
type IconData struct {
	Name         string   // Original kebab-case name
	FunctionName string   // PascalCase function name
	Paths        []string // SVG path data
	Tags         []string
}

func main() {
	fmt.Println("🎨 Lucide Icon Generator for ForgeUI")
	fmt.Println("=====================================")

	// Fetch icon nodes
	fmt.Println("\n📥 Fetching Lucide icon nodes...")
	iconNodes, err := fetchIconNodes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching icon nodes: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Fetched %d icons\n", len(iconNodes))

	// Fetch tags
	fmt.Println("\n📥 Fetching icon tags...")
	tags, err := fetchIconTags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not fetch tags: %v\n", err)
		tags = make(map[string][]string)
	}

	// Process icons
	fmt.Println("\n🔄 Processing icons...")
	icons := processIcons(iconNodes, tags)
	fmt.Printf("✅ Processed %d icons\n", len(icons))

	// Generate Go code
	fmt.Println("\n📝 Generating Go code...")
	code, err := generateGoCode(icons)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	outputPath := filepath.Join("..", "..", outputFile)
	if err := os.WriteFile(outputPath, code, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Generated %s with %d icons\n", outputFile, len(icons))
	fmt.Println("\n✨ Generation complete!")
}

// fetchIconNodes fetches icon node data from the Lucide CDN
func fetchIconNodes() (map[string][]IconNode, error) {
	resp, err := http.Get(lucideIconNodesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch icon nodes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var iconNodes map[string][]IconNode
	if err := json.Unmarshal(body, &iconNodes); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return iconNodes, nil
}

// fetchIconTags fetches icon tags from the Lucide CDN
func fetchIconTags() (map[string][]string, error) {
	resp, err := http.Get(lucideTagsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tags map[string][]string
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return tags, nil
}

// processIcons converts raw icon nodes into processed icon data
func processIcons(iconNodes map[string][]IconNode, tags map[string][]string) []IconData {
	icons := make([]IconData, 0, len(iconNodes))

	for name, nodes := range iconNodes {
		// Extract path data from nodes
		paths := extractPathsFromNodes(nodes)
		if len(paths) == 0 {
			fmt.Printf("⚠️  Warning: No paths found for icon '%s', skipping\n", name)
			continue
		}

		// Convert name to valid Go function name
		funcName := toFunctionName(name)

		// Get tags for this icon
		iconTags := tags[name]

		icons = append(icons, IconData{
			Name:         name,
			FunctionName: funcName,
			Paths:        paths,
			Tags:         iconTags,
		})
	}

	// Sort by function name for consistent output
	sort.Slice(icons, func(i, j int) bool {
		return icons[i].FunctionName < icons[j].FunctionName
	})

	return icons
}

// extractPathsFromNodes extracts path 'd' attributes from icon nodes
// and converts other SVG shapes (circle, rect, line, ellipse) to path commands
func extractPathsFromNodes(nodes []IconNode) []string {
	paths := make([]string, 0)

	for _, node := range nodes {
		if len(node) < 2 {
			continue
		}

		// First element is the tag name
		tagName, ok := node[0].(string)
		if !ok {
			continue
		}

		// Second element is the attributes map
		attrs, ok := node[1].(map[string]interface{})
		if !ok {
			continue
		}

		// For path elements, extract the 'd' attribute
		if tagName == "path" {
			if d, ok := attrs["d"].(string); ok && d != "" {
				paths = append(paths, d)
			}
			continue
		}

		// For polyline and polygon, check for both 'd' and 'points' attributes
		if tagName == "polyline" || tagName == "polygon" {
			// First try 'd' attribute
			if d, ok := attrs["d"].(string); ok && d != "" {
				paths = append(paths, d)
				continue
			}
			// Try 'points' attribute and convert to path
			if points, ok := attrs["points"].(string); ok && points != "" {
				if pathData := pointsToPath(points, tagName == "polygon"); pathData != "" {
					paths = append(paths, pathData)
				}
			}
			continue
		}

		// Convert other shapes to path commands
		switch tagName {
		case "circle":
			if pathData := circleToPath(attrs); pathData != "" {
				paths = append(paths, pathData)
			}
		case "rect":
			if pathData := rectToPath(attrs); pathData != "" {
				paths = append(paths, pathData)
			}
		case "line":
			if pathData := lineToPath(attrs); pathData != "" {
				paths = append(paths, pathData)
			}
		case "ellipse":
			if pathData := ellipseToPath(attrs); pathData != "" {
				paths = append(paths, pathData)
			}
		}
	}

	return paths
}

// pointsToPath converts a points string (from polygon/polyline) to a path command
func pointsToPath(points string, closed bool) string {
	// Parse points string (space or comma separated coordinates)
	points = strings.TrimSpace(points)
	// Replace commas with spaces and normalize whitespace
	points = regexp.MustCompile(`[,\s]+`).ReplaceAllString(points, " ")

	coords := strings.Fields(points)
	if len(coords) < 2 {
		return ""
	}

	// Build path command
	pathParts := make([]string, 0, len(coords)/2+2)
	pathParts = append(pathParts, "M", coords[0], coords[1])

	for i := 2; i < len(coords)-1; i += 2 {
		pathParts = append(pathParts, "L", coords[i], coords[i+1])
	}

	// Close path for polygons
	if closed {
		pathParts = append(pathParts, "Z")
	}

	return strings.Join(pathParts, " ")
}

// circleToPath converts a circle element to a path command
func circleToPath(attrs map[string]interface{}) string {
	cx := getFloatAttr(attrs, "cx")
	cy := getFloatAttr(attrs, "cy")
	r := getFloatAttr(attrs, "r")

	if r <= 0 {
		return ""
	}

	// Create a circle using two arc commands
	return fmt.Sprintf("M %g %g a %g %g 0 1 0 %g 0 a %g %g 0 1 0 -%g 0",
		cx-r, cy, r, r, 2*r, r, r, 2*r)
}

// rectToPath converts a rect element to a path command
func rectToPath(attrs map[string]interface{}) string {
	x := getFloatAttr(attrs, "x")
	y := getFloatAttr(attrs, "y")
	width := getFloatAttr(attrs, "width")
	height := getFloatAttr(attrs, "height")
	rx := getFloatAttr(attrs, "rx")
	ry := getFloatAttr(attrs, "ry")

	if width <= 0 || height <= 0 {
		return ""
	}

	// Handle rounded corners
	if rx > 0 || ry > 0 {
		if rx == 0 {
			rx = ry
		}
		if ry == 0 {
			ry = rx
		}
		// Limit radius to half the smallest dimension
		maxR := width / 2
		if height/2 < maxR {
			maxR = height / 2
		}
		if rx > maxR {
			rx = maxR
		}
		if ry > maxR {
			ry = maxR
		}

		return fmt.Sprintf("M %g %g h %g a %g %g 0 0 1 %g %g v %g a %g %g 0 0 1 -%g %g h -%g a %g %g 0 0 1 -%g -%g v -%g a %g %g 0 0 1 %g -%g Z",
			x+rx, y, width-2*rx, rx, ry, rx, ry, height-2*ry, rx, ry, rx, ry, width-2*rx, rx, ry, rx, ry, height-2*ry, rx, ry, rx, ry)
	}

	// Simple rectangle
	return fmt.Sprintf("M %g %g h %g v %g h -%g Z", x, y, width, height, width)
}

// lineToPath converts a line element to a path command
func lineToPath(attrs map[string]interface{}) string {
	x1 := getFloatAttr(attrs, "x1")
	y1 := getFloatAttr(attrs, "y1")
	x2 := getFloatAttr(attrs, "x2")
	y2 := getFloatAttr(attrs, "y2")

	return fmt.Sprintf("M %g %g L %g %g", x1, y1, x2, y2)
}

// ellipseToPath converts an ellipse element to a path command
func ellipseToPath(attrs map[string]interface{}) string {
	cx := getFloatAttr(attrs, "cx")
	cy := getFloatAttr(attrs, "cy")
	rx := getFloatAttr(attrs, "rx")
	ry := getFloatAttr(attrs, "ry")

	if rx <= 0 || ry <= 0 {
		return ""
	}

	// Create an ellipse using two arc commands
	return fmt.Sprintf("M %g %g a %g %g 0 1 0 %g 0 a %g %g 0 1 0 -%g 0",
		cx-rx, cy, rx, ry, 2*rx, rx, ry, 2*rx)
}

// getFloatAttr safely extracts a float attribute value
func getFloatAttr(attrs map[string]interface{}, key string) float64 {
	val, ok := attrs[key]
	if !ok {
		return 0
	}

	switch v := val.(type) {
	case float64:
		return v
	case string:
		// Parse string to float
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	default:
		return 0
	}
}

// toFunctionName converts kebab-case to PascalCase and handles special cases
func toFunctionName(name string) string {
	// Handle special prefixes (numbers)
	if len(name) > 0 && unicode.IsDigit(rune(name[0])) {
		// Convert leading numbers to words
		name = numberPrefix(name)
	}

	// Split on hyphens and capitalize each part
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		// Capitalize first letter
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}

	funcName := strings.Join(parts, "")

	// Handle any remaining invalid characters
	funcName = regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(funcName, "")

	// Ensure it doesn't conflict with Go keywords
	if isGoKeyword(funcName) {
		funcName = funcName + "Icon"
	}

	// Ensure it doesn't conflict with package-level types
	if isPackageReserved(funcName) {
		funcName = funcName + "Icon"
	}

	return funcName
}

// numberPrefix converts leading numbers to word prefixes
func numberPrefix(name string) string {
	numberWords := map[string]string{
		"0": "zero",
		"1": "one",
		"2": "two",
		"3": "three",
		"4": "four",
		"5": "five",
		"6": "six",
		"7": "seven",
		"8": "eight",
		"9": "nine",
	}

	// Find the first non-digit
	i := 0
	for i < len(name) && unicode.IsDigit(rune(name[i])) {
		i++
	}

	if i == 0 {
		return name
	}

	prefix := name[:i]
	rest := name[i:]

	// Convert number to words
	words := ""
	for _, digit := range prefix {
		if word, ok := numberWords[string(digit)]; ok {
			words += word + "-"
		}
	}

	return strings.TrimSuffix(words, "-") + rest
}

// isGoKeyword checks if a name is a Go keyword
func isGoKeyword(name string) bool {
	keywords := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true,
		"continue": true, "default": true, "defer": true, "else": true,
		"fallthrough": true, "for": true, "func": true, "go": true,
		"goto": true, "if": true, "import": true, "interface": true,
		"map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true,
		"var": true,
	}
	return keywords[strings.ToLower(name)]
}

// isPackageReserved checks if a name conflicts with package-level types/functions
func isPackageReserved(name string) bool {
	reserved := map[string]bool{
		"Option":          true, // Functional option type
		"Props":           true, // Icon props struct
		"Icon":            true, // Icon function
		"MultiPathIcon":   true, // MultiPathIcon function
		"WithSize":        true,
		"WithColor":       true,
		"WithStrokeWidth": true,
		"WithClass":       true,
		"WithAttrs":       true,
	}
	return reserved[name]
}

// generateGoCode generates the complete Go source file
func generateGoCode(icons []IconData) ([]byte, error) {
	var buf bytes.Buffer

	// Write package header
	buf.WriteString("// Code generated by icons/internal/generate/main.go. DO NOT EDIT.\n")
	buf.WriteString("// To regenerate: cd icons/internal/generate && go run main.go\n\n")
	buf.WriteString("package icons\n\n")
	buf.WriteString("import g \"maragu.dev/gomponents\"\n\n")
	buf.WriteString(fmt.Sprintf("// This file contains %d auto-generated icon functions from Lucide Icons\n", len(icons)))
	buf.WriteString("// See https://lucide.dev for the complete icon reference\n\n")

	// Generate function for each icon
	for _, icon := range icons {
		// Function documentation
		buf.WriteString(fmt.Sprintf("// %s creates a %s icon\n", icon.FunctionName, icon.Name))
		if len(icon.Tags) > 0 {
			buf.WriteString(fmt.Sprintf("// Tags: %s\n", strings.Join(icon.Tags, ", ")))
		}

		// Function signature
		buf.WriteString(fmt.Sprintf("func %s(opts ...Option) g.Node {\n", icon.FunctionName))

		// Function body
		if len(icon.Paths) == 1 {
			// Single path icon
			buf.WriteString(fmt.Sprintf("\treturn Icon(%q, opts...)\n", icon.Paths[0]))
		} else {
			// Multi-path icon
			buf.WriteString("\treturn MultiPathIcon([]string{\n")
			for _, path := range icon.Paths {
				buf.WriteString(fmt.Sprintf("\t\t%q,\n", path))
			}
			buf.WriteString("\t}, opts...)\n")
		}

		buf.WriteString("}\n\n")
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, return the unformatted code and the error
		// This helps debug syntax issues
		fmt.Fprintf(os.Stderr, "Warning: Failed to format generated code: %v\n", err)
		fmt.Fprintf(os.Stderr, "Writing unformatted code for debugging\n")
		return buf.Bytes(), nil
	}

	return formatted, nil
}
