---
path: llm-context.md
page-type: overview
summary: Complete codebase summary optimized for LLM consumption.
tags: [llm, context, summary]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# LLM Context: Dracory CRUD

## Project Summary

The `github.com/dracory/crud/v2` package is a Go library that generates server-side rendered CRUD (Create, Read, Update, Trash) admin interfaces for web applications. It produces Bootstrap 5 HTML pages with Vue.js 3 interactivity. The package is callback-driven: consumers provide functions for data access (fetch rows, create, update, trash) and the package handles routing, HTML rendering, form generation, validation, and frontend interactivity. It is database-agnostic and designed to be embedded into existing Go web applications.

The repository contains two versions: v2 (current, recommended) and v1 (legacy, unmaintained). Each version is a separate Go module.

## Key Technologies

- **Go 1.25+** - Primary language
- **Bootstrap 5.3.3** - CSS framework (CDN)
- **Vue.js 3** - Client-side reactivity (CDN)
- **jQuery 3.7.1** - DOM manipulation and AJAX (CDN)
- **HTMX 2.0.0** - HTML-over-the-wire for create modal (CDN)
- **Sweetalert2 11** - Alert dialogs (CDN)
- **jQuery DataTables 1.13.4** - Sortable entity tables (CDN)
- **Element Plus 2.3.8** - Vue 3 date picker (CDN)
- **Trumbowyg 2.27.3** - WYSIWYG HTML editor (CDN)
- **github.com/dracory/hb** - Server-side HTML builder
- **github.com/dracory/form** - Form field interface
- **github.com/dracory/api** - JSON API response helpers
- **github.com/samber/lo** - Functional programming helpers

## Directory Structure

```
crud/
├── docs/              # Documentation (Markdown + LiveWiki)
├── v1/                # Legacy version (separate Go module)
│   ├── go.mod         # module github.com/dracory/crud
│   └── crud.go        # Monolithic implementation
├── v2/                # Current version (separate Go module)
│   ├── go.mod         # module github.com/dracory/crud/v2
│   ├── config.go      # Config struct (public API surface)
│   ├── crud.go        # Crud struct, Handler, routing, layout, forms, URLs
│   ├── new.go         # Constructor with validation
│   ├── constants.go   # Route paths + form field type constants
│   ├── log.go         # Log levels and helper
│   ├── pagination.go  # Pagination rendering
│   ├── breadcrumb.go  # Breadcrumb type
│   ├── key_value.go   # KeyValue type (read view)
│   ├── row.go         # Row type (entity manager table)
│   ├── entity_create_controller.go   # Create modal + AJAX save
│   ├── entity_manager_controller.go  # Entity listing page
│   ├── entity_read_controller.go     # Read-only view
│   ├── entity_update_controller.go   # Edit form + AJAX save
│   ├── entity_trash_controller.go    # Soft-delete via AJAX
│   └── *_test.go      # 76 tests total
└── .github/workflows/ # CI configuration
```

## Core Concepts

1. **Crud struct** - Central struct holding all configuration, routing, rendering, and URL generation. Created via `New(Config)`.
2. **Config struct** - Public configuration passed to `New()`. Contains entity names, endpoint, callback functions, form fields, and optional features.
3. **Handler method** - Single `http.HandlerFunc` entry point. Dispatches to internal controllers based on `path` query parameter.
4. **Callback functions** - All data access is via user-provided functions (FuncRows, FuncCreate, FuncUpdate, FuncTrash, etc.). The package never touches a database.
5. **Controllers** - Internal (unexported) structs that handle specific CRUD operations. Each holds a pointer to the parent Crud.
6. **Form fields** - Use `form.FieldInterface` from `github.com/dracory/form`. 11 field types supported (string, number, password, textarea, select, image, image_inline, htmlarea, blockarea, datetime, raw).
7. **Route dispatch** - Routes are `path` query parameter values: `entity-manager`, `entity-create-modal`, `entity-create-ajax`, `entity-read`, `entity-update`, `entity-update-ajax`, `entity-trash-ajax`.
8. **Middleware pipeline** - FuncLog → FuncBeforeAction → FuncValidateCSRF (POST only) → Controller → FuncAfterAction.
9. **Layout system** - Two modes: custom (FuncLayout wraps content in app shell) or default (standalone Bootstrap page).
10. **Server-side pagination** - When PageSize > 0 and FuncRowsCount is set, Bootstrap pagination controls are rendered.

## Common Patterns

- **XSS prevention** - All values interpolated into inline JS are JSON-encoded via `json.Marshal`
- **HTTP method enforcement** - State-mutating endpoints reject non-POST requests
- **Nil safety** - All optional callbacks are nil-checked before invocation
- **Raw HTML support** - Column names and KeyValue pairs wrapped in `{!! !!}` are rendered as raw HTML
- **Query string handling** - URL helpers detect existing `?` in endpoint and use `&` accordingly
- **Controller factory** - Each controller is created via a factory method on Crud (e.g., `newEntityCreateController()`)
- **Vue.js mounting** - Each page mounts a Vue app on a specific container (`#entity-manager`, `#entity-update`)

## Important Files

| File | Description |
|------|-------------|
| `v2/config.go` | **Config struct** - The public API surface. All configuration fields and callback signatures. |
| `v2/crud.go` | **Crud struct** - Core implementation: Handler, routing, layout, form rendering, URL helpers, breadcrumbs. |
| `v2/new.go` | **Constructor** - `New(Config)` with validation rules. |
| `v2/constants.go` | **Constants** - Route path strings and 11 form field type constants. |
| `v2/entity_manager_controller.go` | **Entity Manager** - Main listing page with DataTable, Vue.js app, pagination. |
| `v2/entity_create_controller.go` | **Create** - HTMX modal + AJAX save with required field validation. |
| `v2/entity_update_controller.go` | **Update** - Vue.js reactive form + AJAX save with Trumbowyg/Element Plus integration. |
| `v2/entity_read_controller.go` | **Read** - Read-only KeyValue table with raw HTML support and extras. |
| `v2/entity_trash_controller.go` | **Trash** - POST-only soft-delete endpoint + confirmation modal HTML. |
| `v2/pagination.go` | **Pagination** - Bootstrap pagination control rendering. |
| `v2/log.go` | **Logging** - Log level constants and nil-safe log helper. |

## See Also

- [Overview](overview.md) - Human-readable project overview
- [Architecture](architecture.md) - System design and patterns
- [API Reference](api_reference.md) - Complete API documentation
