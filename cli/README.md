# ForgeUI CLI

> Production-ready command-line interface for ForgeUI applications

The ForgeUI CLI provides a complete toolchain for creating, developing, and building ForgeUI applications. It features project scaffolding, code generation, development servers, production builds, and plugin management.

---

## Features

- **Project Initialization**: Bootstrap new projects with multiple templates
- **Code Generation**: Generate components and pages with customizable templates
- **Development Server**: Hot reload development environment
- **Production Builds**: Optimized builds with asset processing
- **Plugin Management**: Easy plugin installation and configuration
- **Zero External Dependencies**: Built entirely with Go standard library

---

## Installation

### From Source

```bash
go install github.com/xraph/forgeui/cmd/forgeui@latest
```

### Build Locally

```bash
cd cmd/forgeui
go build -o forgeui
```

---

## Quick Start

### Create a New Project

```bash
# Interactive mode
forgeui init

# With options
forgeui init my-app --template=standard --module=github.com/user/my-app
```

**Available Templates:**
- `minimal` - Basic setup with one page
- `standard` - Full setup with router, assets, examples
- `blog` - Blog template with posts and tags
- `dashboard` - Admin dashboard with charts and tables
- `api` - API-first template with HTMX

### Start Development Server

```bash
cd my-app
forgeui dev

# Custom port
forgeui dev --port 8080

# Open browser automatically
forgeui dev --open
```

### Generate Code

```bash
# Generate a component
forgeui generate component Button
forgeui g c Button --type=compound --with-props --with-test

# Generate a page
forgeui generate page About --path=/about
forgeui g p UserProfile --type=detail --with-loader --with-meta
```

### Build for Production

```bash
# Build assets
forgeui build

# Build with binary
forgeui build --binary --minify --embed

# Custom output directory
forgeui build --output=dist
```

### Manage Plugins

```bash
# List installed plugins
forgeui plugin list

# Add a plugin
forgeui plugin add toast

# Remove a plugin
forgeui plugin remove toast

# Get plugin information
forgeui plugin info charts
```

---

## Commands Reference

### `forgeui init [project-name]`

Initialize a new ForgeUI project.

**Flags:**
- `--template, -t` - Project template (default: minimal)
- `--module, -m` - Go module path
- `--force, -f` - Force initialization even if directory exists

**Examples:**
```bash
forgeui init my-app
forgeui init my-app --template=dashboard
forgeui init my-app --module=github.com/user/my-app
```

---

### `forgeui generate <type> <name>`

Generate code from templates.

**Aliases:** `g`

#### Component Generation

```bash
forgeui generate component <name>
forgeui g c <name>
```

**Flags:**
- `--type, -t` - Component type (basic, compound, form, layout, data)
- `--dir, -d` - Output directory (default: components)
- `--with-variants` - Add CVA variants
- `--with-props` - Generate props struct
- `--with-test` - Generate test file

**Component Types:**
- `basic` - Simple functional component
- `compound` - Compound component (Card, Modal)
- `form` - Form component with validation
- `layout` - Layout component
- `data` - Data display component (Table, List)

**Examples:**
```bash
forgeui g c Button
forgeui g c ProfileCard --type=compound --with-props
forgeui g c ContactForm --type=form --with-test
```

#### Page Generation

```bash
forgeui generate page <name>
forgeui g p <name>
```

**Flags:**
- `--type, -t` - Page type (simple, dynamic, form, list, detail)
- `--path` - Route path
- `--dir, -d` - Output directory (default: pages)
- `--with-loader` - Add data loader
- `--with-meta` - Add SEO meta tags

**Page Types:**
- `simple` - Static page
- `dynamic` - Dynamic page with data
- `form` - Page with form submission
- `list` - List page with pagination
- `detail` - Detail page with params

**Examples:**
```bash
forgeui g p About --path=/about
forgeui g p UserList --type=list --with-loader
forgeui g p UserDetail --type=detail --path=/users/:id
```

---

### `forgeui dev`

Start development server with hot reload.

