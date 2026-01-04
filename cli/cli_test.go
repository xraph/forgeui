package cli

import (
	"bytes"
	"testing"
)

func TestCommand_Execute(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *Command
		args    []string
		wantErr bool
	}{
		{
			name: "simple command",
			cmd: &Command{
				Name:  "test",
				Short: "Test command",
				Run: func(ctx *Context) error {
					ctx.Println("test executed")
					return nil
				},
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "command with flag",
			cmd: &Command{
				Name:  "test",
				Short: "Test command",
				Flags: []Flag{
					StringFlag("name", "n", "Name flag", "default"),
				},
				Run: func(ctx *Context) error {
					name := ctx.GetString("name")
					ctx.Printf("name: %s\n", name)

					return nil
				},
			},
			args:    []string{"--name", "test"},
			wantErr: false,
		},
		{
			name: "command with subcommand",
			cmd: &Command{
				Name:  "test",
				Short: "Test command",
				Subcommands: []*Command{
					{
						Name:  "sub",
						Short: "Subcommand",
						Run: func(ctx *Context) error {
							ctx.Println("subcommand executed")
							return nil
						},
					},
				},
			},
			args:    []string{"sub"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Execute(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContext_Flags(t *testing.T) {
	ctx := &Context{
		Flags: map[string]any{
			"string": "test",
			"bool":   true,
			"int":    42,
		},
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}

	if got := ctx.GetString("string"); got != "test" {
		t.Errorf("GetString() = %v, want %v", got, "test")
	}

	if got := ctx.GetBool("bool"); got != true {
		t.Errorf("GetBool() = %v, want %v", got, true)
	}

	if got := ctx.GetInt("int"); got != 42 {
		t.Errorf("GetInt() = %v, want %v", got, 42)
	}

	// Test missing flags
	if got := ctx.GetString("missing"); got != "" {
		t.Errorf("GetString(missing) = %v, want empty string", got)
	}

	if got := ctx.GetBool("missing"); got != false {
		t.Errorf("GetBool(missing) = %v, want false", got)
	}

	if got := ctx.GetInt("missing"); got != 0 {
		t.Errorf("GetInt(missing) = %v, want 0", got)
	}
}

func TestFlags(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want FlagType
	}{
		{"string flag", StringFlag("test", "t", "Test flag", "default"), FlagTypeString},
		{"bool flag", BoolFlag("test", "t", "Test flag"), FlagTypeBool},
		{"int flag", IntFlag("test", "t", "Test flag", 10), FlagTypeInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.flag.Type != tt.want {
				t.Errorf("Flag.Type = %v, want %v", tt.flag.Type, tt.want)
			}
		})
	}
}

func TestRegisterCommand(t *testing.T) {
	// Save and restore root command
	oldRoot := rootCmd

	defer func() { rootCmd = oldRoot }()

	rootCmd = &Command{
		Name:        "test-root",
		Subcommands: []*Command{},
	}

	testCmd := &Command{
		Name:  "test",
		Short: "Test command",
	}

	RegisterCommand(testCmd)

	if len(rootCmd.Subcommands) != 1 {
		t.Errorf("RegisterCommand() failed to add command")
	}

	if rootCmd.Subcommands[0].Name != "test" {
		t.Errorf("RegisterCommand() added wrong command")
	}
}


