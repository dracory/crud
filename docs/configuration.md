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

## Config Struct

### Core Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Endpoint` | `string` | No | The base URL path for the CRUD interface. Used to generate all internal URLs. |
| `EntityNameSingular` | `string` | No | Singular name of the entity (e.g., "Product"). Used in headings and labels. |
| `EntityNamePlural` | `string` | No | Plural name of the entity (e.g., "Products"). Used in headings. |
| `HomeURL` | `string` | No | URL for the "Home" breadcrumb link. |
| `ColumnNames` | `[]string` | No | Column headers for the entity manager table. Supports raw HTML via `{!! !!}` wrapper. |
| `CreateFields` | `[]form.FieldInterface` | No | Form fields for the create modal. |
| `UpdateFields` | `[]form.FieldInterface` | **Yes** | Must not be `nil`. Can be an empty slice if update is not needed. |
| `FileManagerURL` | `string` | No | URL for a file manager. When set, image fields display a "Browse" link. |

### Callback Functions

All callback functions receive `*http.Request` as their first parameter, giving access to the authenticated user, headers, cookies, and request context.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncRows` | `func(r *http.Request) ([]Row, error)` | **Yes** | Returns rows to display in the entity manager table. |
| `FuncCreate` | `func(r *http.Request, data map[string]string) (string, error)` | No | Creates an entity and returns its ID. If `nil`, create AJAX endpoint returns an error. |
| `FuncFetchReadData` | `func(r *http.Request, entityID string) ([]KeyValue, error)` | No | Returns key-value pairs for the read view. If `nil`, read page returns an error. View button is hidden in the manager. |
| `FuncFetchUpdateData` | `func(r *http.Request, entityID string) (map[string]string, error)` | Conditional | Required when `UpdateFields` has entries. Returns current field values for the entity. |
| `FuncUpdate` | `func(r *http.Request, entityID string, data map[string]string) error` | Conditional | Required when `UpdateFields` has entries. |
| `FuncTrash` | `func(r *http.Request, entityID string) error` | No | Soft-deletes an entity. If `nil`, trash endpoint returns an error. Trash button is hidden in the manager. |
| `FuncLayout` | `func(w http.ResponseWriter, r *http.Request, title, content string, styleFiles []string, style string, jsFiles []string, js string) string` | No | Custom layout function. When `nil`, a default Bootstrap webpage is rendered. |
| `FuncReadExtras` | `func(r *http.Request, entityID string) []hb.TagInterface` | No | Returns additional HTML elements appended below the read view card. |

### Redirect Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `CreateRedirectURL` | `string` | No | URL to redirect to after successful entity creation. Default: the update page for the new entity. |
| `UpdateRedirectURL` | `string` | No | URL to redirect to after successful entity update (Save button). Default: the entity manager page. |

### Server-Side Pagination

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `PageSize` | `int` | No | Number of rows per page. When `0` (default), all rows are loaded and DataTables handles client-side pagination. |
| `FuncRowsCount` | `func(r *http.Request) (int64, error)` | No | Returns the total number of rows. Required when `PageSize > 0` to render pagination controls. |

When `PageSize > 0` and `FuncRowsCount` is set, the entity manager renders Bootstrap pagination controls below the table. The current page is read from the `page` query parameter. Your `FuncRows` implementation should read `page` and use `PageSize` to return only the relevant slice of rows.

### Middleware Hooks

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncBeforeAction` | `func(w http.ResponseWriter, r *http.Request, action string) bool` | No | Called before every controller action. Return `false` to abort the request. The `action` parameter is the route path (e.g., `"entity-manager"`, `"entity-create-ajax"`). |
| `FuncAfterAction` | `func(w http.ResponseWriter, r *http.Request, action string)` | No | Called after every controller action completes. Not called if `FuncBeforeAction` aborted the request. |

### CSRF Validation

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncValidateCSRF` | `func(r *http.Request) error` | No | Called on every POST request before the controller action. Return a non-nil error to reject the request with a JSON error response. When `nil`, no CSRF validation is performed. Only called for POST requests (GET requests are never validated). |

### Structured Logging

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `FuncLog` | `func(level string, message string, attrs map[string]any)` | No | Called for key events: request handling (`"info"`), before-action aborts (`"warn"`), and CSRF failures (`"warn"`). When `nil`, no logging occurs. Log levels: `LogLevelInfo`, `LogLevelWarn`, `LogLevelError`, `LogLevelDebug`. |

## KeyValue Type

Used by `FuncFetchReadData` to return structured key-value pairs for the read view:

```go
type KeyValue struct {
    Key   string
    Value string
}
```

Both `Key` and `Value` support the `{!! !!}` raw HTML syntax.

## Validation Rules

The `New()` function enforces the following validation:

1. `FuncRows` must not be `nil`.
2. `UpdateFields` must not be `nil` (but can be an empty slice).
3. If `UpdateFields` contains entries:
   - `FuncUpdate` must not be `nil`.
   - `FuncFetchUpdateData` must not be `nil`.

## Row Struct

The `Row` struct represents a single row in the entity manager table:

```go
type Row struct {
    ID   string   // Unique entity identifier
    Data []string // Cell values, must match ColumnNames length
}
```

## Callback Function Signatures

### FuncRows

```go
func(r *http.Request) (rows []crud.Row, err error)
```

Returns entities as rows. Each row's `Data` slice must have the same length as `ColumnNames`. When server-side pagination is enabled (`PageSize > 0`), read the `page` query parameter from `r` to return only the relevant page of results.

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

Custom layout wrapper. The package prepends required CDN URLs (Vue.js, Sweetalert2, HTMX, Element Plus) to `jsFiles` and `styleFiles` before calling this function.

### FuncReadExtras

```go
func(r *http.Request, entityID string) []hb.TagInterface
```

Returns additional HTML elements to append below the entity read view card.

### FuncRowsCount

```go
func(r *http.Request) (int64, error)
```

Returns the total number of rows matching the current query. Used to calculate pagination when `PageSize > 0`.

### FuncBeforeAction

```go
func(w http.ResponseWriter, r *http.Request, action string) bool
```

Called before every controller action. Return `true` to continue, `false` to abort. Use this for per-action authorization, rate limiting, or custom headers.

### FuncAfterAction

```go
func(w http.ResponseWriter, r *http.Request, action string)
```

Called after every controller action completes successfully (not called if before-action aborted).

### FuncValidateCSRF

```go
func(r *http.Request) error
```

Called on POST requests before the controller action. Return `nil` to allow the request, or an error to reject it with a JSON error response.

### FuncLog

```go
func(level string, message string, attrs map[string]any)
```

Called for key events during request handling. The `level` parameter is one of `"info"`, `"warn"`, `"error"`, or `"debug"`. The `attrs` map contains structured context (e.g., `"action"`, `"method"`, `"url"`, `"error"`).