**Flags:**
- `--port, -p` - Port to listen on (default: 3000)
- `--host, -h` - Host to bind to (default: localhost)
- `--open, -o` - Open browser automatically

**Examples:**
```bash
forgeui dev
forgeui dev --port 8080 --open
```

**Features:**
- Automatic Go application restart
- Asset watching and rebuilding
- Live browser reload
- Graceful shutdown

---

### `forgeui build`

Build for production deployment.

**Flags:**
- `--output, -o` - Output directory (default: dist)
- `--binary, -b` - Compile Go binary
- `--minify, -m` - Minify assets
- `--embed, -e` - Embed assets in binary

**Examples:**
```bash
forgeui build
forgeui build --binary --minify --embed
forgeui build --output=build
```

**Build Process:**
1. Creates output directory
2. Copies static assets
3. Processes CSS and JS (if configured)
4. Generates fingerprinted files
5. Creates manifest file
6. Optionally compiles Go binary

---

### `forgeui plugin <command>`

Manage ForgeUI plugins.

#### `forgeui plugin list`

List installed plugins.

```bash
forgeui plugin list
```

#### `forgeui plugin add <name>`

Add a plugin to the project.

```bash
forgeui plugin add toast
forgeui plugin add charts
```

**Available Plugins:**
- `toast` - Notification system
- `sortable` - Drag-and-drop sorting
- `charts` - Data visualization (Line, Bar, Pie, Area, Doughnut)
- `analytics` - Analytics tracking
- `seo` - SEO optimization tools
- `htmxplugin` - HTMX integration wrapper

#### `forgeui plugin remove <name>`

Remove a plugin from the project.

**Aliases:** `rm`

```bash
forgeui plugin remove toast
forgeui plugin rm charts
```

#### `forgeui plugin info <name>`

Show detailed plugin information.

```bash
forgeui plugin info toast
forgeui plugin info charts
```

---

## Configuration

The CLI uses a `.forgeui.json` configuration file in the project root.

### Configuration File Structure

```json
{
  "name": "my-app",
  "version": "1.0.0",
  "dev": {
    "port": 3000,
    "host": "localhost",
    "auto_reload": true,
    "open_browser": false
  },
  "build": {
    "output_dir": "dist",
    "public_dir": "public",
    "minify": true,
    "binary": false,
    "embed_assets": true
  },
  "assets": {
    "css": ["public/css/app.css"],
    "js": ["public/js/app.js"]
  },
  "plugins": ["toast", "charts"],
  "router": {
    "base_path": "/",
    "not_found": "pages/404.go"
  }
}
```

### Configuration Options

**Dev Configuration:**
- `port` - Development server port
- `host` - Host to bind to
- `auto_reload` - Enable automatic reload
- `open_browser` - Open browser on start

**Build Configuration:**
- `output_dir` - Output directory for builds
- `public_dir` - Source directory for static assets
- `minify` - Minify CSS and JS
- `binary` - Compile Go binary
- `embed_assets` - Embed assets in binary

**Assets Configuration:**
- `css` - CSS files to process
- `js` - JavaScript files to process

**Plugins Configuration:**
- `plugins` - Array of installed plugin names

**Router Configuration:**
- `base_path` - Base path for routes
- `not_found` - Custom 404 handler

---

## Project Structure

A typical ForgeUI project structure:

```
my-app/
├── main.go              # Application entry point
├── go.mod               # Go module file
├── .forgeui.json        # CLI configuration
├── .gitignore           # Git ignore file
├── README.md            # Project documentation
├── components/          # Reusable components
│   ├── button/
│   │   ├── button.go
│   │   └── button_test.go
│   └── card/
│       └── card.go
├── pages/               # Page handlers
│   ├── home.go
│   ├── about.go
│   └── contact.go
└── public/              # Static assets
    ├── css/
    │   └── app.css
    └── js/
        └── app.js
```

---

## Architecture

The CLI is built with a custom lightweight command framework (no external dependencies):

