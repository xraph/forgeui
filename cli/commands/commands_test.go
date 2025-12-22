package commands

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/xraph/forgeui/cli"
	"github.com/xraph/forgeui/cli/util"
)

func TestInitCommand(t *testing.T) {
	cmd := InitCommand()
	
	if cmd.Name != "init" {
		t.Errorf("InitCommand().Name = %v, want %v", cmd.Name, "init")
	}
	
	if len(cmd.Flags) == 0 {
		t.Error("InitCommand() should have flags")
	}
}

func TestGenerateCommand(t *testing.T) {
	cmd := GenerateCommand()
	
	if cmd.Name != "generate" {
		t.Errorf("GenerateCommand().Name = %v, want %v", cmd.Name, "generate")
	}
	
	if len(cmd.Subcommands) == 0 {
		t.Error("GenerateCommand() should have subcommands")
	}
	
	// Check for component subcommand
	hasComponent := false
	for _, sub := range cmd.Subcommands {
		if sub.Name == "component" {
			hasComponent = true
			break
		}
	}
	
	if !hasComponent {
		t.Error("GenerateCommand() should have component subcommand")
	}
}

func TestDevCommand(t *testing.T) {
	cmd := DevCommand()
	
	if cmd.Name != "dev" {
		t.Errorf("DevCommand().Name = %v, want %v", cmd.Name, "dev")
	}
	
	if len(cmd.Flags) == 0 {
		t.Error("DevCommand() should have flags")
	}
}

func TestBuildCommand(t *testing.T) {
	cmd := BuildCommand()
	
	if cmd.Name != "build" {
		t.Errorf("BuildCommand().Name = %v, want %v", cmd.Name, "build")
	}
	
	if len(cmd.Flags) == 0 {
		t.Error("BuildCommand() should have flags")
	}
}

func TestPluginCommand(t *testing.T) {
	cmd := PluginCommand()
	
	if cmd.Name != "plugin" {
		t.Errorf("PluginCommand().Name = %v, want %v", cmd.Name, "plugin")
	}
	
	if len(cmd.Subcommands) == 0 {
		t.Error("PluginCommand() should have subcommands")
	}
	
	// Check for expected subcommands
	expectedSubs := []string{"list", "add", "remove", "info"}
	for _, expected := range expectedSubs {
		found := false
		for _, sub := range cmd.Subcommands {
			if sub.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PluginCommand() missing subcommand: %s", expected)
		}
	}
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-project", false},
		{"valid with underscore", "my_project", false},
		{"valid alphanumeric", "myproject123", false},
		{"empty name", "", true},
		{"with spaces", "my project", true},
		{"starts with number", "123project", true},
		{"special chars", "my@project", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := util.ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"hello", "hello"},
		{"HTTPServer", "h_t_t_p_server"},
		{"myHTTPServer", "my_h_t_t_p_server"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := util.ToSnakeCase(tt.input); got != tt.want {
				t.Errorf("ToSnakeCase(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"hello", "Hello"},
		{"HELLO_WORLD", "HELLOWORLD"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := util.ToPascalCase(tt.input); got != tt.want {
				t.Errorf("ToPascalCase(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test", "nested", "dir")
	
	if err := util.CreateDir(testDir); err != nil {
		t.Errorf("CreateDir() error = %v", err)
	}
	
	if !util.DirExists(testDir) {
		t.Error("CreateDir() did not create directory")
	}
}

func TestCreateFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "test content"
	
	if err := util.CreateFile(testFile, content); err != nil {
		t.Errorf("CreateFile() error = %v", err)
	}
	
	if !util.FileExists(testFile) {
		t.Error("CreateFile() did not create file")
	}
	
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}
	
	if string(data) != content {
		t.Errorf("File content = %v, want %v", string(data), content)
	}
}

func TestConfig(t *testing.T) {
	config := cli.DefaultConfig()
	
	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	
	if config.Dev.Port == 0 {
		t.Error("DefaultConfig() should have default dev port")
	}
	
	if config.Build.OutputDir == "" {
		t.Error("DefaultConfig() should have default output dir")
	}
	
	// Test save and load
	tmpDir := t.TempDir()
	if err := config.Save(tmpDir); err != nil {
		t.Errorf("Config.Save() error = %v", err)
	}
	
	loaded, err := cli.LoadConfig(tmpDir)
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}
	
	if loaded.Name != config.Name {
		t.Errorf("Loaded config name = %v, want %v", loaded.Name, config.Name)
	}
}

