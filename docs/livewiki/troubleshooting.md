---
path: troubleshooting.md
page-type: reference
summary: Common issues, error messages, and solutions when using the Dracory CRUD package.
tags: [troubleshooting, errors, debugging, faq]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Troubleshooting

## Constructor Errors

### "FuncRows function is required"

**Cause:** `Config.FuncRows` is `nil`.

**Solution:** Provide a `FuncRows` callback that returns entity rows:

```go
crud.New(crud.Config{
    FuncRows: func(r *http.Request) ([]crud.Row, error) {
        return []crud.Row{}, nil
    },
    UpdateFields: []form.FieldInterface{},
})
```

### "UpdateFields is required"

**Cause:** `Config.UpdateFields` is `nil`.

**Solution:** Provide `UpdateFields`, even if empty:

```go
UpdateFields: []form.FieldInterface{},
```

### "FuncUpdate function is required when UpdateFields are provided"

**Cause:** `UpdateFields` has entries but `FuncUpdate` is `nil`.

**Solution:** Provide a `FuncUpdate` callback:

```go
FuncUpdate: func(r *http.Request, entityID string, data map[string]string) error {
    // Update entity in database
    return nil
},
```

### "FuncFetchUpdateData function is required when UpdateFields are provided"

**Cause:** `UpdateFields` has entries but `FuncFetchUpdateData` is `nil`.

**Solution:** Provide a `FuncFetchUpdateData` callback:

```go
FuncFetchUpdateData: func(r *http.Request, entityID string) (map[string]string, error) {
    return map[string]string{"name": "value"}, nil
},
```

## Runtime Errors

### "Entity ID is required"

**Cause:** The `entity_id` query/form parameter is missing or empty.

**Affected endpoints:** `entity-read`, `entity-update`, `entity-update-ajax`, `entity-trash-ajax`

**Solution:** Ensure the `entity_id` parameter is included in the request URL or POST body.

### "Method not allowed"

**Cause:** A non-POST request was sent to a POST-only endpoint.

**Affected endpoints:** `entity-create-ajax`, `entity-update-ajax`, `entity-trash-ajax`

**Solution:** Ensure AJAX calls use `POST` method.

### "Create functionality is not configured"

**Cause:** `FuncCreate` is `nil` but the create AJAX endpoint was called.

**Solution:** Provide a `FuncCreate` callback in the config.

### "Trash functionality is not configured"

**Cause:** `FuncTrash` is `nil` but the trash AJAX endpoint was called.

**Solution:** Provide a `FuncTrash` callback in the config.

### "Update functionality is not configured"

**Cause:** `FuncFetchUpdateData` is `nil` but the update page was requested.

**Solution:** Provide a `FuncFetchUpdateData` callback in the config.

### "CSRF validation failed: ..."

**Cause:** `FuncValidateCSRF` returned a non-nil error for a POST request.

**Solution:** Ensure the CSRF token is included in POST requests, or check your `FuncValidateCSRF` implementation.

### "Save failed: ..."

**Cause:** `FuncCreate` or `FuncUpdate` returned an error.

**Solution:** Check the error message for details from your callback implementation.

## UI Issues

### DataTable not initializing

**Possible causes:**
- jQuery or DataTables CDN is blocked
- JavaScript error in the page preventing Vue.js app from mounting

**Solution:** Check browser console for errors. Ensure CDN resources are accessible.

### Create modal not appearing

**Possible causes:**
- HTMX CDN is blocked (only loaded when custom `FuncLayout` is used)
- The `entity-create-modal` endpoint is returning an error

**Solution:** Check the Network tab in browser DevTools for the HTMX request. Verify `CreateFields` are configured.

### Form fields not binding

**Possible causes:**
- Vue.js CDN is blocked
- Field names contain special characters

**Solution:** Ensure field names are valid JavaScript identifiers. Check browser console for Vue.js errors.

### Pagination not showing

**Possible causes:**
- `PageSize` is `0` (default)
- `FuncRowsCount` is `nil`
- Total rows fit in one page

**Solution:** Set `PageSize > 0` and provide `FuncRowsCount`. Ensure your `FuncRows` respects the `page` query parameter.

### Action buttons not showing

**Possible causes:**
- **View button missing:** `FuncFetchReadData` is `nil`
- **Edit button missing:** `FuncFetchUpdateData` is `nil`
- **Trash button missing:** `FuncTrash` is `nil`

**Solution:** Provide the corresponding callback function to enable each button.

## Layout Issues

### CDN resources not loading

**Cause:** The default layout loads all frontend libraries from CDN. If CDN access is restricted, the page will not render correctly.

**Solution:** Provide a custom `FuncLayout` that serves resources from your own infrastructure.

### Custom layout missing required JS libraries

**Cause:** When using `FuncLayout`, the package prepends required CDN URLs to `jsFiles` and `styleFiles`. If your layout ignores these parameters, interactivity will break.

**Solution:** Ensure your `FuncLayout` includes all provided `jsFiles` and `styleFiles` in the rendered HTML.

## See Also

- [Configuration](configuration.md) - Config struct reference
- [Getting Started](getting_started.md) - Setup guide
- [Development](development.md) - Testing and debugging
