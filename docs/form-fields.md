# Form Fields

## Overview

Form fields define the inputs rendered in the create modal and update form. They use the `form.FieldInterface` from `github.com/dracory/form` and are configured via `form.NewField(form.FieldOptions{...})`.

## Creating Fields

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
        Name:  "description",
        Label: "Description",
        Type:  crud.FORM_FIELD_TYPE_TEXTAREA,
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

## Field Types

The following field type constants are defined in `constants.go`:

| Constant | Value | Rendered As |
|----------|-------|-------------|
| `FORM_FIELD_TYPE_STRING` | `"string"` | `<input type="text">` with `v-model` binding |
| `FORM_FIELD_TYPE_NUMBER` | `"number"` | `<input type="number">` with `v-model` binding |
| `FORM_FIELD_TYPE_PASSWORD` | `"password"` | `<input type="password">` with `v-model` binding |
| `FORM_FIELD_TYPE_TEXTAREA` | `"textarea"` | `<textarea>` with `v-model` binding |
| `FORM_FIELD_TYPE_SELECT` | `"select"` | `<select>` with `<option>` elements and `v-model` binding |
| `FORM_FIELD_TYPE_IMAGE` | `"image"` | Image preview + URL input + optional file manager "Browse" link |
| `FORM_FIELD_TYPE_IMAGE_INLINE` | `"image_inline"` | Image preview + file upload input (base64 via FileReader) + toggle for raw data view |
| `FORM_FIELD_TYPE_HTMLAREA` | `"htmlarea"` | Trumbowyg WYSIWYG editor with `v-model` binding |
| `FORM_FIELD_TYPE_BLOCKAREA` | `"blockarea"` | `<textarea>` + BlockArea script initialization for block-based editing |
| `FORM_FIELD_TYPE_DATETIME` | `"datetime"` | Element Plus `<el-date-picker>` with datetime type |
| `FORM_FIELD_TYPE_RAW` | `"raw"` | Raw HTML output from the field's `Value`. No label is rendered. |

## Field Options

The `form.FieldOptions` struct supports:

| Option | Type | Description |
|--------|------|-------------|
| `ID` | `string` | HTML element ID. Auto-generated if empty. |
| `Name` | `string` | Field name, used as the form key and `v-model` property. Fields with empty names are skipped in `listCreateNames`/`listUpdateNames`. |
| `Label` | `string` | Display label. Falls back to `Name` if empty. |
| `Type` | `string` | One of the `FORM_FIELD_TYPE_*` constants. |
| `Value` | `string` | Default value. For `RAW` type, this is the raw HTML content. |
| `Required` | `bool` | When `true`, a red asterisk is shown next to the label, and the field is validated as non-empty on save. |
| `Help` | `string` | Help text displayed below the field as a `text-info` paragraph. Supports HTML. |
| `Options` | `[]form.FieldOption` | Static options for `SELECT` fields. Each has a `Key` and `Value`. |
| `OptionsF` | `func() []form.FieldOption` | Dynamic options function for `SELECT` fields. Called at render time. |
| `Readonly` | `bool` | Marks the field as read-only. |
| `Disabled` | `bool` | Marks the field as disabled. |
| `Placeholder` | `string` | Placeholder text for input fields. |
| `Invisible` | `bool` | Hides the field. |

## Vue.js Integration

All form fields (except `RAW`) use Vue.js `v-model` binding with the pattern:

```
v-model="entityModel.<fieldName>"
```

This means the field values are reactive and automatically synchronized with the Vue.js app data. The update controller initializes `entityModel` with values from `FuncFetchUpdateData`, while the create modal initializes with default values from the field definitions.

## Required Field Validation

Required fields are validated server-side during save operations:

1. The field's `GetRequired()` returns `true`.
2. The server checks that the field name exists in the posted data.
3. The server checks that the value is not empty.

If validation fails, an error response is returned with the message `"<Label> is required field"`.

## Image Fields

### `FORM_FIELD_TYPE_IMAGE`

Renders:
- An `<img>` tag with `v-bind:src` showing the current image or a placeholder
- A URL input field for direct URL entry
- A "Browse" link to the file manager (only if `FileManagerURL` is configured)

### `FORM_FIELD_TYPE_IMAGE_INLINE`

Renders:
- An `<img>` tag with `v-bind:src` showing the current image or a placeholder
- A file upload input that converts the selected file to base64 via `FileReader`
- A "See Image Data" toggle button to show/hide the raw base64 data
- A `<textarea>` for viewing/editing the base64 data (shown when toggled)

The `uploadImage` method in the Vue.js update controller handles the file-to-base64 conversion.
