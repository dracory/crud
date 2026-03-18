# CRUD v2

A server-side rendered CRUD interface for Go web applications. Generates Bootstrap 5 admin pages with Vue.js 3 interactivity for managing entities through Create, Read, Update, and Trash operations.

## Installation

```bash
go get github.com/dracory/crud/v2
```

## Quick Start

```go
package main

import (
	"log"
	"net/http"

	crud "github.com/dracory/crud/v2"
)

func main() {
	c, err := crud.New(crud.Config{
		Endpoint:           "/admin",
		EntityNameSingular: "Product",
		EntityNamePlural:   "Products",
		HomeURL:            "/",
		ColumnNames:        []string{"ID", "Name", "Status"},
		FuncRows: func(r *http.Request) ([]crud.Row, error) {
			return []crud.Row{
				{ID: "1", Data: []string{"1", "Widget", "Active"}},
				{ID: "2", Data: []string{"2", "Gadget", "Draft"}},
			}, nil
		},
		FuncCreate: func(r *http.Request, data map[string]string) (string, error) {
			// Save to database, return new entity ID
			return "3", nil
		},
		FuncFetchUpdateData: func(r *http.Request, entityID string) (map[string]string, error) {
			// Fetch current values from database
			return map[string]string{"name": "Widget", "status": "active"}, nil
		},
		FuncUpdate: func(r *http.Request, entityID string, data map[string]string) error {
			// Update in database
			return nil
		},
		FuncTrash: func(r *http.Request, entityID string) error {
			// Soft-delete in database
			return nil
		},
		FuncFetchReadData: func(r *http.Request, entityID string) ([]crud.KeyValue, error) {
			return []crud.KeyValue{
				{Key: "Name", Value: "Widget"},
				{Key: "Status", Value: "Active"},
			}, nil
		},
		CreateFields: []crud.FieldInterface{
			crud.NewField(crud.FieldOptions{
				Name:     "name",
				Label:    "Name",
				Type:     crud.FORM_FIELD_TYPE_STRING,
				Required: true,
			}),
			crud.NewField(crud.FieldOptions{
				Name:  "status",
				Label: "Status",
				Type:  crud.FORM_FIELD_TYPE_SELECT,
				Options: []crud.FieldOption{
					{Key: "active", Value: "Active"},
					{Key: "draft", Value: "Draft"},
				},
			}),
		},
		UpdateFields: []crud.FieldInterface{
			crud.NewStringField("name", "Name").WithRequired(),
			crud.NewSelectField("status", "Status", []crud.FieldOption{
				{Key: "active", Value: "Active"},
				{Key: "draft", Value: "Draft"},
			}),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/admin", c.Handler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Features

- **Entity Manager** - Sortable DataTable listing with view, edit, and trash actions
- **Create Modal** - Bootstrap modal loaded via HTMX with form validation
- **Read View** - Read-only entity detail page with `KeyValue` pairs
- **Update Form** - Vue.js reactive form with two-way data binding
- **Trash** - Soft-delete with confirmation modal
- **Custom Layout** - Plug in your own layout function to wrap pages in your app shell
- **23 Field Types** - String, number, password, email, tel, url, date, datetime, color, checkbox, radio, textarea, select, image, inline image, HTML editor, block editor, hidden, file, repeater, table, and raw HTML
- **Repeater Fields** - Dynamic add/remove/reorder rows, supports typed sub-fields (any field type) or generic single-value rows; submitted as standard bracket notation (`field[0][key]=val`) and received as a JSON array string
- **No External Form Dependency** - Field types, constructors, and fluent builders are all built into this package
- **Raw HTML Support** - Column names and read view values support `{!! !!}` syntax for raw HTML rendering
- **Middleware Hooks** - `FuncBeforeAction` / `FuncAfterAction` for per-action authorization, audit logging, and custom headers
- **CSRF Validation** - Optional `FuncValidateCSRF` hook validates POST requests before controller actions
- **Server-Side Pagination** - Configurable `PageSize` and `FuncRowsCount` for large datasets
- **Structured Logging** - Optional `FuncLog` callback for centralized logging of requests, aborts, and errors
- **Configurable Redirects** - `CreateRedirectURL` and `UpdateRedirectURL` for custom post-action navigation
- **Unified JSON Responses** - All AJAX endpoints (create, update, trash) return consistent JSON via `api.Respond`
- **Request Context** - All callbacks receive `*http.Request` for access to auth, headers, cookies, and context

## Form Fields

Fields are defined using types built into this package — no external form dependency needed.

### Struct-based

```go
crud.NewField(crud.FieldOptions{
    Name:     "title",
    Label:    "Title",
    Type:     crud.FORM_FIELD_TYPE_STRING,
    Required: true,
})
```

### Fluent builder

Chain methods directly on the constructor — no struct literal needed:

```go
UpdateFields: []crud.FieldInterface{
    crud.NewStringField("name", "Name").
        WithRequired().
        WithHelp("Full product name"),

    crud.NewSelectField("status", "Status", []crud.FieldOption{
        {Key: "active",   Value: "Active"},
        {Key: "draft",    Value: "Draft"},
        {Key: "archived", Value: "Archived"},
    }),

    crud.NewTextAreaField("description", "Description").
        WithHelp("Supports plain text only"),

    crud.NewNumberField("price", "Price (cents)").
        WithRequired(),

    crud.NewHiddenField("source", "web"),

    crud.NewRepeater(crud.RepeaterOptions{
        Name:  "images",
        Label: "Images",
        Fields: []crud.FieldInterface{
            crud.NewURLField("url", "URL").WithRequired(),
            crud.NewStringField("alt", "Alt Text"),
        },
    }),
},
```

### Convenience constructors

| Constructor | Field Type |
|-------------|-----------|
| `NewStringField(name, label)` | `FORM_FIELD_TYPE_STRING` |
| `NewNumberField(name, label)` | `FORM_FIELD_TYPE_NUMBER` |
| `NewPasswordField(name, label)` | `FORM_FIELD_TYPE_PASSWORD` |
| `NewEmailField(name, label)` | `FORM_FIELD_TYPE_EMAIL` |
| `NewTelField(name, label)` | `FORM_FIELD_TYPE_TEL` |
| `NewURLField(name, label)` | `FORM_FIELD_TYPE_URL` |
| `NewDateField(name, label)` | `FORM_FIELD_TYPE_DATE` |
| `NewDateTimeField(name, label)` | `FORM_FIELD_TYPE_DATETIME` |
| `NewColorField(name, label)` | `FORM_FIELD_TYPE_COLOR` |
| `NewCheckboxField(name, label)` | `FORM_FIELD_TYPE_CHECKBOX` |
| `NewRadioField(name, label, options)` | `FORM_FIELD_TYPE_RADIO` |
| `NewTextAreaField(name, label)` | `FORM_FIELD_TYPE_TEXTAREA` |
| `NewSelectField(name, label, options)` | `FORM_FIELD_TYPE_SELECT` |
| `NewImageField(name, label)` | `FORM_FIELD_TYPE_IMAGE` |
| `NewFileField(name, label)` | `FORM_FIELD_TYPE_FILE` |
| `NewHtmlAreaField(name, label)` | `FORM_FIELD_TYPE_HTMLAREA` |
| `NewHiddenField(name, value)` | `FORM_FIELD_TYPE_HIDDEN` |
| `NewRawField(value)` | `FORM_FIELD_TYPE_RAW` |

### Fluent methods

| Method | Description |
|--------|-------------|
| `.WithRequired()` | Mark field as required |
| `.WithLabel(label)` | Set label |
| `.WithHelp(text)` | Set help text |
| `.WithValue(value)` | Set default value |
| `.WithOptions(opts...)` | Set static options |
| `.WithOptionsF(fn)` | Set dynamic options function |
| `.WithFields(fields...)` | Set sub-fields (repeater) |
| `.WithFuncValuesProcess(fn)` | Set value pre-processor for repeater |

## Field Type Constants

| Constant | Description |
|----------|-------------|
| `FORM_FIELD_TYPE_STRING` | Text input |
| `FORM_FIELD_TYPE_NUMBER` | Number input |
| `FORM_FIELD_TYPE_PASSWORD` | Password input |
| `FORM_FIELD_TYPE_EMAIL` | Email input |
| `FORM_FIELD_TYPE_TEL` | Telephone input |
| `FORM_FIELD_TYPE_URL` | URL input |
| `FORM_FIELD_TYPE_DATE` | Date input |
| `FORM_FIELD_TYPE_DATETIME` | Element Plus date/time picker |
| `FORM_FIELD_TYPE_COLOR` | Color picker |
| `FORM_FIELD_TYPE_CHECKBOX` | Checkbox |
| `FORM_FIELD_TYPE_RADIO` | Radio button group |
| `FORM_FIELD_TYPE_TEXTAREA` | Multi-line text area |
| `FORM_FIELD_TYPE_SELECT` | Dropdown select with static or dynamic options |
| `FORM_FIELD_TYPE_IMAGE` | Image URL input with preview and optional file manager |
| `FORM_FIELD_TYPE_IMAGE_INLINE` | Image upload with base64 encoding |
| `FORM_FIELD_TYPE_HTMLAREA` | Trumbowyg WYSIWYG editor |
| `FORM_FIELD_TYPE_BLOCKAREA` | Block-based content editor |
| `FORM_FIELD_TYPE_BLOCKEDITOR` | Block editor (rendered via blockarea JS) |
| `FORM_FIELD_TYPE_HIDDEN` | Hidden input |
| `FORM_FIELD_TYPE_FILE` | File upload input |
| `FORM_FIELD_TYPE_REPEATER` | Dynamic list of rows with add/remove/reorder controls |
| `FORM_FIELD_TYPE_TABLE` | Pre-rendered HTML table (caller provides HTML via `Value`) |
| `FORM_FIELD_TYPE_RAW` | Raw HTML output (no label) |

## Repeater Fields

Repeater fields let users manage a dynamic list of rows directly in the form.

```go
// Generic — each row is a single plain text value
crud.NewField(crud.FieldOptions{
    Name:  "image_urls",
    Label: "Image URLs",
    Type:  crud.FORM_FIELD_TYPE_REPEATER,
})

