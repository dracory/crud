# URL Helpers

## Overview

The `Crud` struct provides public URL generation methods that build endpoint URLs with the correct `path` query parameter. These are used internally by controllers and can also be used by consumers to generate links.

All URL helpers automatically detect whether the base `Endpoint` already contains a query string (`?`) and append the `path` parameter with the correct separator (`?` or `&`).

## Methods

### `UrlEntityManager() string`

Returns the URL for the entity manager page.

```go
// Endpoint = "/admin"
crud.UrlEntityManager() // "/admin?path=entity-manager"

// Endpoint = "/admin?module=cms"
crud.UrlEntityManager() // "/admin?module=cms&path=entity-manager"
```

### `UrlEntityCreateModal() string`

Returns the URL for the create modal HTML fragment (loaded via HTMX).

```go
crud.UrlEntityCreateModal() // "/admin?path=entity-create-modal"
```

### `UrlEntityCreateAjax() string`

Returns the URL for the create AJAX endpoint (POST).

```go
crud.UrlEntityCreateAjax() // "/admin?path=entity-create-ajax"
```

### `UrlEntityRead() string`

Returns the base URL for the entity read page. Append `&entity_id=ID` to form the full URL.

```go
crud.UrlEntityRead() // "/admin?path=entity-read"
// Full: "/admin?path=entity-read&entity_id=abc-123"
```

### `UrlEntityUpdate() string`

Returns the base URL for the entity update page. Append `&entity_id=ID` to form the full URL.

```go
crud.UrlEntityUpdate() // "/admin?path=entity-update"
// Full: "/admin?path=entity-update&entity_id=abc-123"
```

### `UrlEntityUpdateAjax() string`

Returns the URL for the update AJAX endpoint (POST).

```go
crud.UrlEntityUpdateAjax() // "/admin?path=entity-update-ajax"
```

### `UrlEntityTrashAjax() string`

Returns the URL for the trash AJAX endpoint (POST).

```go
crud.UrlEntityTrashAjax() // "/admin?path=entity-trash-ajax"
```

## Internal Method

### `urlHome() string`

Returns the `HomeURL` configured in the `Config`. Used for the "Home" breadcrumb link. This method is not exported.

## Route Path Constants

The following internal constants define the `path` values (defined in `constants.go`):

| Constant | Value |
|----------|-------|
| `pathHome` | `"home"` |
| `pathEntityCreateAjax` | `"entity-create-ajax"` |
| `pathEntityCreateModal` | `"entity-create-modal"` |
| `pathEntityManager` | `"entity-manager"` |
| `pathEntityRead` | `"entity-read"` |
| `pathEntityUpdate` | `"entity-update"` |
| `pathEntityUpdateAjax` | `"entity-update-ajax"` |
| `pathEntityTrashAjax` | `"entity-trash-ajax"` |
