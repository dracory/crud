package crud

import (
	"github.com/dracory/bs"
	"github.com/dracory/form"
	"github.com/dracory/hb"
)

func (crud *Crud) buildFieldInput(field form.FieldInterface) *hb.Tag {
	return hb.Input().
		Class("form-control").
		Attr("v-model", "entityModel."+field.GetName())
}

func (crud *Crud) buildFieldNumber(field form.FieldInterface) *hb.Tag {
	return hb.Input().
		Class("form-control").
		Type(hb.TYPE_NUMBER).
		Attr("v-model", "entityModel."+field.GetName())
}

func (crud *Crud) buildFieldPassword(field form.FieldInterface) *hb.Tag {
	return hb.Input().
		Class("form-control").
		Type(hb.TYPE_PASSWORD).
		Attr("v-model", "entityModel."+field.GetName())
}

func (crud *Crud) buildFieldTextarea(field form.FieldInterface) *hb.Tag {
	return hb.TextArea().
		Class("form-control").
		Attr("v-model", "entityModel."+field.GetName())
}

func (crud *Crud) buildFieldDatetime(field form.FieldInterface) *hb.Tag {
	return hb.NewTag(`el-date-picker`).
		Attr("type", "datetime").
		Attr("v-model", "entityModel."+field.GetName())
}

func (crud *Crud) buildFieldHtmlarea(field form.FieldInterface) *hb.Tag {
	return hb.NewTag("trumbowyg").
		Attr("v-model", "entityModel."+field.GetName()).
		Attr(":config", "trumbowigConfig").
		Class("form-control")
}

func (crud *Crud) buildFieldSelect(field form.FieldInterface) *hb.Tag {
	sel := hb.Select().Class("form-select").Attr("v-model", "entityModel."+field.GetName())
	for _, opt := range field.GetOptions() {
		sel.AddChild(hb.Option().Value(opt.Key).Text(opt.Value))
	}
	if field.GetOptionsF() != nil {
		for _, opt := range field.GetOptionsF()() {
			sel.AddChild(hb.Option().Value(opt.Key).Text(opt.Value))
		}
	}
	return sel
}

func (crud *Crud) buildFieldImage(field form.FieldInterface) *hb.Tag {
	fieldName := field.GetName()
	return hb.Div().Children([]hb.TagInterface{
		hb.Image("").
			Attr(`v-bind:src`, `entityModel.`+fieldName+`||'https://www.freeiconspng.com/uploads/no-image-icon-11.PNG'`).
			Style(`width:200px;`),
		bs.InputGroup().Children([]hb.TagInterface{
			hb.Input().Type(hb.TYPE_URL).Class("form-control").Attr("v-model", "entityModel."+fieldName),
			hb.If(crud.fileManagerURL != "", bs.InputGroupText().Children([]hb.TagInterface{
				hb.Hyperlink().Text("Browse").Href(crud.fileManagerURL).Target("_blank"),
			})),
		}),
	})
}

func (crud *Crud) buildFieldImageInline(field form.FieldInterface) *hb.Tag {
	fieldName := field.GetName()
	return hb.Div().Children([]hb.TagInterface{
		hb.Image("").
			Attr(`v-bind:src`, `entityModel.`+fieldName+`||'https://www.freeiconspng.com/uploads/no-image-icon-11.PNG'`).
			Style(`width:200px;`),
		hb.Input().
			Type(hb.TYPE_FILE).
			Attr("v-on:change", "uploadImage($event, '"+fieldName+"')").
			Attr("accept", "image/*"),
		hb.Button().
			HTML("See Image Data").
			Attr("v-on:click", "tmp.show_url_"+fieldName+" = !tmp.show_url_"+fieldName),
		hb.TextArea().
			Type(hb.TYPE_URL).
			Class("form-control").
			Attr("v-if", "tmp.show_url_"+fieldName).
			Attr("v-model", "entityModel."+fieldName),
	})
}

