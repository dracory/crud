# Package Overview

## Introduction

The `github.com/dracory/crud/v2` package provides a complete, server-side rendered CRUD (Create, Read, Update, Trash) interface for Go web applications. It generates Bootstrap 5-based HTML pages with Vue.js 3 interactivity, allowing developers to quickly scaffold admin panels and entity management UIs.

## Architecture

The package follows a controller-based architecture where a single `Crud` struct acts as the central router and configuration holder. All routing is handled through a single `Handler` method that dispatches requests to internal controllers based on a `path` query parameter.

```
HTTP Request
    │
    ▼
Crud.Handler(w, r)
    │
    ├── 1. FuncLog ("handling request")        [if configured]
    ├── 2. FuncBeforeAction(w, r, action)       [if configured; return false to abort]
    ├── 3. FuncValidateCSRF(r)                  [if configured; POST only]
    ├── 4. Route dispatch:
    │      ├── path=""  or "home"  ──► entityManagerController.page
    │      ├── path="entity-manager"  ──► entityManagerController.page
    │      ├── path="entity-create-modal" ──► entityCreateController.modalShow
    │      ├── path="entity-create-ajax"  ──► entityCreateController.modalSave  (POST only)
    │      ├── path="entity-read"         ──► entityReadController.page
    │      ├── path="entity-update"       ──► entityUpdateController.page
    │      ├── path="entity-update-ajax"  ──► entityUpdateController.pageSave   (POST only)
    │      └── path="entity-trash-ajax"   ──► entityTrashController.pageEntityTrashAjax (POST only)
    └── 5. FuncAfterAction(w, r, action)        [if configured; skipped if aborted]
```

Unknown paths fall back to the entity manager page (home).

## Core Types

| Type | File | Description |
|------|------|-------------|
| `Crud` | `crud.go` | Central struct holding configuration, routing, layout, form rendering, URL helpers, and breadcrumbs |
| `Config` | `config.go` | Configuration struct passed to `New()` to create a `Crud` instance |
| `Row` | `row.go` | Represents a single row in the entity manager table |
| `KeyValue` | `key_value.go` | Key-value pair used by `FuncFetchReadData` for the read view |
| `Breadcrumb` | `breadcrumb.go` | Represents a breadcrumb navigation item |

## Controllers (internal)

| Controller | File | Responsibility |
|------------|------|----------------|
| `entityManagerController` | `entity_manager_controller.go` | Lists all entities in a DataTable with view/edit/trash actions |
| `entityCreateController` | `entity_create_controller.go` | Shows a Bootstrap modal for creating entities via HTMX; saves via AJAX |
| `entityReadController` | `entity_read_controller.go` | Displays entity details in a read-only table |
| `entityUpdateController` | `entity_update_controller.go` | Renders an edit form with Vue.js two-way binding; saves via AJAX |
| `entityTrashController` | `entity_trash_controller.go` | Handles soft-delete (trash) via AJAX POST |

## Frontend Stack

The default layout (when no custom `FuncLayout` is provided) includes the following CDN-loaded libraries:

- **Bootstrap 5.3.3** - CSS framework and JS components
- **jQuery 3.7.1** - DOM manipulation and AJAX
- **Vue.js 3** - Reactive UI for forms and entity management
- **Sweetalert2 11** - User-friendly alert dialogs
- **HTMX 2.0.0** - HTML-over-the-wire for modal loading (added when custom layout is used)
- **Element Plus 2.3.8** - Vue 3 date picker component (added when custom layout is used)
- **jQuery DataTables 1.13.4** - Sortable/searchable entity table
- **Trumbowyg 2.27.3** - WYSIWYG HTML editor for `htmlarea` fields
- **Vue Trumbowyg 4** - Vue 3 wrapper for Trumbowyg

## Dependencies

See `go.mod` for the full list. Key dependencies:

| Package | Purpose |
|---------|---------|
| `github.com/dracory/api` | JSON API response helpers (`api.Respond`, `api.Error`, `api.SuccessWithData`) |
| `github.com/dracory/bs` | Bootstrap component builder (modals, input groups) |
| `github.com/dracory/cdn` | CDN URL helpers for frontend libraries |
| `github.com/dracory/form` | Form field interface and builder |
| `github.com/dracory/hb` | HTML builder for server-side HTML generation |
| `github.com/dracory/req` | HTTP request parameter extraction |
| `github.com/dracory/str` | String utilities (random ID generation) |
| `github.com/samber/lo` | Functional programming helpers (map, ternary, forEach) |

## Security

- **XSS Prevention**: All user-controlled values interpolated into inline JavaScript are JSON-encoded via `json.Marshal` before embedding.
- **HTTP Method Enforcement**: All state-mutating endpoints (create, update, trash) enforce `POST` method and reject other methods.
- **Nil Safety**: All optional callback functions (`FuncCreate`, `FuncTrash`, `FuncFetchUpdateData`, etc.) are checked for `nil` before invocation.
- **CSRF Validation**: Optional `FuncValidateCSRF` hook validates POST requests before controller actions. Returns a JSON error if validation fails.
- **Middleware Hooks**: `FuncBeforeAction` enables per-action authorization (e.g., "this user can read but not trash").
- **Authentication/Authorization**: Not handled directly by this package. Use `FuncBeforeAction` for per-action checks, or wrap `Crud.Handler` with your own middleware.
