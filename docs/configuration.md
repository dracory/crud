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
    ReadFields:         readFields,
})
if err != nil {
    log.Fatal(err)
}

http.HandleFunc("/admin", c.Handler)
```

## Config Struct

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Endpoint` | `string` | No | The base URL path for the CRUD interface. Used to generate all internal URLs. |
| `EntityNameSingular` | `string` | No | Singular name of the entity (e.g., "Product"). Used in headings and labels. |
| `EntityNamePlural` | `string` | No | Plural name of the entity (e.g., "Products"). Used in headings. |
| `HomeURL` | `string` | No | URL for the "Home" breadcrumb link. |
| `ColumnNames` | `[]string` | No | Column headers for the entity manager table. Supports raw HTML via `{!! !!}` wrapper. |
| `FuncRows` | `func() ([]Row, error)` | **Yes** | Returns all rows to display in the entity manager table. |
| `UpdateFields` | `[]form.FieldInterface` | **Yes** | Must not be `nil`. Can be an empty slice if update is not needed. |
| `FuncUpdate` | `func(entityID string, data map[string]string) error` | Conditional | Required when `UpdateFields` has entries. |
| `FuncFetchUpdateData` | `func(entityID string) (map[string]string, error)` | Conditional | Required when `UpdateFields` has entries. Returns current field values for the entity. |
| `FuncCreate` | `func(data map[string]string) (string, error)` | No | Creates an entity and returns its ID. If `nil`, create AJAX endpoint returns an error. |
| `FuncFetchReadData` | `func(entityID string) ([][2]string, error)` | No | Returns key-value pairs for the read view. If `nil`, read page returns an error. View button is hidden in the manager. |
| `FuncTrash` | `func(entityID string) error` | No | Soft-deletes an entity. If `nil`, trash endpoint returns an error. Trash button is hidden in the manager. |
| `FuncLayout` | `func(w, r, title, content, styleFiles, style, jsFiles, js) string` | No | Custom layout function. When provided, the package prepends its required CDN resources and delegates HTML generation to this function. When `nil`, a default Bootstrap webpage is rendered. |
| `FuncReadExtras` | `func(entityID string) []hb.TagInterface` | No | Returns additional HTML elements appended below the read view card. |
| `CreateFields` | `[]form.FieldInterface` | No | Form fields for the create modal. |
| `ReadFields` | `[]form.FieldInterface` | No | Form fields for the read view. |
| `FileManagerURL` | `string` | No | URL for a file manager. When set, image fields display a "Browse" link. |

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
func() (rows []crud.Row, err error)
```

Returns all entities as rows. Each row's `Data` slice must have the same length as `ColumnNames`.

### FuncCreate

```go
func(data map[string]string) (entityID string, err error)
```

Receives form field values keyed by field name. Returns the new entity's ID on success.

### FuncFetchReadData

```go
func(entityID string) ([][2]string, error)
```

Returns an ordered list of key-value pairs. Keys and values wrapped in `{!! !!}` are rendered as raw HTML.

### FuncFetchUpdateData

```go
func(entityID string) (map[string]string, error)
```

Returns current field values keyed by field name, used to populate the update form.

### FuncUpdate

```go
func(entityID string, data map[string]string) error
```

Receives the entity ID and updated field values.

### FuncTrash

```go
func(entityID string) error
```

Soft-deletes the entity by ID.

### FuncLayout

```go
func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string
```

Custom layout wrapper. The package prepends required CDN URLs (Vue.js, Sweetalert2, HTMX, Element Plus) to `jsFiles` and `styleFiles` before calling this function.

### FuncReadExtras

```go
func(entityID string) []hb.TagInterface
```

Returns additional HTML elements to append below the entity read view card.
