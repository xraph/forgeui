package plugin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Satisfies checks if a version satisfies the constraint.
// Supports: =, >=, <=, >, <, ~, ^
func (d Dependency) Satisfies(version string) bool {
	if d.Version == "" || d.Version == "*" {
		return true
	}

	constraint, err := ParseConstraint(d.Version)
	if err != nil {
		return false
	}

	v, err := ParseVersion(version)
	if err != nil {
		return false
	}

	return constraint.Check(v)
}

// Version represents a semantic version.
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a semantic version string.
func ParseVersion(s string) (*Version, error) {
	// Remove 'v' prefix if present
	s = strings.TrimPrefix(s, "v")

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid version format: %s (expected major.minor.patch)", s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// Compare compares two versions.
// Returns: -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	return 0
}

// String returns the version as a string.
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// VersionConstraint represents a version constraint.
type VersionConstraint struct {
	Operator string
	Version  *Version
}

var constraintRegex = regexp.MustCompile(`^(>=|<=|>|<|=|~|\^)?(.+)$`)

// ParseConstraint parses a version constraint string.
func ParseConstraint(constraint string) (*VersionConstraint, error) {
	constraint = strings.TrimSpace(constraint)

	matches := constraintRegex.FindStringSubmatch(constraint)
	if matches == nil {
		return nil, fmt.Errorf("invalid constraint format: %s", constraint)
	}

	operator := matches[1]
	if operator == "" {
		operator = "="
	}

	version, err := ParseVersion(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid version in constraint: %w", err)
	}

	return &VersionConstraint{
		Operator: operator,
		Version:  version,
	}, nil
}

// Check checks if a version satisfies the constraint.
func (c *VersionConstraint) Check(version *Version) bool {
	cmp := version.Compare(c.Version)

	switch c.Operator {
	case "=":
		return cmp == 0
	case ">":
		return cmp > 0
	case ">=":
		return cmp >= 0
	case "<":
		return cmp < 0
	case "<=":
		return cmp <= 0
	case "~":
		// ~1.2.3 means >=1.2.3 and <1.3.0
		if version.Major != c.Version.Major || version.Minor != c.Version.Minor {
			return false
		}
		return version.Patch >= c.Version.Patch
	case "^":
		// ^1.2.3 means >=1.2.3 and <2.0.0
		if version.Major != c.Version.Major {
			return false
		}
		return cmp >= 0
	default:
		return false
	}
}

// String returns the constraint as a string.
func (c *VersionConstraint) String() string {
	return c.Operator + c.Version.String()
}

