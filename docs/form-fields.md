# Form Fields

## Overview

Form fields define the inputs rendered in the create modal and update form. All field types, constructors, and the `FieldInterface` are built into this package — no external dependency is needed.

## Creating Fields

### Struct-based

```go
import crud "github.com/dracory/crud/v2"

fields := []crud.FieldInterface{
    crud.NewField(crud.FieldOptions{
        Name:     "title",
        Label:    "Title",
        Type:     crud.FORM_FIELD_TYPE_STRING,
        Required: true,
        Help:     "Enter the product title",
    }),
    crud.NewField(crud.FieldOptions{
        Name:  "description",
        Label: "Description",
        Type:  crud.FORM_FIELD_TYPE_TEXTAREA,
    }),
    crud.NewField(crud.FieldOptions{
        Name:  "status",
        Label: "Status",
        Type:  crud.FORM_FIELD_TYPE_SELECT,
        Options: []crud.FieldOption{
            {Key: "active", Value: "Active"},
            {Key: "inactive", Value: "Inactive"},
        },
    }),
}
```

### Fluent builder

The same fields using convenience constructors and fluent methods:

```go
fields := []crud.FieldInterface{
    crud.NewStringField("title", "Title").WithRequired().WithHelp("Enter the product title"),
    crud.NewTextAreaField("description", "Description"),
    crud.NewSelectField("status", "Status", []crud.FieldOption{
        {Key: "active", Value: "Active"},
        {Key: "inactive", Value: "Inactive"},
    }),
}
```

## Convenience Constructors

| Constructor | Field Type |
|-------------|-----------|
| `NewStringField(name, label)` | `FORM_FIELD_TYPE_STRING` |
| `NewNumberField(name, label)` | `FORM_FIELD_TYPE_NUMBER` |
| `NewPasswordField(name, label)` | `FORM_FIELD_TYPE_PASSWORD` |
| `NewEmailField(name, label)` | `FORM_FIELD_TYPE_EMAIL` |
| `NewTelField(name, label)` | `FORM_FIELD_TYPE_TEL` |
| `NewURLField(name, label)` | `FORM_FIELD_TYPE_URL` |
| `NewDateField(name, label)` | `FORM_FIELD_TYPE_DATE` |
| `NewDateTimeField(name, label)` | `FORM_FIELD_TYPE_DATETIME` |
| `NewColorField(name, label)` | `FORM_FIELD_TYPE_COLOR` |
| `NewCheckboxField(name, label)` | `FORM_FIELD_TYPE_CHECKBOX` |
| `NewRadioField(name, label, options)` | `FORM_FIELD_TYPE_RADIO` |
| `NewTextAreaField(name, label)` | `FORM_FIELD_TYPE_TEXTAREA` |
| `NewSelectField(name, label, options)` | `FORM_FIELD_TYPE_SELECT` |
| `NewImageField(name, label)` | `FORM_FIELD_TYPE_IMAGE` |
| `NewFileField(name, label)` | `FORM_FIELD_TYPE_FILE` |
| `NewHtmlAreaField(name, label)` | `FORM_FIELD_TYPE_HTMLAREA` |
| `NewHiddenField(name, value)` | `FORM_FIELD_TYPE_HIDDEN` |
| `NewRawField(value)` | `FORM_FIELD_TYPE_RAW` |
| `NewRepeater(RepeaterOptions{...})` | `FORM_FIELD_TYPE_REPEATER` |

## Fluent Methods

All methods return `*Field` for chaining:

| Method | Description |
|--------|-------------|
| `.WithID(id)` | Set HTML element ID |
| `.WithName(name)` | Set field name |
| `.WithLabel(label)` | Set label text |
| `.WithHelp(text)` | Set help text shown below the input |
| `.WithType(type)` | Set field type |
| `.WithValue(value)` | Set default value |
| `.WithRequired()` | Mark as required (shows asterisk, validated on save) |
| `.WithOptions(opts...)` | Set static options for select/radio fields |
| `.WithOptionsF(fn)` | Set dynamic options function (called at render time) |
| `.WithFields(fields...)` | Set sub-fields for repeater rows |
| `.WithFuncValuesProcess(fn)` | Set value pre-processor for repeater update form |

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

## FieldOptions Reference

