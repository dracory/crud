---
path: getting_started.md
page-type: tutorial
summary: Installation, setup, and quick start guide for the Dracory CRUD v2 package.
tags: [getting-started, installation, tutorial, quick-start]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Getting Started

## Prerequisites

- **Go 1.25+** (as specified in `go.mod`)
- A Go web application using `net/http` or a compatible router

## Installation

```bash
go get github.com/dracory/crud/v2
```

## Import

```go
import crud "github.com/dracory/crud/v2"
```

> **Note:** The import alias `crud` is recommended since the module path ends with `v2`.

## Quick Start

Here is a minimal working example that sets up a CRUD interface for managing products:

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
			return map[string]string{"name": "Widget", "status": "active"}, nil
		},
		FuncUpdate: func(r *http.Request, entityID string, data map[string]string) error {
			return nil
		},
		FuncTrash: func(r *http.Request, entityID string) error {
			return nil
		},
		FuncFetchReadData: func(r *http.Request, entityID string) ([]crud.KeyValue, error) {
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

## Step-by-Step Walkthrough

### 1. Define Your Entity Fields

Create field definitions for both create and update forms:

```go
fields := []form.FieldInterface{
    form.NewField(form.FieldOptions{
        Name:     "title",
        Label:    "Title",
        Type:     crud.FORM_FIELD_TYPE_STRING,
        Required: true,
    }),
    form.NewField(form.FieldOptions{
        Name:  "description",
        Label: "Description",
        Type:  crud.FORM_FIELD_TYPE_TEXTAREA,
    }),
}
```

### 2. Implement Callback Functions

At minimum, you need `FuncRows` and `UpdateFields`:

```go
funcRows := func(r *http.Request) ([]crud.Row, error) {
    // Query your database and return rows
    return rows, nil
}
```

### 3. Create the Crud Instance

```go
c, err := crud.New(crud.Config{
    Endpoint:     "/admin/products",
    FuncRows:     funcRows,
    UpdateFields: fields,
    // ... other config
})
```

### 4. Register the Handler

```go
http.HandleFunc("/admin/products", c.Handler)
```

### 5. Access the Interface

Navigate to `http://localhost:8080/admin/products` to see the entity manager.

## Minimum Required Configuration

The absolute minimum configuration requires only two fields:

```go
c, err := crud.New(crud.Config{
    FuncRows: func(r *http.Request) ([]crud.Row, error) {
        return []crud.Row{}, nil
    },
    UpdateFields: []form.FieldInterface{}, // Can be empty slice
})
```

This creates a read-only entity manager with no create, update, or trash capabilities.

## Adding Features Incrementally

| Feature | Required Config |
|---------|----------------|
| List entities | `FuncRows`, `ColumnNames` |
| Create entities | `FuncCreate`, `CreateFields` |
| View entities | `FuncFetchReadData` |
| Edit entities | `FuncUpdate`, `FuncFetchUpdateData`, `UpdateFields` |
| Trash entities | `FuncTrash` |
| Custom layout | `FuncLayout` |
| Pagination | `PageSize`, `FuncRowsCount` |
| Authorization | `FuncBeforeAction` |
| CSRF protection | `FuncValidateCSRF` |
| Logging | `FuncLog` |

## See Also

- [Configuration](configuration.md) - Full Config struct reference
- [API Reference](api_reference.md) - Complete API documentation
- [Architecture](architecture.md) - System design and patterns
- [Cheatsheet](cheatsheet.md) - Quick reference for common operations