// buildFieldRepeater renders a Vue-driven repeater container for the given field.
//
// If the field declares sub-fields via GetFields(), each row is rendered as a
// typed sub-form (buildRepeaterRowFromFields) with one input per sub-field.
// The "Add Item" button pre-seeds the new row with all sub-field keys set to
// empty strings so Vue's reactivity picks them up immediately.
//
// If no sub-fields are declared, each row falls back to a single plain text
// input (buildRepeaterRowGeneric), treating the row as a simple string value.
//
// The rendered HTML relies on two Vue methods that must exist in the mounted
// app: addRepeaterItem(fieldName, item) and removeRepeaterItem(fieldName, index).
// These are provided by all three Vue apps in this package (EntityManager,
// EntityCreate, EntityUpdate).
func (crud *Crud) buildFieldRepeater(field form.FieldInterface) *hb.Tag {
	fieldName := field.GetName()

	// Typed sub-fields or generic single-value fallback.
	var rowContent hb.TagInterface
	emptyItem := `""` // generic: each row is a plain string
	subFields := field.GetFields()
	if len(subFields) > 0 {
		rowContent = crud.buildRepeaterRowFromFields(fieldName, subFields)
		// Pre-seed the new row JS object with all sub-field keys as empty strings.
		emptyItem = "{"
		for i, sf := range subFields {
			if i > 0 {
				emptyItem += ", "
			}
			emptyItem += `"` + sf.GetName() + `": ""`
		}
		emptyItem += "}"
	} else {
		rowContent = crud.buildRepeaterRowGeneric(fieldName)
	}

	return hb.Div().Class("repeater-container").Children([]hb.TagInterface{
		hb.Div().Class("repeater-items").
			Attr("v-for", "(item, index) in entityModel."+fieldName).
			Attr(":key", "index").
			Children([]hb.TagInterface{
				hb.Div().Class("repeater-item mb-3 border p-3").Children([]hb.TagInterface{
					hb.Div().Class("d-flex gap-1 float-end mb-2").Children([]hb.TagInterface{
						hb.Button().Class("btn btn-sm btn-outline-secondary").
							Attr("v-on:click", "moveRepeaterItemUp('"+fieldName+"', index)").
							Attr(":disabled", "index === 0").
							Child(hb.I().Class("bi-arrow-up")),
						hb.Button().Class("btn btn-sm btn-outline-secondary").
							Attr("v-on:click", "moveRepeaterItemDown('"+fieldName+"', index)").
							Attr(":disabled", "index === entityModel."+fieldName+".length - 1").
							Child(hb.I().Class("bi-arrow-down")),
						hb.Button().Class("btn btn-sm btn-outline-danger").
							Attr("v-on:click", "removeRepeaterItem('"+fieldName+"', index)").
							Child(hb.I().Class("bi-trash")),
					}),
					rowContent,
				}),
			}),
		hb.Button().Class("btn btn-sm btn-outline-primary").
			Attr("v-on:click", "addRepeaterItem('"+fieldName+"', "+emptyItem+")").
			Child(hb.I().Class("bi-plus")).
			HTML(" Add Item"),
	})
}

// buildRepeaterRowFromFields renders one repeater row as a Bootstrap grid row
// with a labeled input per sub-field. The v-model for each input is scoped to
// entityModel.<fieldName>[index].<subFieldName>, so Vue binds each sub-field
// independently within the row object.
//
// All field types supported by buildSubFieldWidget are available as sub-fields.
func (crud *Crud) buildRepeaterRowFromFields(fieldName string, fields []form.FieldInterface) *hb.Tag {
	row := hb.Div().Class("repeater-item-content row g-2")
	for _, subField := range fields {
		subName := subField.GetName()
		subLabel := subField.GetLabel()
		if subLabel == "" {
			subLabel = subName
		}
		vModel := "entityModel." + fieldName + "[index]." + subName
		input := crud.buildSubFieldWidget(subField, vModel)
		row.AddChild(hb.Div().Class("col").Children([]hb.TagInterface{
			hb.Label().Class("form-label").
				Text(subLabel).
				ChildIf(subField.GetRequired(), hb.Sup().Text("*").Class("text-danger ms-1")),
			input,
		}))
	}
	return row
}