| Option | Type | Description |
|--------|------|-------------|
| `ID` | `string` | HTML element ID. Auto-generated if empty. |
| `Name` | `string` | Field name, used as the form key and `v-model` property. Fields with empty names are skipped. |
| `Label` | `string` | Display label. Falls back to `Name` if empty. |
| `Type` | `string` | One of the `FORM_FIELD_TYPE_*` constants. |
| `Value` | `string` | Default value. For `RAW` type, this is the raw HTML content. |
| `Required` | `bool` | When `true`, a red asterisk is shown and the field is validated as non-empty on save. |
| `Help` | `string` | Help text displayed below the field as a `text-info` paragraph. Supports HTML. |
| `Options` | `[]FieldOption` | Static options for `SELECT` and `RADIO` fields. Each has a `Key` and `Value`. |
| `OptionsF` | `func() []FieldOption` | Dynamic options function for `SELECT` fields. Called at render time. |
| `Fields` | `[]FieldInterface` | Sub-fields for `REPEATER` type. |
| `FuncValuesProcess` | `func(raw string) []map[string]string` | Pre-processor for repeater values on the update form. |

## Vue.js Integration

All form fields (except `RAW`) use Vue.js `v-model` binding with the pattern:

```
v-model="entityModel.<fieldName>"
```

The update controller initializes `entityModel` with values from `FuncFetchUpdateData`, while the create modal initializes with default values from the field definitions.

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

## Repeater Fields

`FORM_FIELD_TYPE_REPEATER` renders a dynamic list of rows that the user can add, remove, and reorder.

### Generic (single-value rows)

```go
crud.NewField(crud.FieldOptions{
    Name:  "image_urls",
    Label: "Image URLs",
    Type:  crud.FORM_FIELD_TYPE_REPEATER,
})
```

Each row is a plain text input bound to `entityModel.image_urls[index]`.

### Typed sub-fields (structured rows)

```go
crud.NewRepeater(crud.RepeaterOptions{
    Name:  "links",
    Label: "Links",
    Fields: []crud.FieldInterface{
        crud.NewStringField("label", "Label"),
        crud.NewURLField("url", "URL"),
    },
})
```

All field types supported by the top-level form are also supported as sub-fields.

### Value pre-processing for the update form

When the update form loads, repeater values come back from `FuncFetchUpdateData` as a raw JSON string. Use `FuncValuesProcess` to decode them into the `[]map[string]string` structure Vue expects:

```go
crud.NewRepeater(crud.RepeaterOptions{
    Name:  "links",
    Label: "Links",
    Fields: []crud.FieldInterface{
        crud.NewStringField("label", "Label"),
        crud.NewURLField("url", "URL"),
    },
    FuncValuesProcess: func(raw string) []map[string]string {
        var rows []map[string]string
        json.Unmarshal([]byte(raw), &rows)
        return rows
    },
})
```

### Row controls

Each row renders three buttons:
- Up arrow — moves the row up (disabled on the first row)
- Down arrow — moves the row down (disabled on the last row)
- Trash — removes the row

These rely on four Vue methods present in all three Vue apps (`EntityManager`, `EntityCreate`, `EntityUpdate`):

```
addRepeaterItem(fieldName, item)
removeRepeaterItem(fieldName, index)
moveRepeaterItemUp(fieldName, index)
moveRepeaterItemDown(fieldName, index)
```

### How repeater data is submitted and received

The browser submits repeater rows using standard bracket notation:

```
links[0][label]=Home&links[0][url]=https://example.com
links[1][label]=About&links[1][url]=https://example.com/about
```

The server-side `collectRepeaterField(r, fieldName)` function parses these keys and returns a JSON array string:

```json
[{"label":"Home","url":"https://example.com"},{"label":"About","url":"https://example.com/about"}]
```

This JSON string is what arrives in `data[fieldName]` inside your `FuncCreate` / `FuncUpdate` callback:

```go
FuncUpdate: func(r *http.Request, entityID string, data map[string]string) error {
    var links []struct {
        Label string `json:"label"`
        URL   string `json:"url"`
    }
    json.Unmarshal([]byte(data["links"]), &links)
    return nil
},
```

For generic (no sub-fields) repeaters the array contains plain strings:

```json
["https://example.com/a.jpg","https://example.com/b.jpg"]
```
