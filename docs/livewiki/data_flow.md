---
path: data_flow.md
page-type: overview
summary: How data moves through the CRUD system from HTTP requests to rendered HTML and JSON responses.
tags: [data-flow, request-lifecycle, callbacks, rendering]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Data Flow

## Overview

Data in the CRUD package flows through a well-defined pipeline: HTTP request → middleware → controller → callback functions → rendering → HTTP response. The package never accesses a database directly; all data operations are delegated to user-provided callback functions.

## Entity Manager (List) Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Handler
    participant ManagerCtrl as entityManagerController
    participant FuncRows
    participant FuncRowsCount
    participant Layout

    Browser->>Handler: GET /admin?path=entity-manager
    Handler->>ManagerCtrl: page(w, r)
    ManagerCtrl->>FuncRows: funcRows(r)
    FuncRows-->>ManagerCtrl: []Row, error
    alt PageSize > 0
        ManagerCtrl->>FuncRowsCount: funcRowsCount(r)
        FuncRowsCount-->>ManagerCtrl: totalRows, error
        ManagerCtrl->>ManagerCtrl: renderPagination()
    end
    ManagerCtrl->>Layout: layout(w, r, title, content, ...)
    Layout-->>Browser: Full HTML page
```

**Data transformations:**
1. `FuncRows` returns `[]Row` (each row has `ID` and `Data []string`)
2. `Data` values are mapped to table cells matching `ColumnNames`
3. Column names/values wrapped in `{!! !!}` are rendered as raw HTML
4. Vue.js app data is JSON-encoded via `json.Marshal` for XSS safety

## Entity Create Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Handler
    participant CreateCtrl as entityCreateController
    participant FuncCreate

    Note over Browser: Step 1: Load modal via HTMX
    Browser->>Handler: GET /admin?path=entity-create-modal
    Handler->>CreateCtrl: modalShow(w, r)
    CreateCtrl-->>Browser: HTML fragment (Bootstrap modal)

    Note over Browser: Step 2: Submit form
    Browser->>Handler: POST /admin?path=entity-create-ajax
    Handler->>CreateCtrl: modalSave(w, r)
    CreateCtrl->>CreateCtrl: Extract field values from request
    CreateCtrl->>CreateCtrl: Validate required fields
    CreateCtrl->>FuncCreate: funcCreate(r, data)
    FuncCreate-->>CreateCtrl: entityID, error
    CreateCtrl-->>Browser: JSON {status, entity_id, redirect_url}
```

**Data transformations:**
1. Form field names are extracted from `CreateFields` via `listCreateNames()`
2. Values are read from the POST body using `req.GetString(r, name)`
3. Required fields are validated (non-empty check)
4. `FuncCreate` receives `map[string]string` and returns the new entity ID
5. Response includes `redirect_url` (defaults to update page for the new entity)

## Entity Read Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Handler
    participant ReadCtrl as entityReadController
    participant FuncFetchReadData
    participant FuncReadExtras
    participant Layout

    Browser->>Handler: GET /admin?path=entity-read&entity_id=ID
    Handler->>ReadCtrl: page(w, r)
    ReadCtrl->>FuncFetchReadData: funcFetchReadData(r, entityID)
    FuncFetchReadData-->>ReadCtrl: []KeyValue, error
    alt FuncReadExtras configured
        ReadCtrl->>FuncReadExtras: funcReadExtras(r, entityID)
        FuncReadExtras-->>ReadCtrl: []hb.TagInterface
    end
    ReadCtrl->>Layout: layout(w, r, title, content, ...)
    Layout-->>Browser: Full HTML page
```

**Data transformations:**
1. `entity_id` is extracted from query parameters
2. `FuncFetchReadData` returns `[]KeyValue` pairs
3. Keys and values wrapped in `{!! !!}` are rendered as raw HTML
4. Optional `FuncReadExtras` appends additional HTML elements below the card

## Entity Update Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Handler
    participant UpdateCtrl as entityUpdateController
    participant FuncFetchUpdateData
    participant FuncUpdate
    participant Layout

    Note over Browser: Step 1: Load form
    Browser->>Handler: GET /admin?path=entity-update&entity_id=ID
    Handler->>UpdateCtrl: page(w, r)
    UpdateCtrl->>FuncFetchUpdateData: funcFetchUpdateData(r, entityID)
    FuncFetchUpdateData-->>UpdateCtrl: map[string]string, error
    UpdateCtrl->>Layout: layout(w, r, title, content, ...)
    Layout-->>Browser: Full HTML page with Vue.js form

    Note over Browser: Step 2: Save changes
    Browser->>Handler: POST /admin?path=entity-update-ajax
    Handler->>UpdateCtrl: pageSave(w, r)
    UpdateCtrl->>UpdateCtrl: Extract field values from request
    UpdateCtrl->>UpdateCtrl: Validate required fields
    UpdateCtrl->>FuncUpdate: funcUpdate(r, entityID, data)
    FuncUpdate-->>UpdateCtrl: error
    UpdateCtrl-->>Browser: JSON {status, entity_id}
```

**Data transformations:**
1. `FuncFetchUpdateData` returns `map[string]string` (field name → current value)
2. Values are JSON-encoded and injected into the Vue.js app as `customValues`
3. Vue.js binds values to form inputs via `v-model="entityModel.<fieldName>"`
4. On save, Vue.js sends all `entityModel` data via `$.post()`
5. Server extracts values using `listUpdateNames()` and `req.GetStringTrimmed()`

## Entity Trash Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Handler
    participant TrashCtrl as entityTrashController
    participant FuncTrash

    Browser->>Handler: POST /admin?path=entity-trash-ajax
    Handler->>TrashCtrl: pageEntityTrashAjax(w, r)
    TrashCtrl->>FuncTrash: funcTrash(r, entityID)
    FuncTrash-->>TrashCtrl: error
    TrashCtrl-->>Browser: JSON {status, entity_id}
```

**Data transformations:**
1. `entity_id` is extracted from the POST body
2. `FuncTrash` receives the entity ID and performs soft-deletion
3. On success, the browser reloads the page after a 3-second delay

## Middleware Data Flow

```mermaid
graph TD
    A[Request] --> B[FuncLog: log request info]
    B --> C{FuncBeforeAction}
    C -->|false| D[Abort - no response body]
    C -->|true| E{POST + FuncValidateCSRF?}
    E -->|CSRF error| F[JSON error response]
    E -->|pass| G[Controller execution]
    G --> H[FuncAfterAction]
```

The `action` parameter passed to `FuncBeforeAction` and `FuncAfterAction` is the route path string (e.g., `"entity-manager"`, `"entity-create-ajax"`).

## See Also

- [Architecture](architecture.md) - System design and patterns
- [Modules: Controllers](modules/controllers.md) - Detailed controller documentation
- [Configuration](configuration.md) - Callback function signatures
- [API Reference](api_reference.md) - Complete API documentation