// buildSubFieldWidget returns the appropriate Vue-bound input widget for a
// repeater sub-field, using the provided vModel path instead of the default
// top-level entityModel.<name> path.
//
// Supported types: string, number, password, email, tel, url, date, datetime,
// color, checkbox, hidden, htmlarea, image, radio, select, textarea, blockarea,
// blockeditor. Anything else falls back to a plain text input.
func (crud *Crud) buildSubFieldWidget(field form.FieldInterface, vModel string) *hb.Tag {
	switch field.GetType() {
	case FORM_FIELD_TYPE_TEXTAREA, FORM_FIELD_TYPE_BLOCKAREA, FORM_FIELD_TYPE_BLOCKEDITOR:
		return hb.TextArea().Class("form-control").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_NUMBER:
		return hb.Input().Class("form-control").Type(hb.TYPE_NUMBER).Attr("v-model", vModel)
	case FORM_FIELD_TYPE_PASSWORD:
		return hb.Input().Class("form-control").Type(hb.TYPE_PASSWORD).Attr("v-model", vModel)
	case FORM_FIELD_TYPE_DATE:
		return hb.NewTag(`el-date-picker`).Attr("type", "date").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_DATETIME:
		return hb.NewTag(`el-date-picker`).Attr("type", "datetime").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_EMAIL:
		return hb.Input().Class("form-control").Type("email").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_TEL:
		return hb.Input().Class("form-control").Type("tel").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_URL:
		return hb.Input().Class("form-control").Type("url").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_COLOR:
		return hb.Input().Class("form-control form-control-color").Type("color").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_CHECKBOX:
		return hb.Input().Class("form-check-input").Type("checkbox").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_HIDDEN:
		return hb.Input().Type("hidden").Attr("v-model", vModel)
	case FORM_FIELD_TYPE_HTMLAREA:
		return hb.NewTag("trumbowyg").Attr("v-model", vModel).Attr(":config", "trumbowigConfig").Class("form-control")
	case FORM_FIELD_TYPE_IMAGE:
		return hb.Div().Children([]hb.TagInterface{
			hb.Image("").
				Attr(`v-bind:src`, vModel+`||'https://www.freeiconspng.com/uploads/no-image-icon-11.PNG'`).
				Style(`width:200px;`),
			bs.InputGroup().Children([]hb.TagInterface{
				hb.Input().Type(hb.TYPE_URL).Class("form-control").Attr("v-model", vModel),
				hb.If(crud.fileManagerURL != "", bs.InputGroupText().Children([]hb.TagInterface{
					hb.Hyperlink().Text("Browse").Href(crud.fileManagerURL).Target("_blank"),
				})),
			}),
		})
	case FORM_FIELD_TYPE_RADIO:
		group := hb.Div().Class("d-flex flex-wrap gap-2")
		opts := field.GetOptions()
		if field.GetOptionsF() != nil {
			opts = append(opts, field.GetOptionsF()()...)
		}
		for _, opt := range opts {
			group.AddChild(hb.Div().Class("form-check").Children([]hb.TagInterface{
				hb.Input().Type("radio").Class("form-check-input").
					Attr("v-model", vModel).Value(opt.Key),
				hb.Label().Class("form-check-label").Text(opt.Value),
			}))
		}
		return group
	case FORM_FIELD_TYPE_SELECT:
		sel := hb.Select().Class("form-select").Attr("v-model", vModel)
		for _, opt := range field.GetOptions() {
			sel.AddChild(hb.Option().Value(opt.Key).Text(opt.Value))
		}
		if field.GetOptionsF() != nil {
			for _, opt := range field.GetOptionsF()() {
				sel.AddChild(hb.Option().Value(opt.Key).Text(opt.Value))
			}
		}
		return sel
	default: // FORM_FIELD_TYPE_STRING and anything else
		return hb.Input().Class("form-control").Attr("v-model", vModel)
	}
}

// buildRepeaterRowGeneric renders a single plain text input per row, used when
// the repeater has no declared sub-fields. The entire row is treated as a
// scalar string value bound to entityModel.<fieldName>[index].
func (crud *Crud) buildRepeaterRowGeneric(fieldName string) *hb.Tag {
	return hb.Input().
		Class("form-control").
		Attr("v-model", "entityModel."+fieldName+"[index]")
}

func (crud *Crud) buildFieldRaw(field form.FieldInterface) *hb.Tag {
	return hb.Raw(field.GetValue())
}

