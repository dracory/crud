# Layout and Templates

## Overview

The `Crud` struct provides layout rendering, breadcrumb generation, and a default webpage template. These are used internally by all controllers to produce consistent HTML pages.

## Layout

**Method:** `layout(w, r, title, content, styleFiles, style, jsFiles, js) string`

This is the central rendering method. It produces the final HTML string for a page.

### With Custom Layout (`FuncLayout` configured)

When `FuncLayout` is provided in the `Config`, the layout method:

1. Prepends required CDN resources to the provided `jsFiles` and `styleFiles`:
   - `jsFiles`: HTMX 2.0.0, Sweetalert2 11, Vue.js 3, Element Plus JS 2.3.8
   - `styleFiles`: Element Plus CSS 2.3.8
2. Delegates to `FuncLayout(w, r, title, content, styleFiles, style, jsFiles, js)`.

This allows the consumer to wrap CRUD pages in their own application shell (navigation, sidebar, footer, etc.) while the package ensures its required frontend dependencies are included.

### Without Custom Layout (default)

When `FuncLayout` is `nil`, the layout method:

1. Calls `webpage(title, content)` to create a standalone HTML page.
2. Appends any additional `styleFiles`, `style`, `jsFiles`, and `js` to the webpage.
3. Returns the full HTML document as a string.

## Default Webpage Template

**Method:** `webpage(title, content) *hb.HtmlWebpage`

Generates a complete HTML document with:

- **Favicon**: Embedded base64 icon
- **CSS**: Bootstrap 5.3.3 (CDN)
- **JavaScript**: Bootstrap 5.3.3, jQuery 3.7.1, Vue.js 3, Sweetalert2 11 (all CDN)
- **Base styles**:
  - Full-height `html` and `body`
  - Nunito font family
  - 0.9rem font size
  - Light gray background (`#f8fafc`)
  - Custom `.form-select` styling

The content is injected as raw HTML into the page body.

## Breadcrumbs

**Method:** `renderBreadcrumbs(breadcrumbs []Breadcrumb) string`

Generates Bootstrap-compatible breadcrumb navigation HTML.

### Breadcrumb Struct

```go
type Breadcrumb struct {
    Name string // Display text
    URL  string // Link href
}
```

### Usage

Each controller builds its own breadcrumb trail:

- **Entity Manager**: Home → Entity Manager
- **Entity Read**: Home → Entity Manager → View Entity
- **Entity Update**: Home → Entity Manager → Edit Entity

### Output

```html
<nav aria-label="breadcrumb">
  <ol class="breadcrumb">
    <li class="breadcrumb-item"><a href="/">Home</a></li>
    <li class="breadcrumb-item"><a href="/admin?path=entity-manager">Product Manager</a></li>
  </ol>
</nav>
```

## Form Rendering

**Method:** `form(fields []form.FieldInterface) []hb.TagInterface`

Generates Bootstrap form groups from field definitions. Each field produces:

1. A `.form-group.mb-3` container `<div>`
2. A `.form-label` `<label>` (with red asterisk for required fields)
3. An input element appropriate to the field type (see [Form Fields](form-fields.md))
4. An optional `.text-info` help paragraph

Fields with empty IDs get auto-generated random IDs. Fields with empty labels fall back to the field name.

## Helper Methods

### `listCreateNames() []string`

Returns a list of field names from `createFields`, skipping fields with empty names. Used by the create controller to extract form values from the request.

### `listUpdateNames() []string`

Returns a list of field names from `updateFields`, skipping fields with empty names. Used by the update controller to extract form values from the request.
