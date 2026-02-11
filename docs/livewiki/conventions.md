---
path: conventions.md
page-type: reference
summary: Coding and documentation conventions used throughout the Dracory CRUD project.
tags: [conventions, code-style, patterns, guidelines]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Conventions

## Go Code Conventions

### Naming

- **Package name:** `crud` (both v1 and v2)
- **Exported types:** PascalCase (`Crud`, `Config`, `Row`, `KeyValue`, `Breadcrumb`)
- **Exported methods:** PascalCase (`Handler`, `UrlEntityManager`, `New`)
- **Unexported types:** camelCase (`entityCreateController`, `entityManagerController`)
- **Unexported methods:** camelCase (`page`, `modalShow`, `modalSave`, `getRoute`)
- **Constants (exported):** UPPER_SNAKE_CASE (`FORM_FIELD_TYPE_STRING`, `LogLevelInfo`)
- **Constants (unexported):** camelCase (`pathHome`, `pathEntityManager`)

### File Organization

| Pattern | Convention |
|---------|-----------|
| One type per file | `breadcrumb.go`, `key_value.go`, `row.go` |
| Controller per file | `entity_create_controller.go`, `entity_manager_controller.go` |
| Test file mirrors source | `entity_create_controller_test.go` matches `entity_create_controller.go` |
| Config separate from struct | `config.go` (public Config) vs `crud.go` (private Crud fields) |

### Controller Pattern

Every controller follows this structure:

```go
// 1. Private struct with pointer to parent Crud
type entityFooController struct {
    crud *Crud
}

// 2. Factory method on Crud
func (crud *Crud) newEntityFooController() *entityFooController {
    return &entityFooController{crud: crud}
}

// 3. Handler methods matching http.HandlerFunc signature
func (controller *entityFooController) page(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Error Handling

- Constructor (`New()`) returns `(Crud, error)` for validation failures
- Controllers use `api.Respond(w, r, api.Error("message"))` for JSON error responses
- Page controllers show inline alerts for data fetch errors (still return 200)
- All optional callbacks are nil-checked before invocation

### Security Patterns

- **JSON encoding for inline JS:** Always use `json.Marshal` before interpolating values into JavaScript
- **POST enforcement:** All mutating endpoints check `r.Method != http.MethodPost`
- **Nil safety:** Check `funcCallback != nil` before calling any optional callback

## HTML Generation Conventions

### Builder Pattern

HTML is generated using the `hb` (HTML builder) package with method chaining:

```go
hb.Div().
    Class("container").
    ID("entity-manager").
    Child(hb.Heading1().Text("Title")).
    Child(hb.Raw(content))
```

### Bootstrap Classes

- Containers: `.container`
- Tables: `.table .table-responsive .table-striped`
- Buttons: `.btn .btn-{variant}` (success, primary, secondary, danger, outline-*)
- Cards: `.card .card-header .card-body`
- Forms: `.form-group .mb-3`, `.form-label`, `.form-control`, `.form-select`
- Alerts: `.alert .alert-danger`

### Vue.js Integration

- Each page mounts a Vue app on a specific container element
- Form fields use `v-model="entityModel.<fieldName>"` binding
- Event handlers use `v-on:click` attribute
- Vue apps are created with `Vue.createApp(AppConfig).mount('#container')`

## Testing Conventions

### Test Naming

```
TestNew_MissingFuncRows          // Constructor validation
TestNew_ValidFullConfig          // Constructor success
TestEntityManagerController_Page // Controller behavior
```

### Test Structure

- Use `testing.T` (no external test frameworks)
- Use `httptest.NewRecorder()` and `httptest.NewRequest()` for HTTP tests
- Assert with `t.Fatal()` / `t.Fatalf()` for failures
- Create minimal `Config` for each test (only required fields)

## Documentation Conventions

### File Naming

- Lowercase with underscores: `getting_started.md`, `api_reference.md`
- Hyphens for compound names in existing docs: `form-fields.md`, `url-helpers.md`
- Module docs in `modules/` subfolder: `modules/crud_core.md`

### Metadata Header

Every LiveWiki page includes frontmatter:

```markdown
---
path: filename.md
page-type: overview | reference | tutorial | module | changelog
summary: One-line description.
tags: [tag1, tag2]
created: YYYY-MM-DD
updated: YYYY-MM-DD
version: X.Y.Z
---
```

### Cross-References

Every page ends with a "See Also" section linking to related pages.

## Version Conventions

- **v2** is the current recommended version
- **v1** is legacy and unmaintained
- Each version is a separate Go module with its own `go.mod`
- Documentation covers v2 by default; v1 is referenced only for migration context

## See Also

- [Development](development.md) - Development workflow
- [Architecture](architecture.md) - Design patterns
- [LLM Context](llm-context.md) - Codebase summary for AI