```
cli/
├── cli.go              # Main executor
├── command.go          # Command type
├── context.go          # Execution context
├── flags.go            # Flag parsing
├── config.go           # Configuration
├── commands/           # Command implementations
│   ├── init.go
│   ├── generate.go
│   ├── dev.go
│   ├── build.go
│   └── plugin.go
├── templates/          # Code templates
│   ├── minimal.go
│   ├── standard.go
│   ├── component_types.go
│   └── page_types.go
└── util/               # Utilities
    ├── fs.go
    ├── prompt.go
    ├── spinner.go
    └── color.go
```

---

## Development

### Running Tests

```bash
# Run all tests
go test ./cli/...

# Run with coverage
go test -cover ./cli/...

# Run integration tests
go test -v ./cli/integration_test.go
```

### Building the CLI

```bash
cd cmd/forgeui
go build -o forgeui
```

### Adding a New Command

1. Create command file in `cli/commands/`
2. Implement command using `cli.Command` struct
3. Register command with `cli.RegisterCommand()` in `init()`
4. Add tests in `commands_test.go`

**Example:**

```go
package commands

import "github.com/xraph/forgeui/cli"

func init() {
    cli.RegisterCommand(MyCommand())
}

func MyCommand() *cli.Command {
    return &cli.Command{
        Name:  "mycommand",
        Short: "My custom command",
        Run:   runMyCommand,
    }
}

func runMyCommand(ctx *cli.Context) error {
    ctx.Println("Hello from my command!")
    return nil
}
```

---

## Best Practices

### Project Initialization

1. Use descriptive project names
2. Choose the appropriate template for your use case
3. Use proper Go module paths (e.g., `github.com/user/project`)

### Code Generation

1. Use meaningful component and page names
2. Generate tests for complex components
3. Use props structs for reusable components
4. Follow naming conventions (PascalCase for components)

### Development

1. Use `forgeui dev` for development
2. Keep the dev server running for hot reload
3. Test changes in the browser immediately
4. Use configuration file for project-specific settings

### Production Builds

1. Always test builds before deployment
2. Use `--minify` for smaller assets
3. Use `--embed` for single-binary deployments
4. Verify manifest file is generated correctly

### Plugin Management

1. Only install plugins you need
2. Check plugin info before installing
3. Keep plugins updated
4. Remove unused plugins

---

## Troubleshooting

### Command Not Found

Ensure the CLI is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Port Already in Use

Use a different port:

```bash
forgeui dev --port 8080
```

### Build Fails

1. Check that you're in a Go project directory
2. Verify `go.mod` exists
3. Run `go mod tidy`
4. Check for compilation errors

### Plugin Installation Fails

1. Ensure you have network access
2. Check plugin name is correct
3. Verify Go module can be downloaded
4. Run `go mod tidy` after installation

---

## Examples

### Create a Blog

```bash
forgeui init my-blog --template=blog
cd my-blog
forgeui plugin add seo
forgeui g c PostCard --type=data
forgeui g p PostDetail --type=detail --path=/post/:slug
forgeui dev
```

### Create a Dashboard

```bash
forgeui init admin-dashboard --template=dashboard
cd admin-dashboard
forgeui plugin add charts
forgeui plugin add analytics
forgeui g c StatsCard --with-props
forgeui g p Analytics --type=dynamic --with-loader
forgeui dev --open
```

### Create an API

```bash
forgeui init my-api --template=api
cd my-api
forgeui plugin add htmxplugin
forgeui g p Users --type=list
forgeui g p UserDetail --type=detail --path=/users/:id
forgeui dev
```

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Write tests for new commands
2. Follow Go coding standards
3. Update documentation
4. Add examples for new features

---

## License

ForgeUI CLI is part of the ForgeUI project.

---

## Support

- **Documentation**: [github.com/xraph/forgeui](https://github.com/xraph/forgeui)
- **Issues**: [GitHub Issues](https://github.com/xraph/forgeui/issues)
- **Discussions**: [GitHub Discussions](https://github.com/xraph/forgeui/discussions)

---

**Built with ❤️ by the ForgeUI team**



