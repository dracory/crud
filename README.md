# CRUD

A server-side rendered CRUD (Create, Read, Update, Trash) interface for Go web applications. Generates Bootstrap 5 admin pages with Vue.js 3 interactivity for managing entities.

## Versions

### v2 (Current)

The recommended version with improved security, validation, and code quality.

```bash
go get github.com/dracory/crud/v2
```

```go
import crud "github.com/dracory/crud/v2"
```

Key improvements over v1:

- **Security** — XSS prevention via JSON-encoded inline JS values, HTTP method enforcement on mutating endpoints, nil safety on all callbacks
- **Architecture** — Controller-based routing, separated concerns across dedicated files
- **Form system** — Uses `github.com/dracory/form` package with 11 field types
- **HTMX integration** — Create modal loaded via HTMX instead of inline Vue.js
- **Comprehensive tests** — 56 tests covering all controllers, validation, routing, and rendering

See [`v2/README.md`](v2/README.md) for full documentation, quick start example, and API reference.

### v1 (Legacy)

The original implementation. Use v1 only if you have an existing project that depends on it.

```bash
go get github.com/dracory/crud
```

```go
import "github.com/dracory/crud"
```

See [`v1/README.md`](v1/README.md) for usage instructions.

> **Note:** v1 is no longer actively maintained. New projects should use v2.

## Documentation

Detailed documentation covering both versions is available in the [`docs/`](docs/) folder:

- [Overview](docs/overview.md) — Package architecture, core types, frontend stack, dependencies, and security model
- [Configuration](docs/configuration.md) — `Config` struct reference, `New()` validation rules, callback function signatures
- [Controllers](docs/controllers.md) — Endpoint routing, controller behavior, request/response formats
- [Form Fields](docs/form-fields.md) — Field type constants, options reference, Vue.js integration
- [URL Helpers](docs/url-helpers.md) — URL generation methods and route path constants
- [Layout and Templates](docs/layout-and-templates.md) — Layout rendering, breadcrumbs, default webpage template

## Project Structure

```
crud/
├── docs/          # Shared documentation
├── v1/            # Legacy version (original implementation)
│   ├── README.md
│   ├── go.mod
│   └── ...
├── v2/            # Current version (recommended)
│   ├── README.md
│   ├── go.mod
│   └── ...
├── .gitignore
└── README.md      # This file
```

## License

See the repository for license information.
