---
path: modules/url_helpers.md
page-type: module
summary: Documentation for URL generation methods and route path constants used for internal navigation.
tags: [module, url, helpers, routing, constants]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: URL Helpers

## Purpose

URL helpers generate endpoint URLs with the correct `path` query parameter. They are used internally by controllers for navigation links and AJAX endpoints, and are also available publicly for consumers to generate links.

**File:** `v2/crud.go` (methods), `v2/constants.go` (path constants)

## Query String Detection

All URL helpers automatically detect whether the base `Endpoint` already contains a query string (`?`) and append the `path` parameter with the correct separator:

```go
q := lo.Ternary(strings.Contains(crud.endpoint, "?"), "&", "?")
url := crud.endpoint + q + "path=" + pathValue
```

## Public Methods

### UrlEntityManager

```go
func (crud *Crud) UrlEntityManager() string
```

Returns the URL for the entity manager page.

```go
// Endpoint = "/admin"
crud.UrlEntityManager() // "/admin?path=entity-manager"

// Endpoint = "/admin?module=cms"
crud.UrlEntityManager() // "/admin?module=cms&path=entity-manager"
```

### UrlEntityCreateModal

```go
func (crud *Crud) UrlEntityCreateModal() string
```

Returns the URL for the create modal HTML fragment (loaded via HTMX).

### UrlEntityCreateAjax

```go
func (crud *Crud) UrlEntityCreateAjax() string
```

Returns the URL for the create AJAX endpoint (POST).

### UrlEntityRead

```go
func (crud *Crud) UrlEntityRead() string
```

Returns the base URL for the entity read page. Append `&entity_id=ID` for the full URL.

### UrlEntityUpdate

```go
func (crud *Crud) UrlEntityUpdate() string
```

Returns the base URL for the entity update page. Append `&entity_id=ID` for the full URL.

### UrlEntityUpdateAjax

```go
func (crud *Crud) UrlEntityUpdateAjax() string
```

Returns the URL for the update AJAX endpoint (POST).

### UrlEntityTrashAjax

```go
func (crud *Crud) UrlEntityTrashAjax() string
```

Returns the URL for the trash AJAX endpoint (POST).

## Internal Method

### urlHome

```go
func (crud *Crud) urlHome() string
```

Returns the `HomeURL` configured in the `Config`. Used for the "Home" breadcrumb link. This method is **not exported**.

## Route Path Constants

Defined in `constants.go` (all unexported):

| Constant | Value | Used By |
|----------|-------|---------|
| `pathHome` | `"home"` | Default route, entity manager |
| `pathEntityCreateAjax` | `"entity-create-ajax"` | Create save endpoint |
| `pathEntityCreateModal` | `"entity-create-modal"` | Create modal fragment |
| `pathEntityManager` | `"entity-manager"` | Entity listing page |
| `pathEntityRead` | `"entity-read"` | Read view page |
| `pathEntityUpdate` | `"entity-update"` | Update form page |
| `pathEntityUpdateAjax` | `"entity-update-ajax"` | Update save endpoint |
| `pathEntityTrashAjax` | `"entity-trash-ajax"` | Trash endpoint |

## Usage in Controllers

Controllers use URL helpers to generate links for navigation, AJAX endpoints, and HTMX targets:

```go
// HTMX create modal button
hb.Button().HxGet(controller.crud.UrlEntityCreateModal())

// View link
hb.Hyperlink().Href(controller.crud.UrlEntityRead() + "&entity_id=" + row.ID)

// Edit link
hb.Hyperlink().Href(controller.crud.UrlEntityUpdate() + "&entity_id=" + row.ID)

// AJAX endpoint for JavaScript
urlEntityTrashAjax, _ := json.Marshal(controller.crud.UrlEntityTrashAjax())
```

## See Also

- [Modules: Crud Core](crud_core.md) - URL helper implementations
- [Modules: Controllers](controllers.md) - How URLs are used in controllers
- [API Reference](../api_reference.md) - Complete method reference
