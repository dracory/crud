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

func (crud *Crud) buildFieldRepeater(field form.FieldInterface) *hb.Tag {
	fieldName := field.GetName()

	// Typed sub-fields or generic key/value fallback.
	var rowContent hb.TagInterface
	emptyItem := "{}"
	if rf, ok := field.(*RepeaterField); ok && len(rf.Fields) > 0 {
		rowContent = crud.buildRepeaterRowFromFields(fieldName, rf.Fields)
		// Pre-seed the new row JS object with all sub-field keys as empty strings.
		emptyItem = "{"
		for i, sf := range rf.Fields {
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
					hb.Button().Class("btn btn-sm btn-outline-danger float-end").
						Attr("v-on:click", "removeRepeaterItem('"+fieldName+"', index)").
						Child(hb.I().Class("bi-trash")).
						HTML(" Remove"),
					rowContent,
				}),
			}),
		hb.Button().Class("btn btn-sm btn-outline-primary").
			Attr("v-on:click", "addRepeaterItem('"+fieldName+"', "+emptyItem+")").
			Child(hb.I().Class("bi-plus")).
			HTML(" Add Item"),
	})
}

// buildRepeaterRowFromFields renders a typed sub-form row using the declared sub-fields.
func (crud *Crud) buildRepeaterRowFromFields(fieldName string, fields []form.FieldInterface) *hb.Tag {
	row := hb.Div().Class("repeater-item-content row g-2")
	for _, subField := range fields {
		subName := subField.GetName()
		subLabel := subField.GetLabel()
		if subLabel == "" {
			subLabel = subName
		}
		vModel := "entityModel." + fieldName + "[index]." + subName

		var input hb.TagInterface
		switch subField.GetType() {
		case FORM_FIELD_TYPE_TEXTAREA:
			input = hb.TextArea().Class("form-control").Attr("v-model", vModel)
		case FORM_FIELD_TYPE_NUMBER:
			input = hb.Input().Class("form-control").Type(hb.TYPE_NUMBER).Attr("v-model", vModel)
		case FORM_FIELD_TYPE_SELECT:
			sel := hb.Select().Class("form-select").Attr("v-model", vModel)
			for _, opt := range subField.GetOptions() {
				sel.AddChild(hb.Option().Value(opt.Key).Text(opt.Value))
			}
			if subField.GetOptionsF() != nil {
				for _, opt := range subField.GetOptionsF()() {
					sel.AddChild(hb.Option().Value(opt.Key).Text(opt.Value))
				}
			}
			input = sel
		default:
			input = hb.Input().Class("form-control").Attr("v-model", vModel)
		}

		row.AddChild(hb.Div().Class("col").Children([]hb.TagInterface{
			hb.Label().Class("form-label").
				Text(subLabel).
				ChildIf(subField.GetRequired(), hb.Sup().Text("*").Class("text-danger ms-1")),
			input,
		}))
	}
	return row
}

// buildRepeaterRowGeneric renders a generic key/value row when no sub-fields are declared.
func (crud *Crud) buildRepeaterRowGeneric(fieldName string) *hb.Tag {
	return hb.Div().Class("repeater-item-content").
		Attr("v-for", "(value, key) in item").
		Attr(":key", "key").
		Children([]hb.TagInterface{
			hb.Label().Class("form-label").Attr("v-text", "key"),
			hb.Input().Class("form-control").Attr("v-model", "entityModel."+fieldName+"[index][key]"),
		})
}

func (crud *Crud) buildFieldRaw(field form.FieldInterface) *hb.Tag {
	return hb.Raw(field.GetValue())
}

// buildFieldWidget resolves the correct input widget for a given field type and sets its ID.
func (crud *Crud) buildFieldWidget(field form.FieldInterface, fieldID string) *hb.Tag {
	var widget *hb.Tag
	switch field.GetType() {
	case FORM_FIELD_TYPE_IMAGE:
		widget = crud.buildFieldImage(field)
	case FORM_FIELD_TYPE_IMAGE_INLINE:
		widget = crud.buildFieldImageInline(field)
	case FORM_FIELD_TYPE_DATETIME:
		widget = crud.buildFieldDatetime(field)
	case FORM_FIELD_TYPE_HTMLAREA:
		widget = crud.buildFieldHtmlarea(field)
	case FORM_FIELD_TYPE_NUMBER:
		widget = crud.buildFieldNumber(field)
	case FORM_FIELD_TYPE_PASSWORD:
		widget = crud.buildFieldPassword(field)
	case FORM_FIELD_TYPE_SELECT:
		widget = crud.buildFieldSelect(field)
	case FORM_FIELD_TYPE_TEXTAREA, FORM_FIELD_TYPE_BLOCKAREA:
		widget = crud.buildFieldTextarea(field)
	case FORM_FIELD_TYPE_RAW:
		widget = crud.buildFieldRaw(field)
	case FORM_FIELD_TYPE_REPEATER:
		widget = crud.buildFieldRepeater(field)
	default:
		widget = crud.buildFieldInput(field)
	}
	widget.ID(fieldID)
	return widget
}
