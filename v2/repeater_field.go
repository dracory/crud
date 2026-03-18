package crud

import "github.com/dracory/form"

// RepeaterField wraps a form.FieldInterface and adds a list of sub-fields
// that define the structure of each repeater row.
// Use NewRepeaterField to construct one, then pass it anywhere a
// form.FieldInterface is accepted (createFields / updateFields).
type RepeaterField struct {
	form.FieldInterface
	Fields []form.FieldInterface
}

// NewRepeaterField creates a repeater field with the given base options and
// the sub-fields that make up each row.
func NewRepeaterField(opts form.FieldOptions, fields []form.FieldInterface) *RepeaterField {
	opts.Type = FORM_FIELD_TYPE_REPEATER
	return &RepeaterField{
		FieldInterface: form.NewField(opts),
		Fields:         fields,
	}
}
