package plugin

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "basic version",
			input: "1.2.3",
			want:  &Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:  "with v prefix",
			input: "v1.2.3",
			want:  &Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:  "zeros",
			input: "0.0.0",
			want:  &Version{Major: 0, Minor: 0, Patch: 0},
		},
		{
			name:    "invalid format",
			input:   "1.2",
			wantErr: true,
		},
		{
			name:    "non-numeric",
			input:   "1.2.x",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch {
					t.Errorf("ParseVersion() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		name string
		v1   *Version
		v2   *Version
		want int
	}{
		{
			name: "equal",
			v1:   &Version{1, 2, 3},
			v2:   &Version{1, 2, 3},
			want: 0,
		},
		{
			name: "major greater",
			v1:   &Version{2, 0, 0},
			v2:   &Version{1, 9, 9},
			want: 1,
		},
		{
			name: "major less",
			v1:   &Version{1, 0, 0},
			v2:   &Version{2, 0, 0},
			want: -1,
		},
		{
			name: "minor greater",
			v1:   &Version{1, 3, 0},
			v2:   &Version{1, 2, 9},
			want: 1,
		},
		{
			name: "minor less",
			v1:   &Version{1, 2, 0},
			v2:   &Version{1, 3, 0},
			want: -1,
		},
		{
			name: "patch greater",
			v1:   &Version{1, 2, 4},
			v2:   &Version{1, 2, 3},
			want: 1,
		},
		{
			name: "patch less",
			v1:   &Version{1, 2, 3},
			v2:   &Version{1, 2, 4},
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v1.Compare(tt.v2)
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	v := &Version{Major: 1, Minor: 2, Patch: 3}

	want := "1.2.3"
	if got := v.String(); got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}

func TestParseConstraint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *VersionConstraint
		wantErr bool
	}{
		{
			name:  "exact",
			input: "=1.2.3",
			want:  &VersionConstraint{Operator: "=", Version: &Version{1, 2, 3}},
		},
		{
			name:  "no operator defaults to exact",
			input: "1.2.3",
			want:  &VersionConstraint{Operator: "=", Version: &Version{1, 2, 3}},
		},
		{
			name:  "greater than",
			input: ">1.2.3",
			want:  &VersionConstraint{Operator: ">", Version: &Version{1, 2, 3}},
		},
		{
			name:  "greater or equal",
			input: ">=1.2.3",
			want:  &VersionConstraint{Operator: ">=", Version: &Version{1, 2, 3}},
		},
		{
			name:  "less than",
			input: "<1.2.3",
			want:  &VersionConstraint{Operator: "<", Version: &Version{1, 2, 3}},
		},
		{
			name:  "less or equal",
			input: "<=1.2.3",
			want:  &VersionConstraint{Operator: "<=", Version: &Version{1, 2, 3}},
		},
		{
			name:  "tilde",
			input: "~1.2.3",
			want:  &VersionConstraint{Operator: "~", Version: &Version{1, 2, 3}},
		},
		{
			name:  "caret",
			input: "^1.2.3",
			want:  &VersionConstraint{Operator: "^", Version: &Version{1, 2, 3}},
		},
		{
			name:    "invalid version",
			input:   ">=invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConstraint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConstraint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Operator != tt.want.Operator {
					t.Errorf("ParseConstraint() operator = %v, want %v", got.Operator, tt.want.Operator)
				}

				if got.Version.Compare(tt.want.Version) != 0 {
					t.Errorf("ParseConstraint() version = %v, want %v", got.Version, tt.want.Version)
				}
			}
		})
	}
}

func TestConstraintCheck(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		version    string
		want       bool
	}{
		// Exact
		{"exact match", "=1.2.3", "1.2.3", true},
		{"exact mismatch", "=1.2.3", "1.2.4", false},

		// Greater than
		{">1.2.3 with 1.2.4", ">1.2.3", "1.2.4", true},
		{">1.2.3 with 1.2.3", ">1.2.3", "1.2.3", false},
		{">1.2.3 with 1.2.2", ">1.2.3", "1.2.2", false},

		// Greater or equal
		{">=1.2.3 with 1.2.4", ">=1.2.3", "1.2.4", true},
		{">=1.2.3 with 1.2.3", ">=1.2.3", "1.2.3", true},
		{">=1.2.3 with 1.2.2", ">=1.2.3", "1.2.2", false},

		// Less than
		{"<1.2.3 with 1.2.2", "<1.2.3", "1.2.2", true},
		{"<1.2.3 with 1.2.3", "<1.2.3", "1.2.3", false},
		{"<1.2.3 with 1.2.4", "<1.2.3", "1.2.4", false},

		// Less or equal
		{"<=1.2.3 with 1.2.2", "<=1.2.3", "1.2.2", true},
		{"<=1.2.3 with 1.2.3", "<=1.2.3", "1.2.3", true},
		{"<=1.2.3 with 1.2.4", "<=1.2.3", "1.2.4", false},

		// Tilde (patch updates)
		{"~1.2.3 with 1.2.4", "~1.2.3", "1.2.4", true},
		{"~1.2.3 with 1.2.3", "~1.2.3", "1.2.3", true},
		{"~1.2.3 with 1.3.0", "~1.2.3", "1.3.0", false},
		{"~1.2.3 with 2.0.0", "~1.2.3", "2.0.0", false},

		// Caret (minor updates)
		{"^1.2.3 with 1.2.4", "^1.2.3", "1.2.4", true},
		{"^1.2.3 with 1.3.0", "^1.2.3", "1.3.0", true},
		{"^1.2.3 with 2.0.0", "^1.2.3", "2.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint, err := ParseConstraint(tt.constraint)
			if err != nil {
				t.Fatalf("ParseConstraint() error = %v", err)
			}

			version, err := ParseVersion(tt.version)
			if err != nil {
				t.Fatalf("ParseVersion() error = %v", err)
			}

			got := constraint.Check(version)
			if got != tt.want {
				t.Errorf("Check(%s) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
