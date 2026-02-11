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
	"github.com/dracory/form"
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
			// Return key-value pairs for read view
			return []crud.KeyValue{
				{Key: "Name", Value: "Widget"},
				{Key: "Status", Value: "Active"},
			}, nil
		},
		CreateFields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Name:     "name",
				Label:    "Name",
				Type:     crud.FORM_FIELD_TYPE_STRING,
				Required: true,
			}),
			form.NewField(form.FieldOptions{
				Name:  "status",
				Label: "Status",
				Type:  crud.FORM_FIELD_TYPE_SELECT,
				Options: []form.FieldOption{
					{Key: "active", Value: "Active"},
					{Key: "draft", Value: "Draft"},
				},
			}),
		},
		UpdateFields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Name:     "name",
				Label:    "Name",
				Type:     crud.FORM_FIELD_TYPE_STRING,
				Required: true,
			}),
			form.NewField(form.FieldOptions{
				Name:  "status",
				Label: "Status",
				Type:  crud.FORM_FIELD_TYPE_SELECT,
				Options: []form.FieldOption{
					{Key: "active", Value: "Active"},
					{Key: "draft", Value: "Draft"},
				},
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
- **11 Field Types** - String, number, password, textarea, select, image, inline image, HTML editor, block editor, datetime picker, and raw HTML
- **Raw HTML Support** - Column names and read view values support `{!! !!}` syntax for raw HTML rendering
- **Middleware Hooks** - `FuncBeforeAction` / `FuncAfterAction` for per-action authorization, audit logging, and custom headers
- **CSRF Validation** - Optional `FuncValidateCSRF` hook validates POST requests before controller actions
- **Server-Side Pagination** - Configurable `PageSize` and `FuncRowsCount` for large datasets
- **Structured Logging** - Optional `FuncLog` callback for centralized logging of requests, aborts, and errors
- **Configurable Redirects** - `CreateRedirectURL` and `UpdateRedirectURL` for custom post-action navigation
- **Unified JSON Responses** - All AJAX endpoints (create, update, trash) return consistent JSON via `api.Respond`
- **Request Context** - All callbacks receive `*http.Request` for access to auth, headers, cookies, and context

## Form Field Types

| Constant | Description |
|----------|-------------|
| `FORM_FIELD_TYPE_STRING` | Text input |
| `FORM_FIELD_TYPE_NUMBER` | Number input |
| `FORM_FIELD_TYPE_PASSWORD` | Password input |
| `FORM_FIELD_TYPE_TEXTAREA` | Multi-line text area |
| `FORM_FIELD_TYPE_SELECT` | Dropdown select with static or dynamic options |
| `FORM_FIELD_TYPE_IMAGE` | Image URL input with preview and optional file manager |
| `FORM_FIELD_TYPE_IMAGE_INLINE` | Image upload with base64 encoding |
| `FORM_FIELD_TYPE_HTMLAREA` | Trumbowyg WYSIWYG editor |
| `FORM_FIELD_TYPE_BLOCKAREA` | Block-based content editor |
| `FORM_FIELD_TYPE_DATETIME` | Element Plus date/time picker |
| `FORM_FIELD_TYPE_RAW` | Raw HTML output (no label) |

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
- **[Form Fields](../docs/form-fields.md)** - Field type constants, `FieldOptions` reference, Vue.js integration, validation
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
