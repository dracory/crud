---
path: api_reference.md
page-type: reference
summary: Complete public API reference for the Dracory CRUD v2 package including functions, types, constants, and methods.
tags: [api, reference, functions, types, constants]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# API Reference

## Constructor

### `New(config Config) (Crud, error)`

Creates a new `Crud` instance from the provided configuration. Returns an error if validation fails.

**Validation rules:**
1. `Config.FuncRows` must not be `nil`
2. `Config.UpdateFields` must not be `nil` (but can be an empty slice)
3. If `Config.UpdateFields` has entries:
   - `Config.FuncUpdate` must not be `nil`
   - `Config.FuncFetchUpdateData` must not be `nil`

```go
c, err := crud.New(crud.Config{
    Endpoint:     "/admin",
    FuncRows:     myFetchRows,
    UpdateFields: []form.FieldInterface{},
})
```

## Crud Methods

### `Handler(w http.ResponseWriter, r *http.Request)`

The main HTTP handler. Register this with your router to serve the CRUD interface. Dispatches requests to internal controllers based on the `path` query parameter.

```go
http.HandleFunc("/admin", c.Handler)
```

### URL Helpers

All URL helpers automatically handle query string separator detection (`?` vs `&`).

| Method | Returns | Example |
|--------|---------|---------|
| `UrlEntityManager()` | Entity manager page URL | `/admin?path=entity-manager` |
| `UrlEntityCreateModal()` | Create modal fragment URL | `/admin?path=entity-create-modal` |
| `UrlEntityCreateAjax()` | Create AJAX endpoint URL | `/admin?path=entity-create-ajax` |
| `UrlEntityRead()` | Read page base URL | `/admin?path=entity-read` |
| `UrlEntityUpdate()` | Update page base URL | `/admin?path=entity-update` |
| `UrlEntityUpdateAjax()` | Update AJAX endpoint URL | `/admin?path=entity-update-ajax` |
| `UrlEntityTrashAjax()` | Trash AJAX endpoint URL | `/admin?path=entity-trash-ajax` |

```go
// With simple endpoint
c.UrlEntityManager() // "/admin?path=entity-manager"

// With existing query parameters
// Endpoint = "/admin?module=cms"
c.UrlEntityManager() // "/admin?module=cms&path=entity-manager"
```

## Config Struct

### Core Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Endpoint` | `string` | No | Base URL path for the CRUD interface |
| `EntityNameSingular` | `string` | No | Singular entity name (e.g., "Product") |
| `EntityNamePlural` | `string` | No | Plural entity name (e.g., "Products") |
| `HomeURL` | `string` | No | URL for the "Home" breadcrumb link |
| `ColumnNames` | `[]string` | No | Column headers for the entity manager table |
| `CreateFields` | `[]form.FieldInterface` | No | Form fields for the create modal |
| `UpdateFields` | `[]form.FieldInterface` | **Yes** | Form fields for the update page (can be empty) |
| `FileManagerURL` | `string` | No | URL for file manager (image field "Browse" link) |
| `CreateRedirectURL` | `string` | No | Redirect URL after create (default: update page) |
| `UpdateRedirectURL` | `string` | No | Redirect URL after save (default: entity manager) |
| `PageSize` | `int` | No | Rows per page (0 = no server-side pagination) |

### Callback Functions

