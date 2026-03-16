# AGENTS.md - Guidelines for Working on Gowid

This document provides guidelines for agentic coding agents working on the gowid codebase.

## Project Overview

Gowid is a terminal user interface (TUI) library for Go, inspired by urwid. It provides widgets and a framework for building terminal UIs, built on top of the tcell package.

- **Module**: `github.com/gnuos/gowid` (published at `github.com/gcla/gowid`)
- **Go Version**: 1.25.0 (as of go.mod)
- **Test Framework**: Go's built-in `testing` package + `github.com/stretchr/testify`

---

## Build, Lint, and Test Commands

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run a single test file
go test -v ./canvas_test.go

# Run a specific test function
go test -v -run TestCanvas19 ./...

# Run tests in a specific package
go test -v ./widgets/button/...

# Run tests with coverage
go test -v -cover ./...
```

### Building

```bash
# Build all packages
go build -v ./...

# Build specific binary (e.g., examples)
go build -v ./examples/...
```

### Linting

The project does not have a `.golangci.yml` configuration file. Run standard Go tooling:

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run all checks (includes vet)
go build -v ./...
```

### Running Examples

Examples are in the `examples/` directory. After building, run them directly:

```bash
gowid-fib
gowid-fractal
# etc.
```

---

## Code Style Guidelines

### File Structure

- **License Header**: Every source file must include the MIT license header:
  ```go
  // Copyright 2019-2022 Graham Clark. All rights reserved.  Use of this source
  // code is governed by the MIT license that can be found in the LICENSE
  // file.
  ```

- **Package Declaration**: Immediately after the license header, followed by imports.

### Imports

Standard Go import ordering (blank line between groups):

```go
import (
    // Standard library
    "fmt"
    "os"
    "strings"

    // Third-party packages
    "github.com/gdamore/tcell/v3"
    log "github.com/sirupsen/logrus"

    // Internal packages (if applicable)
)
```

- Use aliases for packages when the package name is unclear (e.g., `log "github.com/sirupsen/logrus"`)

### Naming Conventions

