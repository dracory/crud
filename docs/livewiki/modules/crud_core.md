---
path: modules/crud_core.md
page-type: module
summary: Documentation for the Crud struct - the central type handling routing, layout, form rendering, URL helpers, and breadcrumbs.
tags: [module, crud, core, routing, handler]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: Crud Core

## Purpose

The `Crud` struct is the central type in the package. It holds all configuration, provides the HTTP handler entry point, dispatches requests to controllers, renders layouts and forms, generates URLs, and renders breadcrumbs.

**File:** `v2/crud.go`

## Responsibilities

- HTTP request handling and route dispatch
- Layout rendering (custom or default)
- Form generation from field definitions
- URL generation for all endpoints
- Breadcrumb rendering
- Middleware pipeline (before/after action, CSRF, logging)

## Key Types

### Crud Struct

```go
type Crud struct {
    columnNames         []string
    createFields        []form.FieldInterface
    endpoint            string
    entityNamePlural    string
    entityNameSingular  string
    fileManagerURL      string
    funcCreate          func(r *http.Request, data map[string]string) (string, error)
    funcReadExtras      func(r *http.Request, entityID string) []hb.TagInterface
    funcFetchReadData   func(r *http.Request, entityID string) ([]KeyValue, error)
    funcFetchUpdateData func(r *http.Request, entityID string) (map[string]string, error)
    funcLayout          func(w http.ResponseWriter, r *http.Request, ...) string
    funcRows            func(r *http.Request) ([]Row, error)
    funcTrash           func(r *http.Request, entityID string) error
    funcUpdate          func(r *http.Request, entityID string, data map[string]string) error
    homeURL             string
    createRedirectURL   string
    updateRedirectURL   string
    updateFields        []form.FieldInterface
    pageSize            int
    funcRowsCount       func(r *http.Request) (int64, error)
    funcBeforeAction    func(w http.ResponseWriter, r *http.Request, action string) bool
    funcAfterAction     func(w http.ResponseWriter, r *http.Request, action string)
    funcValidateCSRF    func(r *http.Request) error
    funcLog             func(level string, message string, attrs map[string]any)
}
```

All fields are **unexported**. The struct is created via `New(Config)` and interacted with via public methods.

## Public Methods

### Handler

```go
func (crud Crud) Handler(w http.ResponseWriter, r *http.Request)
```

The main HTTP handler. Executes the middleware pipeline and dispatches to the appropriate controller.

**Pipeline:**
1. Extract `path` query parameter (default: `"home"`)
2. Log the request via `FuncLog` (if configured)
3. Call `FuncBeforeAction` (if configured; abort if returns `false`)
4. Validate CSRF on POST requests via `FuncValidateCSRF` (if configured)
5. Dispatch to controller via `getRoute(path)`
6. Call `FuncAfterAction` (if configured)

### URL Helpers

| Method | Returns |
|--------|---------|
| `UrlEntityManager() string` | Entity manager page URL |
| `UrlEntityCreateModal() string` | Create modal fragment URL |
| `UrlEntityCreateAjax() string` | Create AJAX endpoint URL |
| `UrlEntityRead() string` | Read page base URL |
| `UrlEntityUpdate() string` | Update page base URL |
| `UrlEntityUpdateAjax() string` | Update AJAX endpoint URL |
| `UrlEntityTrashAjax() string` | Trash AJAX endpoint URL |

All URL helpers detect whether the endpoint already contains `?` and use `&` accordingly.

## Internal Methods

### getRoute

```go
func (crud *Crud) getRoute(route string) func(w http.ResponseWriter, r *http.Request)
```

Maps route path strings to controller handler functions. Unknown routes fall back to the entity manager.

### layout

```go
func (crud *Crud) layout(w, r, title, content, styleFiles, style, jsFiles, js) string
```

Renders the final HTML. Delegates to `FuncLayout` (with prepended CDN resources) or falls back to the default `webpage()` template.

### webpage

```go
func (crud *Crud) webpage(title, content string) *hb.HtmlWebpage
```

Generates a standalone HTML page with Bootstrap 5, jQuery, Vue.js, and Sweetalert2 from CDN.

### form

```go
func (crud *Crud) form(fields []form.FieldInterface) []hb.TagInterface
```

Generates Bootstrap form groups from field definitions. Handles all 11 field types.

### renderBreadcrumbs

```go
func (crud *Crud) renderBreadcrumbs(breadcrumbs []Breadcrumb) string
```

Generates Bootstrap breadcrumb navigation HTML.

### listCreateNames / listUpdateNames

```go
func (crud *Crud) listCreateNames() []string
func (crud *Crud) listUpdateNames() []string
```

Extract field names from create/update field slices, skipping fields with empty names.

## Dependencies

- `github.com/dracory/api` - JSON responses
- `github.com/dracory/bs` - Bootstrap components
- `github.com/dracory/cdn` - CDN URLs
- `github.com/dracory/form` - Form field interface
- `github.com/dracory/hb` - HTML builder
- `github.com/dracory/req` - Request parameter extraction
- `github.com/dracory/str` - Random ID generation
- `github.com/samber/lo` - Functional helpers

## See Also

- [Modules: Config](config.md) - Config struct details
- [Modules: Controllers](controllers.md) - Controller implementations
- [Architecture](../architecture.md) - System design
- [API Reference](../api_reference.md) - Public API