| Field | Signature | Required |
|-------|-----------|----------|
| `FuncRows` | `func(r *http.Request) ([]Row, error)` | **Yes** |
| `FuncCreate` | `func(r *http.Request, data map[string]string) (string, error)` | No |
| `FuncFetchReadData` | `func(r *http.Request, entityID string) ([]KeyValue, error)` | No |
| `FuncFetchUpdateData` | `func(r *http.Request, entityID string) (map[string]string, error)` | Conditional |
| `FuncUpdate` | `func(r *http.Request, entityID string, data map[string]string) error` | Conditional |
| `FuncTrash` | `func(r *http.Request, entityID string) error` | No |
| `FuncLayout` | `func(w http.ResponseWriter, r *http.Request, title, content string, styleFiles []string, style string, jsFiles []string, js string) string` | No |
| `FuncReadExtras` | `func(r *http.Request, entityID string) []hb.TagInterface` | No |
| `FuncRowsCount` | `func(r *http.Request) (int64, error)` | No |
| `FuncBeforeAction` | `func(w http.ResponseWriter, r *http.Request, action string) bool` | No |
| `FuncAfterAction` | `func(w http.ResponseWriter, r *http.Request, action string)` | No |
| `FuncValidateCSRF` | `func(r *http.Request) error` | No |
| `FuncLog` | `func(level string, message string, attrs map[string]any)` | No |

## Types

### Row

```go
type Row struct {
    ID   string   // Unique entity identifier
    Data []string // Cell values (must match ColumnNames length)
}
```

### KeyValue

```go
type KeyValue struct {
    Key   string
    Value string
}
```

Both `Key` and `Value` support `{!! !!}` raw HTML syntax.

### Breadcrumb

```go
type Breadcrumb struct {
    Name string // Display text
    URL  string // Link href
}
```

## Constants

### Form Field Types

| Constant | Value | Rendered As |
|----------|-------|-------------|
| `FORM_FIELD_TYPE_STRING` | `"string"` | `<input type="text">` |
| `FORM_FIELD_TYPE_NUMBER` | `"number"` | `<input type="number">` |
| `FORM_FIELD_TYPE_PASSWORD` | `"password"` | `<input type="password">` |
| `FORM_FIELD_TYPE_TEXTAREA` | `"textarea"` | `<textarea>` |
| `FORM_FIELD_TYPE_SELECT` | `"select"` | `<select>` with `<option>` elements |
| `FORM_FIELD_TYPE_IMAGE` | `"image"` | Image preview + URL input + Browse link |
| `FORM_FIELD_TYPE_IMAGE_INLINE` | `"image_inline"` | Image preview + file upload (base64) |
| `FORM_FIELD_TYPE_HTMLAREA` | `"htmlarea"` | Trumbowyg WYSIWYG editor |
| `FORM_FIELD_TYPE_BLOCKAREA` | `"blockarea"` | Block-based content editor |
| `FORM_FIELD_TYPE_DATETIME` | `"datetime"` | Element Plus date/time picker |
| `FORM_FIELD_TYPE_RAW` | `"raw"` | Raw HTML output (no label) |

### Log Levels

| Constant | Value |
|----------|-------|
| `LogLevelInfo` | `"info"` |
| `LogLevelWarn` | `"warn"` |
| `LogLevelError` | `"error"` |
| `LogLevelDebug` | `"debug"` |

## Endpoints

| Method | Path Parameter | Description | Response |
|--------|---------------|-------------|----------|
| GET | (none) / `home` / `entity-manager` | Entity manager listing | HTML |
| GET | `entity-create-modal` | Create modal fragment | HTML fragment |
| POST | `entity-create-ajax` | Create entity | JSON |
| GET | `entity-read` | Read entity view | HTML |
| GET | `entity-update` | Update form | HTML |
| POST | `entity-update-ajax` | Save entity update | JSON |
| POST | `entity-trash-ajax` | Trash entity | JSON |

### JSON Response Format

Success:
```json
{
    "status": "success",
    "message": "Saved successfully",
    "data": { "entity_id": "abc-123" }
}
```

Error:
```json
{
    "status": "error",
    "message": "Entity ID is required",
    "data": null
}
```

## See Also

- [Configuration](configuration.md) - Detailed configuration guide
- [Modules: Crud Core](modules/crud_core.md) - Internal implementation details
- [Modules: Controllers](modules/controllers.md) - Controller behavior details
- [Cheatsheet](cheatsheet.md) - Quick reference