- **Types**: PascalCase (e.g., `App`, `Canvas`, `IWidget`)
- **Interfaces**: Prefer `I` prefix for interface types (e.g., `IApp`, `IWidget`, `IComposite`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Constants**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase; avoid single letters except in short scopes (e.g., `i`, `j`, `ctx`)
- **Acronyms**: Keep uppercase (e.g., `URL`, `HTML`, `API` not `Url`, `Html`)

### Type Annotations

- Use explicit type annotations for struct fields:
  ```go
  type App struct {
      screen    tcell.Screen
      view      IWidget
      colorMode ColorMode
  }
  ```

- Interface implementations should be verified with compile-time checks:
  ```go
  var _ IApp = (*App)(nil)
  ```

### Error Handling

- Return errors as the last return value
- Use `pkg/errors` for error wrapping (the project depends on it)
- Check errors immediately; avoid ignoring them with `_`
- Use descriptive error messages with context

### Comments

- Comment exported types and functions with descriptive comments
- Use doc comments for public APIs
- Comments should be full sentences with proper punctuation
- Place comments above the code they describe

### Testing

- Test files are named `*_test.go` in the same package as the code they test
- Use table-driven tests when testing multiple cases:
  ```go
  tests := []struct {
      name     string
          input    string
          expected string
      }{
          {"case1", "abc", "abc"},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              // test code
          })
      }
  ```
- Use `github.com/stretchr/testify/assert` for assertions
- Use `t.Errorf` for descriptive failure messages

### Widget Development Patterns

- Implement the `IWidget` interface:
  ```go
  type IWidget interface {
      RenderSize() IWidgetDimension
      Render(app IApp, size IWidgetDimension, focus bool) ICanvas
      UserInput(app IApp, ev any, size IWidgetDimension, focus bool) bool
  }
  ```
- Implement `IIdentityWidget` if the widget needs to be identifiable for mouse interactions
- Implement container widgets with `IContainerWidget` interface
- Use composition over inheritance

### Canvas and Rendering

- Widgets return an `ICanvas` from their `Render` method
- Use `NewCanvas()` to create empty canvases
- Use `CellsFromString()` to create cell arrays from strings
- Merge canvases using `MergeUnder()` and `MergeOver()`

## Architecture Overview

### Core Components
- **IWidget** - The fundamental interface all widgets implement (defined in `support.go`)
- **IApp** - Application context providing screen access, palette, and main loop (in `app.go`)
- **ICanvas** - 2D array of cells returned by widget Render() methods (in `canvas.go`)

### App Lifecycle
1. Create widget hierarchy (composite widgets wrap child widgets)
2. Create App with `gowid.NewApp(gowid.AppArgs{View: ..., Palette: ...})`
3. Run main loop: `app.SimpleMainLoop()` or `app.MainLoop(handler)`

### Rendering Pipeline
1. App calls `RenderRoot()` which invokes root widget's `Render()`
2. Widget calculates sub-widget sizes, renders children, merges canvases
3. Canvas operations: `MergeUnder()`, `AppendBelow()`, `AppendRight()`
4. Final canvas drawn to terminal via `Draw()` function

### Key Modules and Packages

| Package | Purpose |
|---------|---------|
| `/` | Core types (App, Canvas, Widget interfaces) |
| `/widgets/button` | Button widget |
| `/widgets/text` | Text display widget |
| `/widgets/edit` | Text editor widget |
| `/widgets/columns` | Horizontal layout |
| `/widgets/pile` | Vertical layout |
| `/widgets/grid` | 2D grid layout |
| `/widgets/boxadapter` | Box model adapter |
| `/widgets/framed` | Frame/box around widgets |
| `/widgets/overlay` | Modal overlay |
| `/widgets/padding` | Add padding around widgets |
| `/widgets/styled` | Apply palette styling |
| `/widgets/vpadding` | Vertical padding |
| `/widgets/divider` | Horizontal divider |
| `/widgets/fill` | Fill with character |
| `/widgets/progress` | Progress bar |
| `/widgets/menu` | Menu widget |
| `/widgets/checkbox` | Checkbox input |
| `/widgets/radio` | Radio button |
| `/widgets/terminal` | VT220 terminal emulator |
| `/widgets/tree` | Tree widget |
| `/widgets/list` | Infinite list widget |
| `/widgets/table` | Table widget |
| `/gwtest` | Test utilities |

### Focus and Selection
- Widgets implement `Selectable()` to indicate if they can receive focus
- Container widgets implement `IFocus` for focus navigation between children
- Use `gowid.ContainerWidget` helper to wrap widgets with dimensions

### Mouse Interaction
- Implement `IIdentityWidget` for widget identification
- Implement `IClickTargets` for mouse click handling
- Use `IMouseCallbacks` for mouse event handlers

### Common Development Patterns

- Use functional options for configuration (see `AppArgs` and similar patterns)
- Use palette-based styling (`IPalette`, `MakePaletteRef`, `MakePaletteEntry`)
- Support color modes: `ColorModeMono`, `ColorMode16`, `ColorMode256`, `ColorModeRGB`
- Handle mouse input via `IClickTargets` and related interfaces
- Use `IContainerWidget` for widget composition

---

## Working with the Codebase

### Adding a New Widget

1. Create a new directory under `/widgets/`
2. Define your widget type implementing `IWidget`
3. Add tests in `<widget>_test.go`
4. Export the widget constructor and options types
5. Update docs if needed

### Running Specific Tests

```bash
# Single test function
go test -v -run TestCanvas19 ./...

# Tests matching a pattern
go test -v -run TestCanvas ./...

# Verbose with timing
go test -v -race ./...
```

### Local Development Tips

- Use `go mod tidy` after adding dependencies
- Test on multiple terminals with different color capabilities
- Use the examples in `/examples/` to test widget functionality visually

### Examples

The `examples/` directory contains runnable examples:
- `gowid-helloworld` - Basic "Hello World" app
- `gowid-fib` - Fibonacci calculator demonstrating interactive widgets
- `gowid-fractal` - Fractal renderer
- `gowid-table` - Table widget demo
- `gowid-edit` - Text editor demo
- `gowid-list` - Infinite list demo
- `gowid-tree` - Tree widget demo

Build and run with:
```bash
go build -v ./examples/...
```

---

## References

- [urwid](http://urwid.org) - Python library that inspired gowid
- [tcell](https://github.com/gdamore/tcell) - Terminal cell library
- [docs/Tutorial.md](docs/Tutorial.md) - Beginner's tutorial
- [docs/Widgets.md](docs/Widgets.md) - Widget documentation
- [docs/FAQ.md](docs/FAQ.md) - Frequently asked questions
