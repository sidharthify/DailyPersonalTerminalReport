# Development Guide

DPTR is built in pure Go and designed to be easily extensible. All modules are fetched concurrently via goroutines, meaning adding a slow module won't slow down the rest of the report.

## Directory Structure

- `cmd/dptr/` — Entry point (`main.go`).
- `internal/config/` — YAML struct definitions.
- `internal/wakeup/` — Wake-up guard logic and state management.
- `internal/renderer/` — ANSI terminal renderer.
- `internal/runner/` — Module orchestrator and parallel execution engine.
- `modules/` — All module implementations organized by category.

---

## Creating a New Module

Adding a module is incredibly straightforward.

### 1. Implement the Interface

A module is any Go struct that implements the `runner.Module` interface:
```go
type Module interface {
    GetData(cfg map[string]any) []string
}
```

Create a new file in the appropriate directory (e.g., `modules/custom/my_module.go`):

```go
package custom

import "fmt"

type MyModule struct{}

func (MyModule) GetData(cfg map[string]any) []string {
    // Read config
    name := "World"
    if v, ok := cfg["name"].(string); ok {
        name = v
    }

    // Fetch data via HTTP, DB, etc.
    // Return a slice of strings (one per line rendered to the terminal)
    return []string{
        fmt.Sprintf("Hello, %s!", name),
        "Data fetched successfully.",
    }
}
```

### 2. Register the Module

Open `internal/runner/runner.go` and add your module to the `moduleFactory` switch statement:

```go
func moduleFactory(name string) Module {
    switch name {
    // ... existing modules ...
    case "my_module":
        return &custom.MyModule{}
    default:
        return nil
    }
}
```

### 3. Add to config.yaml

```yaml
  - module: "my_module"
    title: "MY CUSTOM SECTION"
    config:
      name: "Sidharth"
```

---

## Testing Your Changes

Use the `--force` flag during development to completely bypass the wake-up guard time checks and render the output immediately inline:

```bash
go build ./cmd/dptr && ./dptr --force
```

If you modify configuration structs, be sure to update `internal/config/config.go`.
