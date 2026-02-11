---
path: cheatsheet.md
page-type: reference
summary: Quick reference for common operations, patterns, and code snippets when using the Dracory CRUD package.
tags: [cheatsheet, quick-reference, patterns, snippets]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Cheatsheet

## Installation

```bash
go get github.com/dracory/crud/v2
```

```go
import crud "github.com/dracory/crud/v2"
import "github.com/dracory/form"
```

## Minimum Config

```go
c, err := crud.New(crud.Config{
    FuncRows:     func(r *http.Request) ([]crud.Row, error) { return nil, nil },
    UpdateFields: []form.FieldInterface{},
})
```

## Full Config Template

```go
c, err := crud.New(crud.Config{
    Endpoint:           "/admin/products",
    EntityNameSingular: "Product",
    EntityNamePlural:   "Products",
    HomeURL:            "/dashboard",
    ColumnNames:        []string{"ID", "Name", "Status"},
    CreateFields:       createFields,
    UpdateFields:       updateFields,
    FileManagerURL:     "/files",
    CreateRedirectURL:  "/admin/products",
    UpdateRedirectURL:  "/admin/products",
    PageSize:           25,
    FuncRows:           fetchRows,
    FuncRowsCount:      countRows,
    FuncCreate:         createEntity,
    FuncFetchReadData:  fetchReadData,
    FuncFetchUpdateData: fetchUpdateData,
    FuncUpdate:         updateEntity,
    FuncTrash:          trashEntity,
    FuncLayout:         myLayout,
    FuncReadExtras:     readExtras,
    FuncBeforeAction:   beforeAction,
    FuncAfterAction:    afterAction,
    FuncValidateCSRF:   validateCSRF,
    FuncLog:            logFunc,
})
```

## Field Types

```go
crud.FORM_FIELD_TYPE_STRING       // Text input
crud.FORM_FIELD_TYPE_NUMBER       // Number input
crud.FORM_FIELD_TYPE_PASSWORD     // Password input
crud.FORM_FIELD_TYPE_TEXTAREA     // Multi-line text
crud.FORM_FIELD_TYPE_SELECT       // Dropdown select
crud.FORM_FIELD_TYPE_IMAGE        // Image URL + preview
crud.FORM_FIELD_TYPE_IMAGE_INLINE // Image upload (base64)
crud.FORM_FIELD_TYPE_HTMLAREA     // WYSIWYG editor
crud.FORM_FIELD_TYPE_BLOCKAREA    // Block editor
crud.FORM_FIELD_TYPE_DATETIME     // Date/time picker
crud.FORM_FIELD_TYPE_RAW          // Raw HTML (no label)
```

## Creating Fields

```go
// Required text field
form.NewField(form.FieldOptions{
    Name: "title", Label: "Title",
    Type: crud.FORM_FIELD_TYPE_STRING, Required: true,
})

// Select with static options
form.NewField(form.FieldOptions{
    Name: "status", Label: "Status",
    Type: crud.FORM_FIELD_TYPE_SELECT,
    Options: []form.FieldOption{
        {Key: "active", Value: "Active"},
        {Key: "draft", Value: "Draft"},
    },
})

// Select with dynamic options
form.NewField(form.FieldOptions{
    Name: "category", Label: "Category",
    Type: crud.FORM_FIELD_TYPE_SELECT,
    OptionsF: func() []form.FieldOption {
        return fetchCategoriesFromDB()
    },
})

// Textarea with help text
form.NewField(form.FieldOptions{
    Name: "description", Label: "Description",
    Type: crud.FORM_FIELD_TYPE_TEXTAREA,
    Help: "Enter a brief description (max 500 characters)",
})

// Raw HTML (hidden field, custom content, etc.)
form.NewField(form.FieldOptions{
    Type:  crud.FORM_FIELD_TYPE_RAW,
    Value: `<hr><h5>Advanced Settings</h5>`,
})
```

## Callback Signatures

```go
// FuncRows - List entities
func(r *http.Request) ([]crud.Row, error)

// FuncCreate - Create entity, return ID
func(r *http.Request, data map[string]string) (string, error)

// FuncFetchReadData - Read view data
func(r *http.Request, entityID string) ([]crud.KeyValue, error)

// FuncFetchUpdateData - Update form data
func(r *http.Request, entityID string) (map[string]string, error)

// FuncUpdate - Save entity changes
func(r *http.Request, entityID string, data map[string]string) error

// FuncTrash - Soft-delete entity
func(r *http.Request, entityID string) error

// FuncLayout - Custom page layout
func(w http.ResponseWriter, r *http.Request, title, content string,
    styleFiles []string, style string, jsFiles []string, js string) string

// FuncBeforeAction - Pre-action hook (return false to abort)
func(w http.ResponseWriter, r *http.Request, action string) bool

// FuncAfterAction - Post-action hook
func(w http.ResponseWriter, r *http.Request, action string)

// FuncValidateCSRF - CSRF validation (POST only)
func(r *http.Request) error

// FuncLog - Structured logging
func(level string, message string, attrs map[string]any)

// FuncRowsCount - Total row count for pagination
func(r *http.Request) (int64, error)

// FuncReadExtras - Extra HTML below read view
func(r *http.Request, entityID string) []hb.TagInterface
```

## URL Helpers

```go
c.UrlEntityManager()     // "?path=entity-manager"
c.UrlEntityCreateModal() // "?path=entity-create-modal"
c.UrlEntityCreateAjax()  // "?path=entity-create-ajax"
c.UrlEntityRead()        // "?path=entity-read"
c.UrlEntityUpdate()      // "?path=entity-update"
c.UrlEntityUpdateAjax()  // "?path=entity-update-ajax"
c.UrlEntityTrashAjax()   // "?path=entity-trash-ajax"
```

## Raw HTML in Columns

```go
ColumnNames: []string{"ID", "{!!Status!!}"}
// Cells in the "Status" column will render as raw HTML
```

## Pagination Setup

```go
crud.Config{
    PageSize: 25,
    FuncRowsCount: func(r *http.Request) (int64, error) {
        return db.Count("products")
    },
    FuncRows: func(r *http.Request) ([]crud.Row, error) {
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        if page < 1 { page = 1 }
        offset := (page - 1) * 25
        return db.Query("SELECT ... LIMIT 25 OFFSET ?", offset)
    },
}
```

## Authorization via FuncBeforeAction

```go
FuncBeforeAction: func(w http.ResponseWriter, r *http.Request, action string) bool {
    user := getUser(r)
    if user == nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return false
    }
    if action == "entity-trash-ajax" && !user.IsAdmin {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return false
    }
    return true
},
```

## Running Tests

```bash
cd v2
go test ./... -v
```

## See Also

- [Getting Started](getting_started.md) - Full setup guide
- [Configuration](configuration.md) - Complete config reference
- [API Reference](api_reference.md) - Full API docs
