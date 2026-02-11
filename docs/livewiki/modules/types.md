---
path: modules/types.md
page-type: module
summary: Documentation for the supporting data types - Row, KeyValue, and Breadcrumb.
tags: [module, types, row, keyvalue, breadcrumb]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: Types

## Purpose

The package defines three simple supporting data types used by callback functions and controllers: `Row`, `KeyValue`, and `Breadcrumb`.

## Row

**File:** `v2/row.go`

Represents a single row in the entity manager table.

```go
type Row struct {
    ID   string   // Unique entity identifier
    Data []string // Cell values, must match ColumnNames length
}
```

**Used by:** `FuncRows` callback → entity manager controller

**Example:**

```go
FuncRows: func(r *http.Request) ([]crud.Row, error) {
    return []crud.Row{
        {ID: "1", Data: []string{"1", "Widget", "Active"}},
        {ID: "2", Data: []string{"2", "Gadget", "Draft"}},
    }, nil
},
```

**Notes:**
- `ID` is used to generate View, Edit, and Trash action URLs
- `Data` length must match `ColumnNames` length
- Cell values support `{!! !!}` raw HTML syntax (when the corresponding column name is also wrapped)

## KeyValue

**File:** `v2/key_value.go`

Represents a key-value pair used in the entity read view.

```go
type KeyValue struct {
    Key   string
    Value string
}
```

**Used by:** `FuncFetchReadData` callback → entity read controller

**Example:**

```go
FuncFetchReadData: func(r *http.Request, entityID string) ([]crud.KeyValue, error) {
    return []crud.KeyValue{
        {Key: "Name", Value: "Widget"},
        {Key: "Status", Value: "Active"},
        {Key: "{!!Image!!}", Value: `{!!<img src="https://example.com/img.png" width="200">!!}`},
    }, nil
},
```

**Notes:**
- Both `Key` and `Value` support `{!! !!}` raw HTML syntax
- The order of the slice determines the display order in the read view table
- Keys are rendered in `<th>` elements, values in `<td>` elements

## Breadcrumb

**File:** `v2/breadcrumb.go`

Represents a breadcrumb navigation item.

```go
type Breadcrumb struct {
    Name string // Display text
    URL  string // Link href
}
```

**Used by:** Controllers internally (not exposed to consumers)

**Breadcrumb trails by controller:**

| Controller | Trail |
|------------|-------|
| Entity Manager | Home → Entity Manager |
| Entity Read | Home → Entity Manager → View Entity |
| Entity Update | Home → Entity Manager → Edit Entity |

**Rendered as:**

```html
<nav aria-label="breadcrumb">
  <ol class="breadcrumb">
    <li class="breadcrumb-item"><a href="/">Home</a></li>
    <li class="breadcrumb-item"><a href="/admin?path=entity-manager">Product Manager</a></li>
  </ol>
</nav>
```

## See Also

- [Modules: Controllers](controllers.md) - How types are used in controllers
- [Modules: Config](config.md) - Callback function signatures
- [API Reference](../api_reference.md) - Type definitions
