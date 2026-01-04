# Changelog

All notable changes to ForgeUI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.2] - 2026-01-04

### Changed
- Restructured application configuration and page registration API
- Improved builder pattern for app configuration
- Enhanced type safety in page registration

## [0.0.1] - 2025-12-XX

### Added
- Initial release preparation
- CI/CD pipeline with GitHub Actions
- Automated releases with GoReleaser
- Multi-platform binary distribution
- Homebrew tap support

## [0.1.0] - TBD

### Added

#### Core Framework
- SSR-first rendering with gomponents
- Type-safe component API with functional options
- CVA (Class Variance Authority) for variant management
- Tailwind CSS integration with built-in processing
- 35+ production-ready UI components

#### Frontend Integration
- Alpine.js integration (directives, stores, magic helpers, plugins)
- HTMX support with complete attribute helpers
- 1600+ Lucide icons with customization
- Animation system (Tailwind animations, transitions, keyframes)

#### Backend Features
- Production-ready HTTP router with middleware support
- Bridge system for Go-JavaScript RPC
- Extensible plugin system with dependency management
- Customizable theme system with CSS variables

#### Developer Tools
- Assets pipeline (esbuild, Tailwind CSS, fingerprinting)
- Hot-reload development server
- CLI tools for project scaffolding
- Layout helpers with page builder

#### Components
- Accordion, Alert, Avatar, Badge, Breadcrumb
- Button (with variants and loading states)
- Card (with header, content, footer)
- Checkbox, Dropdown, Empty State
- Form (with validation and field components)
- Input (text, email, password, number, date, time, file)
- Label, List, Menu, Modal
- Navbar, Pagination, Popover, Progress
- Radio, Select, Separator, Sidebar
- Skeleton, Slider, Spinner, Switch
- Table (with sorting, filtering, pagination)
- Tabs, Textarea, Toast, Tooltip

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- N/A (initial release)

---

## Release Process

Releases are automated through GitHub Actions:

1. **Tag-based releases**: Push a version tag (e.g., `v0.1.0`)
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

2. **Manual releases**: Use GitHub Actions workflow dispatch

All releases include:
- Multi-platform binaries (Linux, macOS, Windows)
- Automated changelog generation
- SHA256 checksums
- Homebrew formula updates

---

[Unreleased]: https://github.com/xraph/forgeui/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/xraph/forgeui/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/xraph/forgeui/releases/tag/v0.0.1
[0.1.0]: https://github.com/xraph/forgeui/releases/tag/v0.1.0

