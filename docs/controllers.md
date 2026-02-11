# Controllers and Endpoints

## Routing

All requests are handled by `Crud.Handler(w, r)`. The `path` query parameter determines which controller handles the request. If no `path` is provided, the entity manager page is shown.

```
GET /admin                          → Entity Manager (default)
GET /admin?path=home                → Entity Manager
GET /admin?path=entity-manager      → Entity Manager
GET /admin?path=entity-create-modal → Create Modal (HTMX fragment)
POST /admin?path=entity-create-ajax → Create Save (AJAX)
GET /admin?path=entity-read&entity_id=ID    → Read View
GET /admin?path=entity-update&entity_id=ID  → Update Form
POST /admin?path=entity-update-ajax         → Update Save (AJAX)
POST /admin?path=entity-trash-ajax          → Trash (AJAX)
```

## Entity Manager Controller

**File:** `entity_manager_controller.go`

**Route:** `path=entity-manager` or `path=home` or no path

Renders the main entity listing page with:

- Breadcrumb navigation (Home → Entity Manager)
- "New Entity" button that loads the create modal via HTMX
- A DataTable with sortable columns
- Per-row action buttons (View, Edit, Trash) conditionally shown based on which callback functions are configured:
  - **View** button: shown when `FuncFetchReadData` is not `nil`
  - **Edit** button: shown when `FuncFetchUpdateData` is not `nil`
  - **Trash** button: shown when `FuncTrash` is not `nil`
- A Vue.js 3 app (`EntityManager`) mounted on `#entity-manager` that handles:
  - DataTable initialization
  - Create modal display
  - Trash confirmation modal with AJAX POST

**Response:** Full HTML page with `Content-Type: text/html`.

### Raw HTML in Columns

Column names and cell values can contain raw HTML by wrapping them in `{!! !!}`:

```go
ColumnNames: []string{"ID", "{!!Status!!}"}
```

When a column name is wrapped in `{!! !!}`, the corresponding cell values are rendered as raw HTML instead of escaped text.

---

## Entity Create Controller

**File:** `entity_create_controller.go`

### Modal Show

**Route:** `GET path=entity-create-modal`

Returns an HTML fragment containing a Bootstrap 5 modal with:

- Form fields generated from `CreateFields` using `github.com/dracory/form`
- "Create & Edit" submit button using HTMX (`hx-post` to the create AJAX endpoint)
- "Close" cancel button
- Modal backdrop

This endpoint is designed to be loaded via HTMX into the page body.

### Modal Save

**Route:** `POST path=entity-create-ajax`

Processes entity creation:

1. **Method check**: Rejects non-POST requests with a Sweetalert2 error.
2. **Nil check**: Returns error if `FuncCreate` is not configured.
3. **Field extraction**: Reads form values for all `CreateFields` by name.
4. **Required field validation**: Checks that all required fields have non-empty values.
5. **Creation**: Calls `FuncCreate` with the collected data.
6. **Success response**: Returns a Sweetalert2 success message and a JavaScript redirect to the update page for the new entity (after 2 seconds).

**Response:** HTML fragment with Sweetalert2 script tags (not JSON).

---

## Entity Read Controller

**File:** `entity_read_controller.go`

**Route:** `GET path=entity-read&entity_id=ID`

Renders a read-only view of an entity:

1. **Validation**: Requires `entity_id` parameter and `FuncFetchReadData` to be configured.
2. **Data fetch**: Calls `FuncFetchReadData(entityID)` to get key-value pairs.
3. **Rendering**: Displays data in a striped table within a Bootstrap card.
4. **Raw HTML support**: Keys and values wrapped in `{!! !!}` are rendered as raw HTML.
5. **Extras**: If `FuncReadExtras` is configured, appends additional HTML elements below the card.
6. **Navigation**: Edit button and Back button in the heading.

**Response:** Full HTML page with `Content-Type: text/html`.

**Error handling:** If `FuncFetchReadData` returns an error, an alert is shown instead of the table (the page still renders with status 200).

---

## Entity Update Controller

**File:** `entity_update_controller.go`

### Update Page

**Route:** `GET path=entity-update&entity_id=ID`

Renders an edit form for an entity:

1. **Validation**: Requires `entity_id` and `FuncFetchUpdateData` to be configured.
2. **Data fetch**: Calls `FuncFetchUpdateData(entityID)` to get current field values.
3. **Form rendering**: Generates form fields from `UpdateFields` with Vue.js `v-model` bindings.
4. **Vue.js app**: Mounts `EntityUpdate` on `#entity-update` with:
   - Two-way data binding for all fields
   - "Save" button (saves and redirects to manager)
   - "Apply" button (saves without redirect)
   - Image upload via FileReader API
   - Trumbowyg WYSIWYG editor integration
   - Element Plus date picker integration

**Response:** Full HTML page with `Content-Type: text/html`.

### Update Save

**Route:** `POST path=entity-update-ajax`

Processes entity update:

1. **Method check**: Rejects non-POST requests.
2. **Validation**: Requires `entity_id` parameter.
3. **Field extraction**: Reads form values for all `UpdateFields` by name.
4. **Required field validation**: Checks that all required fields have non-empty values.
5. **Update**: Calls `FuncUpdate(entityID, data)`.
6. **Response**: JSON via `api.Respond` with success/error status.

**Response:** JSON (`api.Response` format).

---

## Entity Trash Controller

**File:** `entity_trash_controller.go`

### Trash AJAX

**Route:** `POST path=entity-trash-ajax`

Processes entity soft-deletion:

1. **Method check**: Rejects non-POST requests.
2. **Validation**: Requires `entity_id` parameter.
3. **Nil check**: Returns error if `FuncTrash` is not configured.
4. **Trash**: Calls `FuncTrash(entityID)`.
5. **Response**: JSON via `api.Respond` with success/error status.

**Response:** JSON (`api.Response` format).

### Trash Modal

The trash controller also provides `pageEntitiesEntityTrashModal()` which generates a Bootstrap modal HTML fragment for the trash confirmation dialog. This modal is embedded directly in the entity manager page and controlled by Vue.js.

---

## Response Formats

The controllers use two different response formats:

| Controller | Endpoint | Format |
|------------|----------|--------|
| Entity Manager | `page` | Full HTML page |
| Entity Create | `modalShow` | HTML fragment (HTMX) |
| Entity Create | `modalSave` | HTML fragment with Sweetalert2 scripts |
| Entity Read | `page` | Full HTML page |
| Entity Update | `page` | Full HTML page |
| Entity Update | `pageSave` | JSON (`api.Response`) |
| Entity Trash | `pageEntityTrashAjax` | JSON (`api.Response`) |

### JSON Response Format

Endpoints using `api.Respond` return:

```json
{
    "status": "success",
    "message": "Saved successfully",
    "data": {
        "entity_id": "abc-123"
    }
}
```

Or on error:

```json
{
    "status": "error",
    "message": "Entity ID is required",
    "data": null
}
```
