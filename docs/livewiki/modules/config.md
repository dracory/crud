---
path: modules/config.md
page-type: module
summary: Documentation for the Config struct and the New() constructor with validation rules.
tags: [module, config, constructor, validation]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: Config

## Purpose

The `Config` struct is the public API surface for configuring a CRUD instance. It defines all entity metadata, callback functions, form fields, and optional features. The `New()` constructor validates the config and creates a `Crud` instance.

**Files:** `v2/config.go`, `v2/new.go`

## Config Struct

```go
type Config struct {
    ColumnNames         []string
    CreateFields        []form.FieldInterface
    Endpoint            string
    EntityNamePlural    string
    EntityNameSingular  string
    FileManagerURL      string
    FuncCreate          func(r *http.Request, data map[string]string) (string, error)
    FuncFetchReadData   func(r *http.Request, entityID string) ([]KeyValue, error)
    FuncFetchUpdateData func(r *http.Request, entityID string) (map[string]string, error)
    FuncLayout          func(w http.ResponseWriter, r *http.Request, title string, content string, styleFiles []string, style string, jsFiles []string, js string) string
    FuncRows            func(r *http.Request) (rows []Row, err error)
    FuncTrash           func(r *http.Request, entityID string) error
    FuncUpdate          func(r *http.Request, entityID string, data map[string]string) error
    HomeURL             string
    CreateRedirectURL   string
    UpdateRedirectURL   string
    UpdateFields        []form.FieldInterface
    FuncReadExtras      func(r *http.Request, entityID string) []hb.TagInterface
    PageSize            int
    FuncRowsCount       func(r *http.Request) (int64, error)
    FuncBeforeAction    func(w http.ResponseWriter, r *http.Request, action string) bool
    FuncAfterAction     func(w http.ResponseWriter, r *http.Request, action string)
    FuncValidateCSRF    func(r *http.Request) error
    FuncLog             func(level string, message string, attrs map[string]any)
}
```

## Constructor

### New(config Config) (Crud, error)

Creates a `Crud` instance by copying all config fields into the private struct. Returns an error if validation fails.

**Validation rules:**

| Rule | Error Message |
|------|---------------|
| `FuncRows` is `nil` | `"FuncRows function is required"` |
| `UpdateFields` is `nil` | `"UpdateFields is required"` |
| `UpdateFields` has entries but `FuncUpdate` is `nil` | `"FuncUpdate function is required when UpdateFields are provided"` |
| `UpdateFields` has entries but `FuncFetchUpdateData` is `nil` | `"FuncFetchUpdateData function is required when UpdateFields are provided"` |

**Note:** An empty `UpdateFields` slice (`[]form.FieldInterface{}`) is valid and does not require `FuncUpdate` or `FuncFetchUpdateData`.

## Field Categories

### Required Fields

| Field | Why Required |
|-------|-------------|
| `FuncRows` | The entity manager always needs to list entities |
| `UpdateFields` | Must be explicitly set (even if empty) to avoid accidental nil |

### Conditionally Required

| Field | Condition |
|-------|-----------|
| `FuncUpdate` | Required when `UpdateFields` has entries |
| `FuncFetchUpdateData` | Required when `UpdateFields` has entries |

### Optional Fields

All other fields are optional. When `nil` or zero-value:
- **FuncCreate nil** → Create AJAX returns error; create button still shows
- **FuncFetchReadData nil** → View button hidden in entity manager
- **FuncTrash nil** → Trash button hidden in entity manager
- **FuncLayout nil** → Default standalone Bootstrap page used
- **FuncReadExtras nil** → No extra content below read view
- **PageSize 0** → No server-side pagination; DataTables handles client-side
- **FuncBeforeAction nil** → No pre-action middleware
- **FuncAfterAction nil** → No post-action middleware
- **FuncValidateCSRF nil** → No CSRF validation
- **FuncLog nil** → No logging

## Usage Example

```go
c, err := crud.New(crud.Config{
    Endpoint:           "/admin/users",
    EntityNameSingular: "User",
    EntityNamePlural:   "Users",
    HomeURL:            "/admin",
    ColumnNames:        []string{"ID", "Name", "Email", "Role"},
    FuncRows:           userRepository.ListRows,
    FuncCreate:         userRepository.Create,
    FuncFetchReadData:  userRepository.FetchReadData,
    FuncFetchUpdateData: userRepository.FetchUpdateData,
    FuncUpdate:         userRepository.Update,
    FuncTrash:          userRepository.SoftDelete,
    CreateFields:       userCreateFields,
    UpdateFields:       userUpdateFields,
    PageSize:           50,
    FuncRowsCount:      userRepository.Count,
})
```

## See Also

- [Modules: Crud Core](crud_core.md) - Crud struct details
- [Configuration](../configuration.md) - Full configuration guide
- [Getting Started](../getting_started.md) - Quick start
- [Troubleshooting](../troubleshooting.md) - Constructor error solutions
