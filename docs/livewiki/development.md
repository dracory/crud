---
path: development.md
page-type: tutorial
summary: Development workflow, testing, contributing guidelines, and project structure for the Dracory CRUD package.
tags: [development, testing, contributing, workflow]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Development

## Project Structure

```
crud/
├── docs/              # Shared documentation (Markdown)
│   └── livewiki/      # LiveWiki documentation (this site)
├── v1/                # Legacy version (original implementation)
│   ├── README.md
│   ├── go.mod
│   └── ...
├── v2/                # Current version (recommended)
│   ├── README.md
│   ├── go.mod
│   └── ...
├── .github/
│   └── workflows/
│       └── tests.yml  # CI test workflow
├── .gitignore
└── README.md          # Root README
```

## Go Module Structure

The project uses a **multi-module** layout. Each version is a separate Go module:

- `v2/go.mod` → `module github.com/dracory/crud/v2`
- `v1/go.mod` → `module github.com/dracory/crud`

This allows independent versioning and dependency management per version.

## Running Tests

### v2 Tests

```bash
cd v2
go test ./... -v
```

The v2 test suite includes **76 tests** covering:

- Constructor validation (`new_test.go`)
- Entity manager controller (`entity_manager_controller_test.go`)
- Entity create controller (`entity_create_controller_test.go`)
- Entity read controller (`entity_read_controller_test.go`)
- Entity update controller (`entity_update_controller_test.go`)
- Entity trash controller (`entity_trash_controller_test.go`)
- Routing, URL helpers, breadcrumbs, layout rendering (`crud_test.go`)

### v1 Tests

```bash
cd v1
go test ./... -v
```

## Test Patterns

Tests follow standard Go testing patterns using `testing.T`:

```go
func TestNew_MissingFuncRows(t *testing.T) {
    _, err := New(Config{})
    if err == nil {
        t.Fatal("expected error when FuncRows is nil")
    }
    if err.Error() != "FuncRows function is required" {
        t.Fatalf("unexpected error message: %s", err.Error())
    }
}
```

Controller tests use `httptest.NewRecorder` and `httptest.NewRequest` to simulate HTTP requests:

```go
func TestEntityManagerController_Page(t *testing.T) {
    c, _ := New(Config{
        FuncRows: func(r *http.Request) ([]Row, error) {
            return []Row{{ID: "1", Data: []string{"1", "Test"}}}, nil
        },
        UpdateFields: []form.FieldInterface{},
        ColumnNames:  []string{"ID", "Name"},
    })
    w := httptest.NewRecorder()
    r := httptest.NewRequest("GET", "/admin?path=entity-manager", nil)
    c.Handler(w, r)
    // Assert response...
}
```

## CI/CD

The project uses GitHub Actions for continuous integration. The workflow is defined in `.github/workflows/tests.yml`.

## Dependencies

All dependencies are managed via Go modules. Key dependencies:

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/dracory/api` | v1.7.0 | JSON API response helpers |
| `github.com/dracory/bs` | v0.16.0 | Bootstrap component builder |
| `github.com/dracory/cdn` | v1.9.0 | CDN URL helpers |
| `github.com/dracory/form` | v0.19.0 | Form field interface and builder |
| `github.com/dracory/hb` | v1.88.0 | HTML builder |
| `github.com/dracory/req` | v0.1.0 | HTTP request parameter extraction |
| `github.com/dracory/str` | v0.17.0 | String utilities |
| `github.com/samber/lo` | v1.52.0 | Functional programming helpers |

## Adding a New Feature

1. **Plan** - Identify which files need changes (config, controller, tests)
2. **Config** - Add new fields to `Config` struct in `config.go`
3. **Constructor** - Update `New()` in `new.go` to copy new config fields
4. **Crud struct** - Add corresponding private fields to `Crud` in `crud.go`
5. **Implementation** - Implement the feature in the appropriate controller
6. **Tests** - Add tests covering the new functionality
7. **Documentation** - Update docs in `docs/` and `docs/livewiki/`

## Code Style

- Standard Go formatting (`gofmt`)
- No external linter configuration (use Go defaults)
- Private types for controllers (unexported)
- Public types for configuration and data structures
- Callback functions for all data access (no direct DB access)
- `json.Marshal` for all values interpolated into inline JavaScript

## See Also

- [Architecture](architecture.md) - System design and patterns
- [Conventions](conventions.md) - Coding conventions
- [Troubleshooting](troubleshooting.md) - Common issues and solutions
