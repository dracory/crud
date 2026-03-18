package crud

// FieldOption represents a key-value pair for select, radio, and other option-based fields.
type FieldOption struct {
	Key   string
	Value string
}

// FieldInterface defines the contract for fields used in crud create/update forms.
type FieldInterface interface {
	GetID() string
	GetLabel() string
	GetName() string
	GetHelp() string
	GetOptions() []FieldOption
	GetOptionsF() func() []FieldOption
	GetRequired() bool
	GetType() string
	GetValue() string
	// GetFields returns sub-fields for repeater fields; nil for all other types.
	GetFields() []FieldInterface
	// GetFuncValuesProcess returns a function that converts a raw string value
	// (e.g. JSON from FuncFetchUpdateData) into []map[string]string rows for Vue binding.
	// Returns nil for field types that do not need pre-processing.
	GetFuncValuesProcess() func(raw string) []map[string]string
}

// FieldOptions configures a new Field instance.
type FieldOptions struct {
	ID                string
	Type              string
	Name              string
	Label             string
	Help              string
	Options           []FieldOption
	OptionsF          func() []FieldOption
	Value             string
	Required          bool
	Fields            []FieldInterface
	FuncValuesProcess func(raw string) []map[string]string
}

// Field is the default FieldInterface implementation.
type Field struct {
	id                string
	fieldType         string
	name              string
	label             string
	help              string
	options           []FieldOption
	optionsF          func() []FieldOption
	value             string
	required          bool
	fields            []FieldInterface
	funcValuesProcess func(raw string) []map[string]string
}

// NewField creates a new Field with the given options.
func NewField(opts FieldOptions) *Field {
	return &Field{
		id:                opts.ID,
		fieldType:         opts.Type,
		name:              opts.Name,
		label:             opts.Label,
		help:              opts.Help,
		options:           opts.Options,
		optionsF:          opts.OptionsF,
		value:             opts.Value,
		required:          opts.Required,
		fields:            opts.Fields,
		funcValuesProcess: opts.FuncValuesProcess,
	}
}

func (f *Field) GetID() string                                        { return f.id }
func (f *Field) GetType() string                                      { return f.fieldType }
func (f *Field) GetName() string                                      { return f.name }
func (f *Field) GetLabel() string                                     { return f.label }
func (f *Field) GetHelp() string                                      { return f.help }
func (f *Field) GetOptions() []FieldOption                            { return f.options }
func (f *Field) GetOptionsF() func() []FieldOption                    { return f.optionsF }
func (f *Field) GetValue() string                                     { return f.value }
func (f *Field) GetRequired() bool                                    { return f.required }
func (f *Field) GetFields() []FieldInterface                          { return f.fields }
func (f *Field) GetFuncValuesProcess() func(raw string) []map[string]string {
	return f.funcValuesProcess
}

// RepeaterOptions configures a new repeater field instance.
type RepeaterOptions struct {
	Label             string
	Name              string
	Help              string
	Fields            []FieldInterface
	// FuncValuesProcess converts a raw JSON string (from FuncFetchUpdateData)
	// into []map[string]string rows for Vue binding. Optional.
	FuncValuesProcess func(raw string) []map[string]string
}

// NewRepeater creates a repeater FieldInterface for use in crud create/update forms.
func NewRepeater(opts RepeaterOptions) *Field {
	return &Field{
		fieldType:         FORM_FIELD_TYPE_REPEATER,
		name:              opts.Name,
		label:             opts.Label,
		help:              opts.Help,
		fields:            opts.Fields,
		funcValuesProcess: opts.FuncValuesProcess,
	}
}

// == Fluent builder methods ==================================================

func (f *Field) WithID(id string) *Field                { f.id = id; return f }
func (f *Field) WithName(name string) *Field            { f.name = name; return f }
func (f *Field) WithLabel(label string) *Field          { f.label = label; return f }
func (f *Field) WithHelp(help string) *Field            { f.help = help; return f }
func (f *Field) WithType(t string) *Field               { f.fieldType = t; return f }
func (f *Field) WithValue(v string) *Field              { f.value = v; return f }
func (f *Field) WithRequired() *Field                   { f.required = true; return f }
func (f *Field) WithOptions(opts ...FieldOption) *Field { f.options = opts; return f }
func (f *Field) WithOptionsF(fn func() []FieldOption) *Field {
	f.optionsF = fn
	return f
}
func (f *Field) WithFields(fields ...FieldInterface) *Field { f.fields = fields; return f }
func (f *Field) WithFuncValuesProcess(fn func(raw string) []map[string]string) *Field {
	f.funcValuesProcess = fn
	return f
}

// == Convenience constructors ================================================

func NewStringField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_STRING})
}
func NewEmailField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_EMAIL})
}
func NewNumberField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_NUMBER})
}
func NewPasswordField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_PASSWORD})
}
func NewHiddenField(name, value string) *Field {
	return NewField(FieldOptions{Name: name, Value: value, Type: FORM_FIELD_TYPE_HIDDEN})
}
func NewDateField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_DATE})
}
func NewDateTimeField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_DATETIME})
}
func NewSelectField(name, label string, options []FieldOption) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_SELECT, Options: options})
}
func NewTextAreaField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_TEXTAREA})
}
func NewCheckboxField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_CHECKBOX})
}
func NewRadioField(name, label string, options []FieldOption) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_RADIO, Options: options})
}
func NewFileField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_FILE})
}
func NewImageField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_IMAGE})
}
func NewColorField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_COLOR})
}
func NewTelField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_TEL})
}
func NewURLField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_URL})
}
func NewHtmlAreaField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_HTMLAREA})
}
func NewRawField(value string) *Field {
	return NewField(FieldOptions{Value: value, Type: FORM_FIELD_TYPE_RAW})
}

// == Element Plus specific constructors ======================================

func NewColorElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_COLOR_EL})
}
func NewDateElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_DATE_EL})
}
func NewDateTimeElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_DATETIME_EL})
}
func NewInputNumberElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_INPUT_NUMBER_EL})
}
func NewRateElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_RATE_EL})
}
func NewSelectElField(name, label string, options []FieldOption) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_SELECT_EL, Options: options})
}
func NewSliderElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_SLIDER_EL})
}
func NewSwitchElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_SWITCH_EL})
}
func NewTimeElField(name, label string) *Field {
	return NewField(FieldOptions{Name: name, Label: label, Type: FORM_FIELD_TYPE_TIME_EL})
}