func (crud *Crud) buildFieldRadio(field form.FieldInterface) *hb.Tag {
	group := hb.Div().Class("d-flex flex-wrap gap-3")
	opts := field.GetOptions()
	if field.GetOptionsF() != nil {
		opts = append(opts, field.GetOptionsF()()...)
	}
	for _, opt := range opts {
		id := field.GetName() + "_" + opt.Key
		group.AddChild(hb.Div().Class("form-check").Children([]hb.TagInterface{
			hb.Input().Type("radio").Class("form-check-input").
				ID(id).
				Attr("v-model", "entityModel."+field.GetName()).
				Value(opt.Key),
			hb.Label().Class("form-check-label").Attr("for", id).Text(opt.Value),
		}))
	}
	return group
}

// buildFieldWidget resolves the correct input widget for a given field type and sets its ID.
func (crud *Crud) buildFieldWidget(field form.FieldInterface, fieldID string) *hb.Tag {
	var widget *hb.Tag
	switch field.GetType() {
	case FORM_FIELD_TYPE_BLOCKAREA:
		widget = crud.buildFieldTextarea(field)
	case FORM_FIELD_TYPE_BLOCKEDITOR:
		widget = crud.buildFieldTextarea(field) // rendered via blockarea JS init
	case FORM_FIELD_TYPE_CHECKBOX:
		widget = hb.Input().Type("checkbox").Class("form-check-input").Attr("v-model", "entityModel."+field.GetName())
	case FORM_FIELD_TYPE_COLOR:
		widget = hb.Input().Type("color").Class("form-control form-control-color").Attr("v-model", "entityModel."+field.GetName())
	case FORM_FIELD_TYPE_DATE:
		widget = hb.Input().Type("date").Class("form-control").Attr("v-model", "entityModel."+field.GetName())
	case FORM_FIELD_TYPE_DATETIME:
		widget = crud.buildFieldDatetime(field)
	case FORM_FIELD_TYPE_EMAIL:
		widget = hb.Input().Type("email").Class("form-control").Attr("v-model", "entityModel."+field.GetName())
	case FORM_FIELD_TYPE_FILE:
		widget = hb.Input().Type("file").Class("form-control").Attr("v-on:change", "uploadFile($event, '"+field.GetName()+"')")
	case FORM_FIELD_TYPE_HIDDEN:
		widget = hb.Input().Type("hidden").Attr("v-model", "entityModel."+field.GetName())
	case FORM_FIELD_TYPE_HTMLAREA:
		widget = crud.buildFieldHtmlarea(field)
	case FORM_FIELD_TYPE_IMAGE:
		widget = crud.buildFieldImage(field)
	case FORM_FIELD_TYPE_IMAGE_INLINE:
		widget = crud.buildFieldImageInline(field)
	case FORM_FIELD_TYPE_NUMBER:
		widget = crud.buildFieldNumber(field)
	case FORM_FIELD_TYPE_PASSWORD:
		widget = crud.buildFieldPassword(field)
	case FORM_FIELD_TYPE_RADIO:
		widget = crud.buildFieldRadio(field)
	case FORM_FIELD_TYPE_RAW:
		widget = crud.buildFieldRaw(field)
	case FORM_FIELD_TYPE_REPEATER:
		widget = crud.buildFieldRepeater(field)
	case FORM_FIELD_TYPE_SELECT:
		widget = crud.buildFieldSelect(field)
	case FORM_FIELD_TYPE_TABLE:
		widget = hb.Raw(field.GetValue()) // caller provides pre-rendered HTML
	case FORM_FIELD_TYPE_TEL:
		widget = hb.Input().Type("tel").Class("form-control").Attr("v-model", "entityModel."+field.GetName())
	case FORM_FIELD_TYPE_TEXTAREA:
		widget = crud.buildFieldTextarea(field)
	case FORM_FIELD_TYPE_URL:
		widget = hb.Input().Type("url").Class("form-control").Attr("v-model", "entityModel."+field.GetName())
	default: // FORM_FIELD_TYPE_STRING and anything unknown
		widget = crud.buildFieldInput(field)
	}
	widget.ID(fieldID)
	return widget
}
