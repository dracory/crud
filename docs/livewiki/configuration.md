---
path: configuration.md
page-type: reference
summary: Complete reference for the Config struct, validation rules, callback function signatures, and configuration options.
tags: [configuration, config, callbacks, validation, setup]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Configuration

## Creating a Crud Instance

Use the `New()` function to create a configured `Crud` instance:

```go
import crud "github.com/dracory/crud/v2"

c, err := crud.New(crud.Config{
    Endpoint:           "/admin",
    EntityNameSingular: "Product",
    EntityNamePlural:   "Products",
    HomeURL:            "/dashboard",
    ColumnNames:        []string{"ID", "Name", "Status"},
    FuncRows:           fetchAllProducts,
    FuncCreate:         createProduct,
    FuncFetchReadData:  fetchProductReadData,
    FuncFetchUpdateData: fetchProductUpdateData,
    FuncUpdate:         updateProduct,
    FuncTrash:          trashProduct,
    CreateFields:       createFields,
    UpdateFields:       updateFields,
})
if err != nil {
    log.Fatal(err)
}

http.HandleFunc("/admin", c.Handler)
```

## Config Struct Reference

### Core Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Endpoint` | `string` | No | Base URL path for the CRUD interface. Used to generate all internal URLs. |
| `EntityNameSingular` | `string` | No | Singular name (e.g., "Product"). Used in headings and labels. |
| `EntityNamePlural` | `string` | No | Plural name (e.g., "Products"). Used in headings. |
| `HomeURL` | `string` | No | URL for the "Home" breadcrumb link. |
| `ColumnNames` | `[]string` | No | Column headers for the entity manager table. Supports `{!! !!}` raw HTML. |
| `CreateFields` | `[]form.FieldInterface` | No | Form fields for the create modal. |
| `UpdateFields` | `[]form.FieldInterface` | **Yes** | Must not be `nil`. Can be an empty slice. |
| `FileManagerURL` | `string` | No | URL for file manager. Image fields display a "Browse" link when set. |

### Redirect Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `CreateRedirectURL` | `string` | No | Redirect URL after entity creation. Default: update page for the new entity. |
| `UpdateRedirectURL` | `string` | No | Redirect URL after entity update (Save button). Default: entity manager page. |

### Server-Side Pagination

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `PageSize` | `int` | No | Rows per page. `0` (default) = all rows loaded, DataTables handles client-side pagination. |
| `FuncRowsCount` | `func(r *http.Request) (int64, error)` | No | Total row count. Required when `PageSize > 0`. |

When `PageSize > 0` and `FuncRowsCount` is set, Bootstrap pagination controls are rendered below the table. The current page is read from the `page` query parameter. Your `FuncRows` implementation should read `page` and use `PageSize` to return only the relevant slice.

### Middleware Hooks

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncBeforeAction` | `func(w http.ResponseWriter, r *http.Request, action string) bool` | No | Called before every action. Return `false` to abort. |
| `FuncAfterAction` | `func(w http.ResponseWriter, r *http.Request, action string)` | No | Called after every action. Skipped if before-action aborted. |

### CSRF Validation

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncValidateCSRF` | `func(r *http.Request) error` | No | Called on POST requests before the controller action. Return non-nil error to reject. |

### Structured Logging

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncLog` | `func(level string, message string, attrs map[string]any)` | No | Called for key events. Levels: `LogLevelInfo`, `LogLevelWarn`, `LogLevelError`, `LogLevelDebug`. |

## Callback Functions

All callbacks receive `*http.Request` as their first parameter, giving access to the authenticated user, headers, cookies, and request context.

### FuncRows

```go
func(r *http.Request) (rows []crud.Row, err error)
```

Returns entities as rows. Each row's `Data` slice must match `ColumnNames` length. When server-side pagination is enabled, read the `page` query parameter to return only the relevant page.

### FuncCreate

```go
func(r *http.Request, data map[string]string) (entityID string, err error)
```

Receives form field values keyed by field name. Returns the new entity's ID on success.

### FuncFetchReadData

```go
func(r *http.Request, entityID string) ([]crud.KeyValue, error)
```

Returns an ordered list of key-value pairs. Keys and values wrapped in `{!! !!}` are rendered as raw HTML.

### FuncFetchUpdateData

```go
func(r *http.Request, entityID string) (map[string]string, error)
```

Returns current field values keyed by field name, used to populate the update form.

### FuncUpdate

```go
func(r *http.Request, entityID string, data map[string]string) error
```

Receives the entity ID and updated field values.

### FuncTrash

```go
func(r *http.Request, entityID string) error
```

Soft-deletes the entity by ID.

### FuncLayout

```go
func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string
```

Custom layout wrapper. The package prepends required CDN URLs to `jsFiles` and `styleFiles` before calling this function.

### FuncReadExtras

```go
func(r *http.Request, entityID string) []hb.TagInterface
```

Returns additional HTML elements to append below the entity read view card.

### FuncBeforeAction

```go
func(w http.ResponseWriter, r *http.Request, action string) bool
```

Return `true` to continue, `false` to abort. The `action` parameter is the route path (e.g., `"entity-manager"`, `"entity-create-ajax"`).

### FuncAfterAction

```go
func(w http.ResponseWriter, r *http.Request, action string)
```

Called after every controller action completes. Not called if before-action aborted.

### FuncValidateCSRF

```go
func(r *http.Request) error
```

Return `nil` to allow the request, or an error to reject with a JSON error response.

### FuncLog

```go
func(level string, message string, attrs map[string]any)
```

The `attrs` map contains structured context (e.g., `"action"`, `"method"`, `"url"`, `"error"`).

## Validation Rules

The `New()` function enforces:

1. `FuncRows` must not be `nil`
2. `UpdateFields` must not be `nil` (but can be an empty slice)
3. If `UpdateFields` contains entries:
   - `FuncUpdate` must not be `nil`
   - `FuncFetchUpdateData` must not be `nil`

## See Also

- [Getting Started](getting_started.md) - Quick start guide
- [API Reference](api_reference.md) - Complete API documentation
- [Modules: Config](modules/config.md) - Config module details
- [Modules: Form System](modules/form_system.md) - Form field types and options
