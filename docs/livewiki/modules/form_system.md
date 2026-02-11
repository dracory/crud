---
path: modules/form_system.md
page-type: module
summary: Documentation for the form system including field types, options, rendering, validation, and Vue.js integration.
tags: [module, forms, fields, validation, vue]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: Form System

## Purpose

The form system generates Bootstrap form groups from field definitions. It supports 11 field types, required field validation, Vue.js two-way binding, and specialized widgets (WYSIWYG editor, date picker, block editor, image upload).

**Primary file:** `v2/crud.go` (the `form()` method)  
**Constants:** `v2/constants.go`  
**External dependency:** `github.com/dracory/form`

## Field Types

| Constant | Value | HTML Output |
|----------|-------|-------------|
| `FORM_FIELD_TYPE_STRING` | `"string"` | `<input type="text" class="form-control" v-model="entityModel.name">` |
| `FORM_FIELD_TYPE_NUMBER` | `"number"` | `<input type="number" class="form-control" v-model="entityModel.name">` |
| `FORM_FIELD_TYPE_PASSWORD` | `"password"` | `<input type="password" class="form-control" v-model="entityModel.name">` |
| `FORM_FIELD_TYPE_TEXTAREA` | `"textarea"` | `<textarea class="form-control" v-model="entityModel.name">` |
| `FORM_FIELD_TYPE_SELECT` | `"select"` | `<select class="form-select" v-model="entityModel.name">` with `<option>` elements |
| `FORM_FIELD_TYPE_IMAGE` | `"image"` | Image preview + URL input + optional "Browse" link |
| `FORM_FIELD_TYPE_IMAGE_INLINE` | `"image_inline"` | Image preview + file upload input (base64 via FileReader) |
| `FORM_FIELD_TYPE_HTMLAREA` | `"htmlarea"` | Trumbowyg WYSIWYG `<trumbowyg>` component |
| `FORM_FIELD_TYPE_BLOCKAREA` | `"blockarea"` | `<textarea>` + BlockArea script initialization |
| `FORM_FIELD_TYPE_DATETIME` | `"datetime"` | Element Plus `<el-date-picker type="datetime">` |
| `FORM_FIELD_TYPE_RAW` | `"raw"` | Raw HTML from field's `Value`. No label rendered. |

## Creating Fields

Fields use `form.NewField(form.FieldOptions{...})` from `github.com/dracory/form`:

```go
import "github.com/dracory/form"
import crud "github.com/dracory/crud/v2"

fields := []form.FieldInterface{
    form.NewField(form.FieldOptions{
        Name:     "title",
        Label:    "Title",
        Type:     crud.FORM_FIELD_TYPE_STRING,
        Required: true,
        Help:     "Enter the product title",
    }),
    form.NewField(form.FieldOptions{
        Name:  "status",
        Label: "Status",
        Type:  crud.FORM_FIELD_TYPE_SELECT,
        Options: []form.FieldOption{
            {Key: "active", Value: "Active"},
            {Key: "inactive", Value: "Inactive"},
        },
    }),
}
```

## Field Options

| Option | Type | Description |
|--------|------|-------------|
| `ID` | `string` | HTML element ID. Auto-generated (32-char random) if empty. |
| `Name` | `string` | Field name. Used as form key and `v-model` property. Fields with empty names are skipped. |
| `Label` | `string` | Display label. Falls back to `Name` if empty. |
| `Type` | `string` | One of the `FORM_FIELD_TYPE_*` constants. |
| `Value` | `string` | Default value. For `RAW` type, this is the raw HTML content. |
| `Required` | `bool` | Shows red asterisk; validated as non-empty on save. |
| `Help` | `string` | Help text below the field (supports HTML). |
| `Options` | `[]form.FieldOption` | Static options for `SELECT` fields. |
| `OptionsF` | `func() []form.FieldOption` | Dynamic options function for `SELECT` fields. |
| `Readonly` | `bool` | Marks field as read-only. |
| `Disabled` | `bool` | Marks field as disabled. |
| `Placeholder` | `string` | Placeholder text. |
| `Invisible` | `bool` | Hides the field. |

## Form Rendering

The `form()` method on `Crud` generates an array of `hb.TagInterface` elements. Each field produces:

1. A `.form-group.mb-3` container `<div>`
2. A `.form-label` `<label>` (with red `<sup>*</sup>` for required fields)
3. An input element appropriate to the field type
4. An optional `.text-info` help paragraph

### Special Field Behaviors

**IMAGE field:**
- Image preview with `v-bind:src` (placeholder if empty)
- URL text input with `v-model`
- "Browse" link to `FileManagerURL` (only if configured)

**IMAGE_INLINE field:**
- Image preview with `v-bind:src`
- File upload input with `v-on:change="uploadImage($event, 'fieldName')"`
- "See Image Data" toggle button
- Hidden textarea showing base64 data

**BLOCKAREA field:**
- Textarea with `v-model`
- Deferred `<script>` that initializes BlockArea with registered blocks (Heading, Text, Image, Code, RawHtml)

**RAW field:**
- Outputs the field's `Value` as raw HTML
- No label is rendered
- No `v-model` binding

## Vue.js Integration

All form fields (except `RAW`) use Vue.js `v-model` binding:

```
v-model="entityModel.<fieldName>"
```

This means field values are reactive and synchronized with the Vue.js app data:
- **Create modal:** Initialized with default values from field definitions
- **Update form:** Initialized with values from `FuncFetchUpdateData`

## Required Field Validation

Required fields are validated server-side during save operations (create and update):

1. Check `field.GetRequired() == true`
2. Check field name exists in posted data
3. Check value is not empty (using `lo.IsEmpty`)

If validation fails: `api.Error("<Label> is required field")`

## Select Field Options

Select fields support two option sources:

```go
// Static options
Options: []form.FieldOption{
    {Key: "active", Value: "Active"},
}

// Dynamic options (called at render time)
OptionsF: func() []form.FieldOption {
    return fetchFromDB()
}
```

Both can be used together - static options render first, then dynamic options are appended.

## See Also

- [Modules: Controllers](controllers.md) - How forms are used in controllers
- [Modules: Crud Core](crud_core.md) - The form() method implementation
- [Configuration](../configuration.md) - CreateFields and UpdateFields config
- [Cheatsheet](../cheatsheet.md) - Quick field creation examples