// Typed — each row is a mini-form with multiple sub-fields
crud.NewRepeater(crud.RepeaterOptions{
    Name:  "links",
    Label: "Links",
    Fields: []crud.FieldInterface{
        crud.NewStringField("label", "Label"),
        crud.NewURLField("url", "URL"),
    },
})

// With a value pre-processor for the update form
crud.NewRepeater(crud.RepeaterOptions{
    Name:  "links",
    Label: "Links",
    Fields: []crud.FieldInterface{
        crud.NewStringField("label", "Label"),
        crud.NewURLField("url", "URL"),
    },
    FuncValuesProcess: func(raw string) []map[string]string {
        var rows []map[string]string
        json.Unmarshal([]byte(raw), &rows)
        return rows
    },
})
```

The browser submits rows using standard bracket notation:

```
links[0][label]=Home&links[0][url]=https://example.com
links[1][label]=About&links[1][url]=https://example.com/about
```

Inside `FuncCreate` / `FuncUpdate`, the repeater field arrives as a JSON array string:

```go
FuncUpdate: func(r *http.Request, entityID string, data map[string]string) error {
    var links []struct {
        Label string `json:"label"`
        URL   string `json:"url"`
    }
    json.Unmarshal([]byte(data["links"]), &links)
    return nil
},
```

Each row has up/down/remove buttons. Up is disabled on the first row, down on the last.

## Endpoints

| Method | Path Parameter | Description |
|--------|---------------|-------------|
| GET | (none) / `home` / `entity-manager` | Entity manager listing |
| GET | `entity-create-modal` | Create modal HTML fragment |
| POST | `entity-create-ajax` | Create entity |
| GET | `entity-read` | Read entity (requires `entity_id`) |
| GET | `entity-update` | Update form (requires `entity_id`) |
| POST | `entity-update-ajax` | Save entity update |
| POST | `entity-trash-ajax` | Trash entity |

## URL Helpers

```go
crud.UrlEntityManager()     // "/admin?path=entity-manager"
crud.UrlEntityCreateModal() // "/admin?path=entity-create-modal"
crud.UrlEntityCreateAjax()  // "/admin?path=entity-create-ajax"
crud.UrlEntityRead()        // "/admin?path=entity-read"
crud.UrlEntityUpdate()      // "/admin?path=entity-update"
crud.UrlEntityUpdateAjax()  // "/admin?path=entity-update-ajax"
crud.UrlEntityTrashAjax()   // "/admin?path=entity-trash-ajax"
```

## Testing

```bash
cd v2
go test ./... -v
```

76 tests covering all controllers, validation, URL generation, breadcrumbs, routing, layout rendering, middleware hooks, CSRF validation, pagination, and structured logging.

## Documentation

Detailed documentation is available in the [`docs/`](../docs/) folder:

- **[Overview](../docs/overview.md)** - Package architecture, core types, frontend stack, dependencies, and security model
- **[Configuration](../docs/configuration.md)** - `Config` struct reference, `New()` validation rules, callback function signatures
- **[Controllers](../docs/controllers.md)** - Endpoint routing, controller behavior, request/response formats
- **[Form Fields](../docs/form-fields.md)** - Field type constants, constructors, fluent builders, Vue.js integration, validation
- **[URL Helpers](../docs/url-helpers.md)** - URL generation methods and route path constants
- **[Layout and Templates](../docs/layout-and-templates.md)** - Layout rendering, breadcrumbs, default webpage template, form rendering

## Frontend Stack

The default layout includes (via CDN):

- Bootstrap 5.3.3
- jQuery 3.7.1
- Vue.js 3
- Sweetalert2 11
- jQuery DataTables 1.13.4
- Trumbowyg 2.27.3 (WYSIWYG editor)
- Element Plus 2.3.8 (date picker)
- HTMX 2.0.0

## Security Notes

- All user-controlled values in inline JavaScript are JSON-encoded to prevent XSS
- State-mutating endpoints enforce POST method
- Optional callbacks are nil-checked before invocation
- Optional `FuncValidateCSRF` validates POST requests before controller actions
- `FuncBeforeAction` enables per-action authorization (e.g., read-only users)
- Authentication is the consumer's responsibility (use `FuncBeforeAction` or external middleware)

## License

See the repository root for license information.
